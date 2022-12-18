// ginman
// For the full copyright and license information, please view the LICENSE.txt file.

// Package ginman provides a simple wrapper for Gin Web Framework.
package ginman

import (
	"errors"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	xRequestIDHeaderKey = "X-Request-ID"
)

var (
	validate = validator.New()
)

// Options represents the options which can be set when creating a new engine.
type Options struct {
	ContextMetadata   map[string]any
	CORSConfig        cors.Config
	EnableCompression bool
	EnableLocation    bool
	EnableRecovery    bool
	EnableRequestID   bool
	Mode              string
	ValidationFuncs   map[string]validator.Func
	Validations       []string
}

// NewWithOptions returns a new gin.Engine instance by the given options.
func NewWithOptions(o Options) (*gin.Engine, error) {
	// Set the server mode
	// Default value is set by GIN_MODE environment variable.
	switch o.Mode {
	case "release", "production", "prod", "stage", "staging", "dev", "development":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	}

	// Check validations
	if len(o.Validations) > 0 || len(o.ValidationFuncs) > 0 {
		validatorEngine, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return nil, errors.New("couldn't get the validator engine")
		}
		for _, tag := range o.Validations {
			if o.ValidationFuncs == nil {
				o.ValidationFuncs = make(map[string]validator.Func)
			}
			switch tag {
			case "duration":
				o.ValidationFuncs[tag] = validateDuration
			case "json":
				o.ValidationFuncs[tag] = validateJSON
			case "base64Any":
				o.ValidationFuncs[tag] = validateBase64Any
			}
		}
		for tag, fn := range o.ValidationFuncs {
			validatorEngine.RegisterValidation(tag, fn)
		}
	}

	// Init the server instance
	r := gin.New()
	if o.EnableRecovery {
		r.Use(gin.Recovery())
	}
	if o.EnableCompression {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	if o.EnableLocation {
		r.Use(location.Default())
	}
	if o.EnableRequestID {
		r.Use(func(c *gin.Context) {
			// Check the request id
			rid := c.GetHeader(xRequestIDHeaderKey)
			if rid != "" {
				// Check and replace the request id if necessary
				if _, err := uuid.Parse(rid); err != nil {
					rid = uuid.New().String()
				}
			} else {
				// Create a new request id
				rid = uuid.New().String()
				c.Request.Header.Add(xRequestIDHeaderKey, rid)
			}
			// Set the header
			c.Header(xRequestIDHeaderKey, rid)
			c.Next()
		})
	}

	// Ref: https://github.com/gin-contrib/cors/issues/29
	if checkCORSConfigChanges(o.CORSConfig) {
		if err := o.CORSConfig.Validate(); err != nil {
			return nil, err
		}
		r.Use(cors.New(o.CORSConfig))
	}

	// Ref: https://github.com/gin-gonic/gin/issues/932#issuecomment-305992716
	if len(o.ContextMetadata) > 0 {
		r.Use(func(c *gin.Context) {
			for k, v := range o.ContextMetadata {
				c.Set(k, v)
			}
			c.Next()
		})
	}

	return r, nil
}

// checkCORSConfigChanges checks the given CORS config and returns true if it has been changed.
func checkCORSConfigChanges(config cors.Config) bool {
	if config.AllowAllOrigins {
		return true
	}
	if len(config.AllowOrigins) > 0 {
		return true
	}
	if config.AllowOriginFunc != nil {
		return true
	}
	if len(config.AllowMethods) > 0 {
		return true
	}
	if len(config.AllowHeaders) > 0 {
		return true
	}
	if config.AllowCredentials {
		return true
	}
	if len(config.ExposeHeaders) > 0 {
		return true
	}
	if config.MaxAge != 0 {
		return true
	}
	if config.AllowWildcard {
		return true
	}
	if config.AllowBrowserExtensions {
		return true
	}
	if config.AllowWebSockets {
		return true
	}
	if config.AllowFiles {
		return true
	}
	return false
}
