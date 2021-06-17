package service

type HTTPRequest struct {
	Path string
	Method string
}

type Service struct {
	Name string
	HTTP HTTPRequest
}
