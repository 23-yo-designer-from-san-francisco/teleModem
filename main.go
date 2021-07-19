package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type SMS struct {
	ID string
	Number string
	Content string
	Tag string
	Date string
	DraftGroupID string
	ReceivedAllConcatSMS string
	ConcatSMSTotal string
	ConcatSMSReceived string
	SMSClass string
	SMSMem string
	SMSSubmitMsgRef string
}

func main() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://192.168.0.1/goform/goform_get_cmd_process?cmd=sms_data_total&page=0&data_per_page=500&mem_store=1&tags=10&order_by=order+by+id+asc", nil)
	req.Header.Set("Referer", "http://192.168.0.1/index.html")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
    body = body[12:len(body)-1]  // Cut the {"messages": ... } part, keeping the [] of messages

	var sms []SMS

	err := json.Unmarshal(body, &sms)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for _, msg := range sms {
		fmt.Printf("ID: %s\n%s\n", msg.ID, utf8ToString(msg.Content))
	}
	deleteMessage("338")
}

func utf8ToString(str string) string {
	var buf string
	for i := 0; i < len(str) - 1; i += 4 {
		if i, err := strconv.ParseInt(str[i:i+4], 16, 0); err != nil {
			fmt.Println(err)
		} else {
			buf += fmt.Sprintf("%c", i)
		}
	}
	return buf
}

func deleteMessage(id string) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://192.168.0.1/goform/goform_set_cmd_process", nil)

	req.Header.Set("Referer", "http://192.168.0.1/index.html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Content-Length", "63")
	req.Header.Set("Host", "192.168.0.1")
	req.Header.Set("DNT", "1")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")

	req.PostForm = url.Values{}
	req.PostForm.Add("isTest", "false")
	req.PostForm.Add("goformId", "DELETE_SMS")
	req.PostForm.Add("msg_id", id)
	req.PostForm.Add("notCallback", "true")
	fmt.Println(req.PostForm)
	if resp, err := client.Do(req); err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
	}
}