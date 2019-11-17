package cmd

import (
	"github.com/LeeWaiHo/workflows/pkg/otp"
	"github.com/LeeWaiHo/workflows/pkg/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(otpCommand)
}

var otpCommand = &cobra.Command{
	Use:     "otp",
	Short:   "One Time Passcode 快速生成工具",
	Version: "v1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		keys, e := readKeys()
		if e != nil {
			workflow.HandleError(e, "读取otp配置失败")
			return
		}
		keyName := ""
		if len(args) > 0 {
			keyName = args[0]
		}
		workflow.SendItems(otp.FilterKeys(keys, keyName)...)
	},
}

func readKeys() (map[string]*otp.Key, error) {
	keyMap := make(map[string]*otp.Key)
	err := viper.UnmarshalKey("otp", &keyMap)
	if err != nil {
		return nil, err
	}
	return keyMap, nil
}
