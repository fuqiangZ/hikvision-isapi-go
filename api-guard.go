package hikvision

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"

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
	resp   *http.Response
	ctx    context.Context
	Output chan Content
}

func NewGuard(ctx context.Context, host, username, password string) (*Guard, error) {
	r := NewDigestRequest(ctx, username, password)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", host, "/ISAPI/Event/notification/alertStream"), nil)
	// req.WithContext(ctx)
	resp, err := r.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp is nil")
	}
	// defer resp.Body.Close()

	g := &Guard{resp: resp, ctx: ctx}
	g.Output = make(chan Content, 1)
	return g, nil
}
func (g *Guard) Start() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from:", err)
		}

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
			// fmt.Printf("%v \r\n", *c)

			g.Output <- *c
		}
	}
	fmt.Printf("%v\n", g.resp.StatusCode)
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
