package cmd

import (
	"fmt"
	"github.com/LeeWaiHo/workflows/pkg/qiniu"
	"github.com/LeeWaiHo/workflows/pkg/workflow"
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path"
	"strconv"
	"time"
)

var (
	storageCommand = &cobra.Command{
		Use:   "storage",
		Short: "七牛对象存储快捷工具",
	}
	storagePutCmd = &cobra.Command{
		Use:   "put",
		Short: "上传文件",
		Run:   execStoragePut,
	}
	storageListCmd = &cobra.Command{
		Use:   "list",
		Short: "文件列表",
		Run:   execStorageList,
	}
	storageDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "删除文件",
		Run:   execStorageDelete,
	}
)

func init() {
	storageCommand.AddCommand(storageListCmd, storagePutCmd, storageDeleteCmd)
}

func execStorageList(cmd *cobra.Command, args []string) {
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
		item := workflow.NewItem(v.Key, subTitle, map[string]string{
			"key": v.Key,
			"url": storage.MakePublicURL(viper.GetString("qiniu.domain"), v.Key),
		}, true)
		item.Arg(v.Key)
		wf.AddItem(item)
	}
	wf.Send()
}

func execStoragePut(cmd *cobra.Command, args []string) {
	s := new(StorageHandler)
	client, e := s.newQiNiuClient()
	if e != nil {
		workflow.SendItem("初始七牛客户端异常", e.Error(), nil, false)
		return
	}
	total := len(args)
	succeed := 0
	failed := 0
	for _, v := range args {
		_, e := client.UploadFile(path.Base(v), v)
		if e != nil {
			log.Println(v + "保存失败, " + e.Error())
			failed++
		} else {
			log.Println(v + "保存成功")
			succeed++
		}
	}
	result := fmt.Sprintf("总共%d个,成功%d个,失败%d个", total, succeed, failed)
	workflow.SendItem("上传文件完成", result, map[string]string{
		"result": result,
	}, false)
}

func execStorageDelete(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		errMsg := fmt.Sprintf("期望参数:1个, 实际参数%d个", len(args))
		item := workflow.NewItem("删除文件失败", errMsg, nil, false)
		workflow.SendItems(item)
		return
	}
	s := new(StorageHandler)
	client, e := s.newQiNiuClient()
	if e != nil {
		result := fmt.Sprintf("初始七牛客户端异常:%s", e.Error())
		workflow.SendItem("初始七牛客户端异常", e.Error(), map[string]string{
			"result": result,
		}, false)
		return
	}
	filename := args[0]
	err := client.DeleteFile(filename)
	title := fmt.Sprintf("删除文件(%s)完成", filename)
	if err != nil {
		result := "删除失败," + err.Error()
		workflow.SendItem(title, result, map[string]string{
			"result": result,
		}, false)
	} else {
		result := "删除成功"
		workflow.SendItem(title, result, map[string]string{
			"result": result,
		}, false)
	}
	return
}

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
	bucketName, e := s.getQiniuConfig("bucketName")
	if e != nil {
		return nil, e
	}
	accessKey, e := s.getQiniuConfig("accessKey")
	if e != nil {
		return nil, e
	}
	secretKey, e := s.getQiniuConfig("secretKey")
	if e != nil {
		return nil, e
	}
	client := qiniu.NewClient(bucketName, accessKey, secretKey)
	return client, nil
}

func (s *StorageHandler) getQiniuConfig(name string) (string, error) {
	const keyPrefix = "qiniu."
	v := viper.GetString(keyPrefix + name)
	if len(v) > 0 {
		return v, nil
	}
	return "", errors.New(fmt.Sprintf("key [%s.%s] is empty", keyPrefix, name))
}
