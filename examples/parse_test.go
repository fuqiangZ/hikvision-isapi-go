package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"testing"

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
		fmt.Printf("%q \n", c.Body)
	}
}
