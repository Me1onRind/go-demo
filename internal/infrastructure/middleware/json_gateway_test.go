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
		ExpectCode int32
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

	for _, tt := range []struct {
		Raw                string
		Method             string
		Scenario           string
		Protocol           any
		ContentType        string
		ExpectValue        any
		ExpectEmptyErr     bool
		OriginZeroProtocol any
	}{
		{
			Raw:         `{"name":"l1","value":123,"slice":[{"name":"l2","value":321}]}`,
			Method:      "POST",
			ContentType: "application/json",
			Scenario:    "post_json",
			Protocol:    &P{},
			ExpectValue: &P{Name: "l1", Value: 123, Slice: []*P{
				{Name: "l2", Value: 321},
			}},
			OriginZeroProtocol: &P{},
		},
		{
			Raw:                `name=l1&value=123&slice_int=1&slice_int=2`,
			Method:             "GET",
			ContentType:        "application/json",
			Scenario:           "get",
			Protocol:           &P{},
			ExpectValue:        &P{Name: "l1", Value: 123, SliceInt: []int{1, 2}},
			OriginZeroProtocol: &P{},
		},
	} {
		t.Run(fmt.Sprintf("%s", tt.Scenario), func(t *testing.T) {
			c := &gin.Context{
				Request: &http.Request{
					Method: tt.Method,
					Header: http.Header{
						"Content-Type": []string{tt.ContentType},
					},
					Body: io.NopCloser(strings.NewReader(tt.Raw)),
					URL: &url.URL{
						RawQuery: tt.Raw,
					},
				},
			}
			newValue, err := initProtocol(c, tt.Protocol)
			assert.Empty(t, err)
			assert.Equal(t, true, assert.ObjectsAreEqualValues(tt.ExpectValue, newValue))
			assert.Equal(t, true, assert.ObjectsAreEqualValues(tt.Protocol, tt.OriginZeroProtocol))
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
