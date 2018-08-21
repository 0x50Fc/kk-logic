package logic

type View struct {
	Content     []byte
	ContentType string
	Headers     map[string]string
}
