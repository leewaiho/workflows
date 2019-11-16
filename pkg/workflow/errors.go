package workflow

type ErrorHandler struct {
}

func (e *ErrorHandler) Write(p []byte) (n int, err error) {
	errMsg := string(p)
	SendItems(NewItem("运行异常", errMsg, nil, false))
	return len(p), nil
}
