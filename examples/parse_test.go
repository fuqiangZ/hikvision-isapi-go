package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/antchfx/xmlquery"
	hikvision "github.com/fuqiangZ/hikvision-isapi-go"
)

func TestDataParse(t *testing.T) {
	fmt.Println("dd")
	mr := hikvision.NewMultipart(context.Background(), bytes.NewReader(MockFromFile()), "-----------------------7e13971310878")

	for {
		err, c := mr.NextPart()
		if err == io.EOF {
			fmt.Println("finish ")
			break
		}
		if err != nil {
			log.Fatal("2 :", err)
		}
		// fmt.Printf("%v \r\n", *c)
		// fmt.Printf("%q \n", c.Body)

		if c.ContentType == "xml" {
			doc, err := xmlquery.Parse(bytes.NewReader(c.Body))
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

	}
}
