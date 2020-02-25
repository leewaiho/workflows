package cmd

import (
	"github.com/LeeWaiHo/workflows/pkg/otp"
	aw "github.com/deanishe/awgo"
	"github.com/spf13/cobra"
	"log"
)

var otpCommand = &cobra.Command{
	Use:     "otp",
	Short:   "One Time PassCode 快速生成工具",
	Version: "v1.1.0",
}

func init() {
	otpCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "模糊搜索otp",
		Run:   commandListOTP,
	})
	otpCommand.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "输出指定otp",
		Run:   commandGetOTP,
	})
}

func commandListOTP(_ *cobra.Command, args []string) {
	entries, err := otp.LoadEntry()
	if err != nil {
		log.Fatal("读取otp配置失败")
		return
	}
	var name string
	if len(args) > 0 {
		name = args[0]
	}
	feedback := aw.NewFeedback()
	for _, entry := range otp.SearchEntriesByName(name, entries) {
		feedback.NewItem(entry.Name).Subtitle(entry.Desc).Arg(otp.TOTPFromSecret(entry.Secret)).Valid(true)
	}
	if err := feedback.Send(); err != nil {
		log.Fatal("发送workflow异常: ", err)
	}
	return
}

func commandGetOTP(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatalln("未指定OTP Name")
		return
	}
	entries, err := otp.LoadEntry()
	if err != nil {
		log.Fatalln("读取otp配置失败")
		return
	}
	for _, entry := range entries {
		if entry.Name == args[0] {
			aw.NewArgVars().Arg(otp.TOTPFromSecret(entry.Secret)).Send()
			return
		}
	}
	log.Fatalln("请检查name是否配置错误")
	return
}
