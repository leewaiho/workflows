package otp

import (
	"github.com/LeeWaiHo/workflows/pkg/workflow"
	"github.com/deanishe/awgo"
	"github.com/pquerna/otp/totp"
	"strings"
	"time"
)

type Key struct {
	Hint   string `json:"hint"`
	Secret string `json:"secret"`
}

func FilterKeys(keys map[string]*Key, name string) []*aw.Item {
	items := make([]*aw.Item, 0)
	if name == "" {
		for k, v := range keys {
			items = append(items, newKeyItem(k, v))
		}
	} else {
		for k, v := range keys {
			if strings.HasPrefix(k, name) {
				items = append(items, newKeyItem(k, v))
			}
		}
	}
	return items
}

func newKeyItem(name string, key *Key) *aw.Item {
	passcode, _ := totp.GenerateCode(key.Secret, time.Now())
	return workflow.NewItem(name, key.Hint, map[string]string{
		"passcode": passcode,
	}, true)
}
