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
	KeyContent xmlquery.Node //主数据
	Attachment []Content     //依据主数据读取出来的附加内容
}
type Content struct {
	ContentType string
	ContentLen  int
	Body        []byte
}

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
			var err error
			h.g, err = NewGuard(h.ctx, h.host, h.username, h.password)
			if err != nil {
				fmt.Println(err)
				time.Sleep(1 * time.Second)
				continue
			}
			go h.g.Start()
			for {
				select {
				case b := <-h.g.Output:
					fmt.Println(b.ContentType)
					if b.ContentType == TYPE_XML {
						//解析xml内容，再读取n个图片
					}
				// m := Message{
				// 	EventType: "hearbeat",
				// 	Key
				// }
				// h.Message <- b
				case <-h.ctx.Done():
					fmt.Println("for data canceled")
					return
				}

			}

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
