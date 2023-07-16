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

func TestLedDisplay(t *testing.T) {

	xmlData := `{
		"SingleSceneLEDConfigurations": {	
			"sid":  1,	
			"mode":  "passingVehicle",	
			"passingIconType":  "noPassing",	
			"showFreeEnabled":  true,	
			"displayTime":  15,	
			"vehicleDisplayEnabled":  true,	
			"allowListDisplayEnabled":  false,	
			"blockListDisplayEnabled":  false,	
			"temporaryListDisplayEnabled":  false,	
			"LEDConfigurationList": [	
				{
					"LEDConfiguration": {	
						"id":  1,	
						"brightness":  10,	
						"enabled":  true,	
						"ShowInfoList": [	
							{
								"ShowInfo": {	
									"id":  1,	
									"fontSize":  16,	
									"fontColor":  "red",	
									"speedType":  "medium",	
									"displayMode":  "left",	
									"LineInfoList": [	
										{
											"LineInfo": {	
												"id":  1,	
												"value":  "这一行是什么",	
												"customValue":  "aaa应该可以了吧"	
											}
										}
									]
								}
							}
						]
					}
				}
			],
		   "LedInfo": {
		"communicateMode": "network"},
			"speedMode":  "slowMode"	
		}
	}
	`
	client := hikvision.NewClient("http://192.168.100.248", "admin", "Abc12345")

	n, err := client.LedShowInfo(xmlData)
	fmt.Println(err)
	// xmlquery.Query()
	fmt.Println(n.Data)

}
func TestGateawayControl(t *testing.T) {
	client := hikvision.NewClient("http://192.168.100.248", "admin", "Abc12345")
	xmldata := `
<?xml version="1.0" encoding="UTF-8"?>
<BarrierGate xmlns="http://www.isapi.org/ver20/XMLSchema" version="2.0">
  <ctrlMode opt="open,close,lock,unlock">open</ctrlMode>
</BarrierGate>`
	n, err := client.GateawayControl(xmldata)
	fmt.Println(err)
	// xmlquery.Query()
	fmt.Println(n.Data)

}

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

		fmt.Println(b.Header)
		contentType := b.Header.Get("Content-Type")
		if contentType == "xml" {
			fmt.Println(string(b.Body))
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

		if contentType == "xml" {
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
				anprFilter, err := xmlquery.Query(doc, "/EventNotificationAlert/ANPR/licensePlate")
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(anprFilter.InnerText())
				//  anprFilter:= root.SelectElement("/ANPR/licensePlate")
				// if n == nil {
				// 	fmt.Println("no filter licensePlate")
				// 	continue
				// }
				fmt.Println(anprFilter.InnerText())
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
		} else if contentType == hikvision.TYPE_IMAGE {
			if m == nil || m.EventType != hikvision.EVENT_TYPE_ANPR {
				continue
			}
			h := make(hikvision.HeaderType, 1)
			h[hikvision.ContentT[:len(hikvision.ContentT)-1]] = contentType
			h[hikvision.ContentL[:len(hikvision.ContentL)-1]] = len(b.Body)
			nc := hikvision.Content{Header: h, Body: b.Body}
			m.Attachment = append(m.Attachment, nc)
			// m.Attachment = append(m.Attachment, hikvision.Content{ContentType: b.ContentType, ContentLen: len(b.Body), Body: b.Body})
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
