package http_proto

type GreetProxyRequest struct {
	Name string `bind:"required,min=1"`
	Msg  string `bind:"required,min=1"`
}
