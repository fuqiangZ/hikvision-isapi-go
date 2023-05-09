package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"testing"

	"github.com/antchfx/xmlquery"
	hikvision "github.com/fuqiangZ/hikvision-isapi-go"
)

func TestDataParse(t *testing.T) {
	fmt.Println("dd")
	mr := hikvision.NewMultipart(context.Background(), bytes.NewReader(MockFromFile()), "-----------------------7e13971310878")

	var m *hikvision.Message
	for {
		err, b := mr.NextPart()
		if err == io.EOF {
			fmt.Println("finish ")
			break
		}
		if err != nil {
			log.Fatal("2 :", err)
		}
		// fmt.Printf("%v \r\n", *c)
		// fmt.Printf("%q \n", c.Body)

		if b.ContentType == "xml" {
			doc, err := xmlquery.Parse(bytes.NewReader(b.Body))
			fmt.Println(err)

			e, err := xmlquery.Query(doc, "EventNotificationAlert")
			fmt.Println(err)

			if n := e.SelectElement("eventType"); n != nil {
				fmt.Println("eventType", n.InnerText())
				if n.InnerText() == "ANPR" {
					n = e.SelectElement("picNum")
					fmt.Println(n.InnerText())
				}
			}

		}

		if b.ContentType == "xml" {
			var err error
			doc, err := xmlquery.Parse(bytes.NewReader(b.Body))
			if err != nil {
				fmt.Println(err)
				continue
			}
			root, err := xmlquery.Query(doc, "EventNotificationAlert")
			if err != nil {
				fmt.Println(err)
				continue
			}
			n := root.SelectElement("eventType")
			if n == nil {
				fmt.Println("not find eventType field")
				continue
			}
			eventType := n.InnerText()

			picNum := 0
			if eventType == "ANPR" {
				n = root.SelectElement("picNum")
				if n == nil {
					fmt.Println("not find picNum field")
					continue
				}
				picNum, err = strconv.Atoi(n.InnerText())
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
			m = &hikvision.Message{EventType: eventType, KeyContent: b.Body, AttachNum: picNum}
		} else if b.ContentType == hikvision.TYPE_IMAGE {
			if m == nil || m.EventType != hikvision.EVENT_TYPE_ANPR {
				continue
			}
			m.Attachment = append(m.Attachment, hikvision.Content{ContentType: b.ContentType, ContentLen: len(b.Body), Body: b.Body})
		}
		if m.EventType == "heatBeat" || (m.EventType == "ANPR" && len(m.Attachment) == m.AttachNum) {
			//将数据output出去
			fmt.Println(m)
			b, err := json.Marshal(m)
			fmt.Println(err)
			fmt.Println(string(b))
		}

	}
}
