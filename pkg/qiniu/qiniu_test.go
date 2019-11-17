package qiniu

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var (
	c = NewClient("leewaiho", "LkyS9lQUBfviYh-2XwUGAv-IsNAlFUO2Zes8TcZN", "8j4l80AzmM7yN13lviha8MjPtPYUE9Wa-b1vZC0k")
)

func TestUploadFile(t *testing.T) {
	resp, e := c.UploadFile("99B847EF-DF8B-4F8E-92C9-280CB40697BD.png", "/Users/leewaiho/Downloads/99B847EF-DF8B-4F8E-92C9-280CB40697BD.png")
	if e != nil {
		t.Log(e)
	}
	bs, _ := json.Marshal(resp)
	t.Log(string(bs))
}

func TestClient_ListFiles(t *testing.T) {
	items, e := c.ListFile("", 5)
	if e != nil {
		t.Error(e)
	}
	for _, v := range items {
		bs, _ := json.Marshal(v)
		fmt.Println(string(bs))
	}
}

func TestData(t *testing.T) {
	var putTime int64 = 1573925796236782600
	targetTime := time.Unix(0, putTime)
	t.Log(targetTime.Format(time.RFC3339))
}
