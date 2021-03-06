package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SMS struct {
	ID                   string `json:"id"`
	Number               string `json:"number"`
	Content              string `json:"content"`
	Tag                  string `json:"tag"`
	Date                 string `json:"date"`
	DraftGroupID         string `json:"draft_group_id,omitempty"`
	ReceivedAllConcatSMS string `json:"received_all_concat_sms,omitempty"`
	ConcatSMSTotal       string `json:"concat_sms_total"`
	ConcatSMSReceived    string `json:"concat_sms_received"`
	SMSClass             string `json:"sms_class"`
	SMSMem               string `json:"sms_mem"`
	SMSSubmitMsgRef      string `json:"sms_submit_msg_ref,omitempty"`
}

func getMessages() []SMS {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://192.168.0.1/goform/goform_get_cmd_process?"+
		"cmd=sms_data_total&page=0&data_per_page=500&mem_store=1&tags=10&order_by=order+by+id+asc", nil)
	req.Header.Set("Referer", "http://192.168.0.1/index.html")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	body = body[12 : len(body)-1] // Cut the {"messages": ... } part, keeping the [] of messages

	var sms []SMS

	if len(body) > 2 {
		if err := json.Unmarshal(body, &sms); err != nil {
			fmt.Println("Error:", err)
		}
	}

	return sms
}

func utf8ToString(str string) string {
	var buf string
	for i := 0; i < len(str)-1; i += 4 {
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

	data := url.Values{}
	data.Set("goformId", "DELETE_SMS")
	data.Set("msg_id", id)
	data.Set("notCallback", "true")

	req, _ := http.NewRequest("POST", "http://192.168.0.1/goform/goform_set_cmd_process",
		strings.NewReader(data.Encode()))

	req.Header.Set("Referer", "http://192.168.0.1/index.html")

	if resp, err := client.Do(req); err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
	}
}

// Get new modem messages and send them to channel
func modemHandler(updates chan string) {
	for {
		msgs := getMessages()
		if len(msgs) != 0 {
			for _, msg := range msgs {
				if msg.ReceivedAllConcatSMS == "1" {
					updates <- "[" + msg.Number + "]" + "\n\n" + utf8ToString(msg.Content)
					deleteMessage(msg.ID)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
