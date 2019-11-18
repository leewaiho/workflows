package cmd

import (
	"fmt"
	"github.com/LeeWaiHo/workflows/pkg/qiniu"
	"github.com/LeeWaiHo/workflows/pkg/workflow"
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path"
	"strconv"
	"time"
)

func init() {
	rootCmd.AddCommand(storageCmd)
	storageCmd.AddCommand(storageListCmd)
	storageCmd.AddCommand(storagePutCmd)
}

var (
	storageCmd = &cobra.Command{
		Use:   "storage",
		Short: "QiNiu 对象存储快捷工具",
	}
	storagePutCmd = &cobra.Command{
		Use: "put",
		Run: func(cmd *cobra.Command, args []string) {
			wf := workflow.New()
			s := new(StorageHandler)
			client, e := s.newQiNiuClient()
			if e != nil {
				wf.Error(e, "初始七牛客户端异常")
				return
			}
			for _, v := range args {
				_, e := client.UploadFile(path.Base(v), v)
				result := v + "保存成功"
				if e != nil {
					result = v + "保存失败, " + e.Error()
				}
				wf.AddItem(workflow.NewItem(v, result, map[string]string{
					"name": v,
				}, false))
			}
			wf.Send()
		},
	}
	storageListCmd = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			prefix := ""
			if len(args) > 0 {
				prefix = args[0]
			}
			wf := workflow.New()
			s := new(StorageHandler)
			client, e := s.newQiNiuClient()
			if e != nil {
				wf.Error(e, "初始七牛客户端异常")
				return
			}
			items, e := client.ListFile(prefix, getListLimit(cmd))
			if e != nil {
				wf.Error(e, "获取资源列表失败")
				return
			}
			for _, v := range items {
				subTitle := fmt.Sprintf("创建时间:%s\t文件大小:%dB\t文件类型:%s",
					qiniu.FormatPutTime(v.PutTime).Format(time.RFC3339),
					v.Fsize, v.MimeType,
				)
				wf.AddItem(workflow.NewItem(v.Key, subTitle, map[string]string{
					"key": v.Key,
					"url": storage.MakePublicURL(viper.GetString("qiniu.domain"), v.Key),
				}, false))
			}
			wf.Send()
			return
		},
	}
)

func getListLimit(cmd *cobra.Command) int {
	limit := 5
	limitFlag := cmd.PersistentFlags().Lookup("limit")
	if limitFlag != nil {
		v, e := strconv.Atoi(limitFlag.Value.String())
		if e == nil {
			limit = v
		}
	}
	return limit
}

type StorageHandler struct {
}

func (s *StorageHandler) newQiNiuClient() (*qiniu.Client, error) {
	bucketName, e := s.getQiNiuConfig("bucketName")
	if e != nil {
		return nil, e
	}
	accessKey, e := s.getQiNiuConfig("accessKey")
	if e != nil {
		return nil, e
	}
	secretKey, e := s.getQiNiuConfig("secretKey")
	if e != nil {
		return nil, e
	}
	client := qiniu.NewClient(bucketName, accessKey, secretKey)
	return client, nil
}

func (s *StorageHandler) getQiNiuConfig(name string) (string, error) {
	const keyPrefix = "qiniu."
	v := viper.GetString(keyPrefix + name)
	if len(v) > 0 {
		return v, nil
	}
	return "", errors.New(fmt.Sprintf("key [%s.%s] is empty", keyPrefix, name))
}
