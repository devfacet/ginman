// ginman
// For the full copyright and license information, please view the LICENSE.txt file.

package ginman_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/devfacet/ginman"
	"github.com/gin-gonic/gin"
)

func TestValidations(t *testing.T) {
	r, err := ginman.NewWithOptions(ginman.Options{
		Mode:        "test",
		Validations: []string{"duration", "json", "base64Any"},
	})
	if err != nil {
		t.Errorf("got '%v', want nil", err)
	} else if r == nil {
		t.Error("got nil, want not nil")
	}

	type Req struct {
		Foo string `json:"foo" binding:"omitempty,duration"`
		Bar any    `json:"bar" binding:"omitempty,json"`
		Baz string `json:"baz" binding:"omitempty,base64Any"`
	}
	type Res struct {
		Error string `json:"error,omitempty"`
	}

	table := []struct {
		req  any
		path string
		err  error
	}{
		{
			req:  Req{Foo: "5s"},
			path: "/duration",
			err:  nil,
		},
		{
			req:  Req{Foo: "5x"},
			path: "/duration-error",
			err:  errors.New("Key: 'Req.Foo' Error:Field validation for 'Foo' failed on the 'duration' tag"),
		},
		{
			req:  Req{Bar: `{"foo": "bar"}`},
			path: "/json-string",
			err:  nil,
		},
		{
			req:  Req{Bar: map[string]any{"foo": "bar"}},
			path: "/json-map",
			err:  nil,
		},
		{
			req:  Req{Bar: []string{"foo", "bar"}},
			path: "/json-slice",
			err:  nil,
		},
		{
			req:  Req{Bar: "foo"},
			path: "/json-error",
			err:  errors.New("Key: 'Req.Bar' Error:Field validation for 'Bar' failed on the 'json' tag"),
		},
		{
			req:  Req{Baz: "Zm9vK2Jhci1iYXo9Lw=="},
			path: "/base64-std",
			err:  nil,
		},
		{
			req:  Req{Baz: "Zm9vK2Jhci1iYXo9Lw"},
			path: "/base64-url",
			err:  nil,
		},
		{
			req:  Req{Baz: "error"},
			path: "/base64-error",
			err:  errors.New("Key: 'Req.Baz' Error:Field validation for 'Baz' failed on the 'base64Any' tag"),
		},
	}
	for _, v := range table {
		r.Handle(http.MethodPost, v.path, func(c *gin.Context) {
			var req Req
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, Res{})
		})
		var res Res
		b := testRequest(r, http.MethodPost, v.path, v.req)
		if err := json.Unmarshal(b, &res); err != nil {
			t.Errorf("got '%v', want nil", err)
		} else if (v.err != nil && v.err.Error() != res.Error) || (v.err == nil && res.Error != "") {
			t.Errorf("got '%v', want '%v'", res.Error, v.err)
		}
	}
}
