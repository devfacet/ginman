// ginman
// For the full copyright and license information, please view the LICENSE.txt file.

package ginman_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devfacet/ginman"
	"github.com/gin-contrib/cors"
)

func TestNewWithOptions(t *testing.T) {
	cc := cors.DefaultConfig()
	cc.AllowOrigins = []string{"http://localhost:8080"}

	table := []struct {
		options ginman.Options
		err     error
	}{
		{
			options: ginman.Options{Mode: "test"},
			err:     nil,
		},
		{
			options: ginman.Options{
				ContextMetadata:   map[string]any{"env": "dev"},
				CORSConfig:        cc,
				EnableCompression: true,
				EnableLocation:    true,
				EnableRecovery:    true,
				EnableRequestID:   true,
				Mode:              "test",
				Validations:       []string{"duration", "json", "base64Any"},
			},
			err: nil,
		},
		{
			options: ginman.Options{
				Mode:       "test",
				CORSConfig: cors.DefaultConfig(),
			},
			err: errors.New("conflict settings: all origins disabled"),
		},
	}
	for _, v := range table {
		r, err := ginman.NewWithOptions(v.options)
		if err != v.err {
			if err.Error() != v.err.Error() {
				t.Errorf("got '%v', want %v", err, v.err)
			}
		} else if r == nil {
			t.Error("got nil, want not nil")
		}
	}
}

func testRequest(r http.Handler, method, path string, body any) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	b, err = io.ReadAll(w.Body)
	if err != nil {
		return nil
	}
	return b
}
