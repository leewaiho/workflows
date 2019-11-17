package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "workflows",
		Short: "alfred workflows 生产力工具",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "conf", "c", "", "配置文件路径 默认路径:$HOME/.workflows.json")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, e := homedir.Dir()
		if e != nil {
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".workflows")
	}
	if err := viper.ReadInConfig(); err == nil {
		//log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
