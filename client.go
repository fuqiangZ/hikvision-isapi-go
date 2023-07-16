package hikvision

import (
	"context"
	"fmt"
	"time"

	"github.com/antchfx/xmlquery"
)

type HikvisionClient struct {
	host     string
	username string
	password string
	ctx      context.Context
	cancel   context.CancelFunc
	Message  chan Message
	g        *Guard
}

type Message struct {
	EventType  string
	KeyContent []byte
	AttachNum  int
	Attachment []Content //依据主数据读取出来的附加内容
}
type Content struct {
	Header HeaderType
	// ContentType string
	// ContentLen int
	Body []byte
}

type HeaderType map[string]interface{}

func (h HeaderType) Get(key string) interface{} {
	if v, ok := h[key]; ok {
		return v
	}
	return nil
}

const (
	EVENT_TYPE_HEARTBEAT = "heartBeat"
	EVENT_TYPE_ANPR      = "ANPR"
)

func NewClient(host, username, password string) *HikvisionClient {
	h := &HikvisionClient{host: host, username: username, password: password}
	h.ctx, h.cancel = context.WithCancel(context.Background())
	h.Message = make(chan Message, 1)
	return h
}

// 布防能力检查
func (h *HikvisionClient) GetCapablities() {
	// h.cancel()
}

// 布防能力检查
func (h *HikvisionClient) GateawayControl(xmldata string) (*xmlquery.Node, error) {
	// h.cancel()
	return GateawayControl(h.host, h.username, h.password, xmldata)
}

func (h *HikvisionClient) LedShowInfo(xmldata string) (*xmlquery.Node, error) {
	// h.cancel()
	return LedShowInfo(h.host, h.username, h.password, xmldata)
}

// 布防
func (h *HikvisionClient) StartAlarmGuard() {
	go func() {
		for {

			select {
			case <-h.ctx.Done():
				fmt.Println("alarm start canceled")
				return
			default:

			}
			fmt.Println("start guard")
			g, err := NewGuard(h.ctx, h.host, h.username, h.password)
			if err != nil {
				fmt.Println(err)
				time.Sleep(1 * time.Second)
				continue
			}
			g.Start(h.Message)

			// var m *Message
			// for {
			// 	select {
			// 	case b := <-h.g.Output:
			// 		fmt.Println(b.Header)
			// 		contentType := b.Header.Get("Content-Type")
			// 		if contentType == TYPE_XML {
			// 			var err error
			// 			doc, err := xmlquery.Parse(bytes.NewReader(b.Body))
			// 			if err != nil {
			// 				fmt.Println(err)
			// 				continue
			// 			}
			// 			root, err := xmlquery.Query(doc, "EventNotificationAlert")
			// 			if err != nil {
			// 				fmt.Println(err)
			// 				continue
			// 			}
			// 			n := root.SelectElement("eventType")
			// 			if n == nil {
			// 				fmt.Println("not find eventType field")
			// 				continue
			// 			}
			// 			eventType := n.InnerText()

			// 			picNum := 0
			// 			if eventType == EVENT_TYPE_ANPR {
			// 				n = root.SelectElement("picNum")
			// 				if n == nil {
			// 					fmt.Println("not find picNum field")
			// 					continue
			// 				}
			// 				picNum, err = strconv.Atoi(n.InnerText())
			// 				if err != nil {
			// 					fmt.Println(err)
			// 					continue
			// 				}
			// 			}
			// 			m = &Message{EventType: eventType, KeyContent: b.Body, AttachNum: picNum}
			// 		} else if contentType == TYPE_IMAGE {
			// 			if m == nil || m.EventType != EVENT_TYPE_ANPR {
			// 				continue
			// 			}
			// 			h := make(HeaderType, 1)
			// 			h[ContentT] = contentType
			// 			h[ContentL] = len(b.Body)
			// 			nc := Content{Header: h, Body: b.Body}
			// 			m.Attachment = append(m.Attachment, nc)
			// 		}
			// 		if m.EventType == EVENT_TYPE_HEARTBEAT || (m.EventType == EVENT_TYPE_ANPR && len(m.Attachment) == m.AttachNum) {
			// 			//将数据output出去
			// 			h.Message <- *m
			// 		}
			// 	case <-h.ctx.Done():
			// 		fmt.Println("for data canceled")
			// 		return
			// 	}
			// }
		}

	}()
}

// 撤销布防
func (h *HikvisionClient) StopAlarmGuard() {
	h.cancel()
	if h.g != nil {
		h.g.Stop()
	}

	close(h.Message)
}

// 设备信息
func (h *HikvisionClient) DeviceInfo() (*xmlquery.Node, error) {
	return GetDeviceInfo(h.host, h.username, h.password)
}

// 设备信息
func (h *HikvisionClient) GetSystemCapability() (*xmlquery.Node, error) {
	return GetSystemCapability(h.host, h.username, h.password)
}
