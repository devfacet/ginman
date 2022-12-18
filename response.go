// ginman
// For the full copyright and license information, please view the LICENSE.txt file.

package ginman

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// Response represents an HTTP response.
type Response struct {
	Code     int           `json:"code,omitempty"`
	Status   string        `json:"status,omitempty"`
	Error    interface{}   `json:"error,omitempty"`
	Redirect string        `json:"-"`
	Sleep    time.Duration `json:"-"`
}

// Reply replies an HTTP request by the Response instance and the given trailing arguments.
func (resp *Response) Reply(c *gin.Context, responses ...interface{}) {
	// Check fields
	if resp.Code == 0 {
		resp.Code = http.StatusInternalServerError
	}
	resp.Status = http.StatusText(resp.Code)
	if resp.Error != nil {
		if err, ok := resp.Error.(error); ok {
			resp.Error = err.Error()
		}
	}

	// Handle logs
	_, file, line, _ := runtime.Caller(1)
	c.Set("log.caller", fmt.Sprintf("%s:%d", file, line))
	// Check the request id header
	if rid := c.GetHeader(xRequestIDHeaderKey); rid != "" {
		c.Set("log.rid", rid)
	}

	// Merge responses
	var obj any
	for _, v := range responses {
		if err := mapstructure.Decode(v, &obj); err != nil {
			c.Set("log.errors", append(c.GetStringSlice("log.errors"), fmt.Sprintf("couldn't merge reply responses: %s", err)))
		}
	}
	if err := mapstructure.Decode(resp, &obj); err != nil {
		c.Set("log.errors", append(c.GetStringSlice("log.errors"), fmt.Sprintf("couldn't merge reply responses: %s", err)))
	}

	// Check the sleep
	if resp.Sleep > 0 {
		time.Sleep(resp.Sleep)
	}

	// Check the redirect
	if resp.Redirect != "" {
		b, err := json.Marshal(obj)
		if err != nil {
			c.Redirect(302, resp.Redirect)
			return
		}
		c.Redirect(302, fmt.Sprintf("%s?response=%s", resp.Redirect, base64.RawURLEncoding.EncodeToString(b)))
		return
	}

	// Reply
	c.JSON(resp.Code, obj)
}
