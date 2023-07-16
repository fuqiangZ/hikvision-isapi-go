package hikvision

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
	"golang.org/x/net/context"
)

/**
* 布防的接口是长连接
* 第一条数据是布防成功；
* 如果没有数据，则每隔几秒钟有一个心跳数据；
* 如果有车牌数据上报，则第一条是xml数据，包含有n个图片，紧接着n个图片都数据。
* 每段数据都是以\r\n结束，所以解析时候需要两层配套
* 第一层 nextPart，
* 第二层 parsePart, 根据contentLen读取数据，再根据xml内容读取后续字段
**/
type Guard struct {
	eventId string
	resp    *http.Response
	ctx     context.Context
	Output  chan Content
}

func NewGuard(ctx context.Context, host, username, password string) (*Guard, error) {
	r := NewDigestRequest(ctx, username, password)
	reqBody := bytes.NewBuffer([]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<SubscribeEvent >
		<heartbeat>5</heartbeat>
		<channelMode>all</channelMode>
		<eventMode>list</eventMode>
		<EventList>
			<Event>
				<type>ANPR</type>
				<pictureURLType opt="binary,localURL,cloudStorageURL" def="binary"/>
			</Event>
		</EventList>
		<level>middle</level>
		<pictureURLType opt="binary,localURL,cloudStorageURL" def="binary"/>
	</SubscribeEvent>`))

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", host, "/ISAPI/Event/notification/subscribeEvent"), reqBody)
	// req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Connection", "keep-alive")

	resp, err := r.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp is nil")
	}
	// defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("alertStream failed, statusCode: %v", resp.StatusCode)
	}

	g := &Guard{resp: resp, ctx: ctx}
	g.Output = make(chan Content, 1)
	return g, nil
}

// func NewGuard2(ctx context.Context, host, username, password string) (*Guard, error) {
// 	r := NewDigestRequest(ctx, username, password)
// 	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", host, "/ISAPI/Event/notification/alertStream"), nil)
// 	// req.WithContext(ctx)
// 	resp, err := r.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp == nil {
// 		return nil, fmt.Errorf("resp is nil")
// 	}
// 	// defer resp.Body.Close()
// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("alertStream failed, statusCode: %v", resp.StatusCode)
// 	}

//		g := &Guard{resp: resp, ctx: ctx}
//		g.Output = make(chan Content, 1)
//		return g, nil
//	}
func (g *Guard) Start(msg chan<- Message) error {
	defer func() {
		if g.resp.Body != nil {
			g.resp.Body.Close()
		}
		if g.Output != nil {
			close(g.Output)
		}
	}()

	//buf转化一把，确认订阅成功
	mediaType, params, err := mime.ParseMediaType(g.resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	// fmt.Printf("%v \n", params)

	//布防之后的响应处理事件
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := NewMultipart(g.ctx, g.resp.Body, params["boundary"])
		var m *Message

		//red first
		err, c := mr.NextPart()
		if err != nil {
			fmt.Println("finish ")
			return err
		}
		contentT := c.Header.Get("Content-Type").(string)
		if contentT != TYPE_XML {
			return errors.New("invalid content type")
		}
		if sid, err := readSubscribeId(c.Body); err != nil || sid == "" {
			return errors.New("invalid first message")
		} else {
			g.eventId = sid
		}

		for {
			select {
			case <-g.ctx.Done():
				fmt.Println("user canceled")
				return nil
			default:

			}
			err, c := mr.NextPart()
			if err == io.EOF {
				fmt.Println("finish ")
				break
			}
			if err != nil {
				log.Fatal("2 :", err)
			}
			fmt.Printf("%s \r\n", c.Body)

			err = Parse(c.Header.Get("Content-Type").(string), c.Body, &m)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// g.Output <- *c
			if m.EventType == EVENT_TYPE_HEARTBEAT || (m.EventType == EVENT_TYPE_ANPR && len(m.Attachment) == m.AttachNum) {
				//将数据output出去
				// h.Message <- *m
				msg <- *m
			}
		}
	}
	fmt.Printf("%v\n", g.resp.StatusCode)
	return nil
}

func readSubscribeId(body []byte) (string, error) {
	doc, err := xmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	root, err := xmlquery.Query(doc, "SubscribeEventResponse")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	n := root.SelectElement("id")
	if n == nil {
		fmt.Println("not find id field")
		return "", errors.New("not find id field")
	}
	return n.InnerText(), nil
}
func Parse(contentType string, body []byte, m **Message) error {
	if contentType == TYPE_XML {
		var err error
		doc, err := xmlquery.Parse(bytes.NewReader(body))
		if err != nil {
			fmt.Println(err)
			return err
		}
		root, err := xmlquery.Query(doc, "EventNotificationAlert")
		if err != nil {
			fmt.Println(err)
			return err
		}
		n := root.SelectElement("eventType")
		if n == nil {
			fmt.Println("not find eventType field")
			return errors.New("not find eventType field")
		}
		eventType := n.InnerText()

		picNum := 0
		if eventType == EVENT_TYPE_ANPR {
			n = root.SelectElement("picNum")
			if n == nil {
				fmt.Println("not find picNum field")
				return errors.New("not find picNum field")
			}
			picNum, err = strconv.Atoi(n.InnerText())
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		newMsg := Message{EventType: eventType, KeyContent: body, AttachNum: picNum}
		*m = &newMsg
	} else if contentType == TYPE_IMAGE {
		if *m == nil || (*m).EventType != EVENT_TYPE_ANPR {
			return nil
		}
		h := make(HeaderType, 1)
		h[ContentT] = contentType
		h[ContentL] = len(body)
		nc := Content{Header: h, Body: body}
		(*m).Attachment = append((*m).Attachment, nc)
	}
	return nil
}
func (g *Guard) Stop() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print("[guard stop]recover from error")
		}
	}()
	// if g.resp.Body != nil {
	// 	g.resp.Body.Close()
	// }

}
