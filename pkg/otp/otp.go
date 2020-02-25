package otp

import (
	"fmt"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/viper"
	"regexp"
	"time"
)

type Entry struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Secret string `json:"secret"`
}

func LoadEntry() ([]*Entry, error) {
	var result []*Entry
	return result, viper.UnmarshalKey("otp", &result)
}

func SearchEntriesByName(name string, entries []*Entry) []*Entry {
	if len(name) == 0 {
		return entries
	}
	pattern := regexp.MustCompile(fmt.Sprintf(".*?%s.*?", name))
	var result []*Entry
	for _, entry := range entries {
		if pattern.MatchString(entry.Name) {
			result = append(result, entry)
		}
	}
	return result
}

func TOTPFromSecret(secret string) string {
	passcode, _ := totp.GenerateCode(secret, time.Now())
	return passcode
}
