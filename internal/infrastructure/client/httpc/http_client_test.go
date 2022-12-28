package httpc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	port := 15243
	host := fmt.Sprintf("localhost:%d", port)
	http.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			fmt.Printf("http header:%s, value:%s\n", k, v)
		}
		writer.Write([]byte(`{"code":1000,"message":"msg"}`))
		//io.WriteString(writer, `{"code":1000,"message":"msg"}`)
	})

	go http.ListenAndServe(host, nil)
	time.Sleep(time.Millisecond * 20)
	os.Exit(m.Run())
}

func Test_Send_Json_Request(t *testing.T) {
	client := NewHttpClient()
	req := struct {
		Param string `json:"param"`
	}{
		Param: "p",
	}
	resp := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{}
	ctx := context.Background()
	err := client.SendJsonRequest(ctx, "http://localhost:15243/json", &req, &resp)
	assert.Empty(t, err)
}
