package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/Me1onRind/go-demo/internal/model/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_GetResponse(t *testing.T) {
	for _, tt := range []struct {
		Data       any
		Err        error
		ExpectCode int
		ExpectData any
		Scenario   string
	}{
		{
			Scenario:   "success",
			Data:       map[string]any{"name": "twerwen"},
			ExpectData: map[string]any{"name": "twerwen"},
		},
		{
			Scenario:   "error",
			Data:       map[string]any{"name": "rwejrew"},
			Err:        errors.NewError(-1000, "test_error"),
			ExpectCode: -1000,
			ExpectData: nil,
		},
		{
			Scenario:   "warning",
			Data:       map[string]any{"name": "werjhwe"},
			Err:        errors.NewError(1, "test_waring_"),
			ExpectCode: 1,
			ExpectData: map[string]any{"name": "werjhwe"},
		},
		{
			Scenario:   "unexpect",
			Data:       map[string]any{"name": "werjhwe"},
			Err:        fmt.Errorf("????????"),
			ExpectCode: -100000,
		},
	} {
		t.Run(fmt.Sprintf("%s_response", tt.Scenario), func(t *testing.T) {
			r := getResponse(tt.Data, tt.Err)
			t.Log(*r, tt.Err)
			assert.Equal(t, tt.ExpectCode, r.Code)
			assert.EqualValues(t, tt.ExpectData, r.Data)
		})
	}
}

func Test_NewProtocolStruct(t *testing.T) {
	type P struct {
		Name     string `json:"name" form:"name"`
		Value    int64  `json:"value" form:"value"`
		Slice    []*P   `json:"slice" form:"slice"`
		SliceInt []int  `json:"slice_int" form:"slice_int"`
	}

	for _, test := range []struct {
		name               string
		raw                string
		method             string
		protocol           any
		contentType        string
		expectValue        any
		expectEmptyErr     bool
		originZeroProtocol any
	}{
		{
			name:        "post_json",
			raw:         `{"name":"l1","value":123,"slice":[{"name":"l2","value":321}]}`,
			method:      "POST",
			contentType: "application/json",
			protocol:    &P{},
			expectValue: &P{Name: "l1", Value: 123, Slice: []*P{
				{Name: "l2", Value: 321},
			}},
			originZeroProtocol: &P{},
		},
		{
			raw:                `name=l1&value=123&slice_int=1&slice_int=2`,
			method:             "GET",
			contentType:        "application/json",
			name:               "get",
			protocol:           &P{},
			expectValue:        &P{Name: "l1", Value: 123, SliceInt: []int{1, 2}},
			originZeroProtocol: &P{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := &gin.Context{
				Request: &http.Request{
					Method: test.method,
					Header: http.Header{
						"Content-Type": []string{test.contentType},
					},
					Body: io.NopCloser(strings.NewReader(test.raw)),
					URL: &url.URL{
						RawQuery: test.raw,
					},
				},
			}
			newValue, err := initProtocol(c, test.protocol)
			assert.Empty(t, err)
			assert.Equal(t, true, assert.ObjectsAreEqualValues(test.expectValue, newValue))
			assert.Equal(t, true, assert.ObjectsAreEqualValues(test.protocol, test.originZeroProtocol))
		})
	}

}

func Test_newProtocol(t *testing.T) {
	for _, tt := range []struct {
		Protocol     any
		ExpectSucess bool
	}{
		{
			Protocol:     1,
			ExpectSucess: false,
		},
		{
			Protocol:     "string",
			ExpectSucess: false,
		},
		{
			Protocol:     34.234,
			ExpectSucess: false,
		},
		{
			Protocol:     struct{}{},
			ExpectSucess: true,
		},
		{
			Protocol:     &struct{}{},
			ExpectSucess: true,
		},
	} {
		t.Run(reflect.TypeOf(tt.Protocol).Kind().String(), func(t *testing.T) {
			n, err := newProtocol(tt.Protocol)
			if tt.ExpectSucess {
				assert.Equal(t, true, n != nil)
				assert.Empty(t, err)
			} else {
				assert.Empty(t, n)
				assert.NotEmpty(t, err)
			}
		})
	}
}

func Test_jsonGateWay(t *testing.T) {
}
