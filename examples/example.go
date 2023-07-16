package main

import (
	"encoding/json"
	"fmt"

	hikvision "github.com/fuqiangZ/hikvision-isapi-go"
)

func main() {
	client := hikvision.NewClient("http://192.168.100.248:80", "admin", "Abc12345")

	fmt.Println("alarm guard start")
	client.StartAlarmGuard()

	for {
		m := <-client.Message
		fmt.Println(m)
		b, err := json.Marshal(m)
		fmt.Println(err)
		fmt.Println(string(b))
	}

	// time.Sleep(6 * time.Second)

	// client.StopAlarmGuard()
	// fmt.Println("alarm guard stopped")

	// time.Sleep(2 * time.Second)
	// guard2()

	// select {}
}

func guard2() {
	client := hikvision.NewClient("http://localhost:8800", "admin", "Abc12345")

	fmt.Println("alarm guard start")
	client.StartAlarmGuard()

	for {
		m := <-client.Message
		fmt.Println(m)
	}

	// time.Sleep(6 * time.Second)

	// // client.StopAlarmGuard()
	// fmt.Println("alarm guard stopped")

}
