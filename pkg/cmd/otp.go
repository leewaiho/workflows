package cmd

import (
	"github.com/LeeWaiHo/workflows/pkg/workflow"
	aw "github.com/deanishe/awgo"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(otpCommand)
	otpCommand.SetErr(&workflow.ErrorHandler{})
}

type Key struct {
	Hint   string `json:"hint"`
	Secret string `json:"secret"`
}

func ReadKeys() (map[string]*Key, error) {
	keyMap := make(map[string]*Key)
	//fmt.Printf("%v", viper.Get("otp"))
	err := viper.UnmarshalKey("otp", &keyMap)
	if err != nil {
		return nil, err
	}
	return keyMap, nil
}

var otpCommand = &cobra.Command{
	Use: "otp",
	Run: func(cmd *cobra.Command, args []string) {
		keyName := ""
		if len(args) > 0 {
			keyName = args[0]
		}
		workflow.SendItems(filterByKeyName(keyName)...)
	},
}

func filterByKeyName(name string) []*aw.Item {
	items := make([]*aw.Item, 0)
	keys, e := ReadKeys()
	if e != nil {
		items = append(items, workflow.NewItem("读取配置异常", e.Error(), nil, false))
		return items
	}
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
