package qiniu

import (
	"context"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"time"
)

type Credential struct {
	AccessKey string
	SecretKey string
}

type Client struct {
	mc         *qbox.Mac
	credential *Credential
	bucketName string
}

func NewClient(bucketName, accessKey, secretKey string) *Client {
	mc := qbox.NewMac(accessKey, secretKey)
	return &Client{
		mc:         mc,
		bucketName: bucketName,
		credential: &Credential{
			AccessKey: accessKey,
			SecretKey: secretKey,
		},
	}
}

func (c *Client) UploadFile(filename, filepath string) (*Response, error) {
	policy := storage.PutPolicy{
		Scope:      c.bucketName,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fileSize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	token := policy.UploadToken(c.mc)
	cfg := &storage.Config{}
	formUploader := storage.NewFormUploader(cfg)
	resp := new(Response)
	err := formUploader.PutFile(context.Background(), resp, token, filename, filepath, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) ListFile(prefix string, limit int) ([]storage.ListItem, error) {
	storageCfg, e := c.newStorageConfig()
	if e != nil {
		return nil, e
	}
	bucketManager := storage.NewBucketManager(c.mc, storageCfg)
	delimiter := ""
	marker := ""
	result := make([]storage.ListItem, 0)
	for {
		entries, _, nextMarker, hashNext, err := bucketManager.ListFiles(c.bucketName, prefix, delimiter, marker, limit)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			result = append(result, entry)
		}
		if !hashNext {
			break
		}
		marker = nextMarker
	}
	return result, nil
}

func (c *Client) newStorageConfig() (*storage.Config, error) {
	region, err := storage.GetRegion(c.credential.AccessKey, c.bucketName)
	if err != nil {
		return nil, err
	}
	storageCfg := &storage.Config{
		Region: region,
	}
	return storageCfg, nil
}

func FormatPutTime(putTime int64) time.Time {
	return time.Unix(0, putTime*100)
}
