package qiniu

import (
	"encoding/json"
	"fmt"
	"github.com/LeeWaiHo/workflows/utils"
	"testing"
	"time"
)

func newTestClient() *Client {
	return NewClient("leewaiho", utils.GetEnvOrDie("QINIU_ACCESS_TOKEN"), utils.GetEnvOrDie("QINIU_SECRET_TOKEN"))
}

func TestUploadFile(t *testing.T) {
	resp, e := newTestClient().UploadFile("99B847EF-DF8B-4F8E-92C9-280CB40697BD.png", "/Users/leewaiho/Downloads/99B847EF-DF8B-4F8E-92C9-280CB40697BD.png")
	if e != nil {
		t.Log(e)
	}
	bs, _ := json.Marshal(resp)
	t.Log(string(bs))
}

func TestListFiles(t *testing.T) {
	items, e := newTestClient().ListFile("", 5)
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
