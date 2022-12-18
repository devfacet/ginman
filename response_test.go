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

func TestResponse(t *testing.T) {
	r, err := ginman.NewWithOptions(ginman.Options{Mode: "test"})
	if err != nil {
		t.Errorf("got '%v', want nil", err)
	} else if r == nil {
		t.Error("got nil, want not nil")
	}

	type Req struct {
		Foo string `json:"foo"`
	}
	type Res struct {
		ginman.Response
		Bar string `json:"bar,omitempty"`
	}

	table := []struct {
		req  any
		res  string
		path string
		err  error
	}{
		{
			req:  Req{Foo: "bar"},
			res:  "bar",
			path: "/test-01",
			err:  nil,
		},
		{
			req:  Req{},
			res:  "",
			path: "/test-02",
			err:  errors.New("invalid Foo"),
		},
	}
	for _, v := range table {
		r.Handle(http.MethodPost, v.path, func(c *gin.Context) {
			var req Req
			if err := c.ShouldBindJSON(&req); err != nil {
				r := ginman.Response{Code: http.StatusBadRequest, Error: err}
				r.Reply(c)
				return
			}
			if req.Foo == "" {
				r := ginman.Response{Code: http.StatusBadRequest, Error: "invalid Foo"}
				r.Reply(c)
				return
			}
			r := Res{
				Response: ginman.Response{Code: http.StatusOK},
				Bar:      v.res,
			}
			r.Reply(c, r)
		})
		var res Res
		b := testRequest(r, http.MethodPost, v.path, v.req)
		if err := json.Unmarshal(b, &res); err != nil {
			t.Errorf("got '%v', want nil", err)
		} else if (v.err != nil && v.err.Error() != res.Error) || (v.err == nil && res.Error != nil) {
			t.Errorf("got '%v', want '%v'", res.Error, v.err)
		}
	}
}
