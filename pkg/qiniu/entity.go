package qiniu

type Response struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	FileSize int    `json:"fileSize"`
	Bucket   string `json:"bucket"`
	Name     string `json:"name"`
}
