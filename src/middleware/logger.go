package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (m *Middleware) Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			reqBody := []byte{}
			if req.Body != nil {
				reqBody, _ = io.ReadAll(req.Body)
			}
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			if err := next(c); err != nil {
				c.Error(err)
			}

			m.logger.Infow("INBOUND LOG",
				"request_id", m.getReqID(req.Context()),
				"remote_ip", c.RealIP(),
				"host", req.Host,
				"uri", req.RequestURI,
				"method", req.Method,
				"user_agent", req.UserAgent(),
				"body", m.compactJSON(reqBody),
				"status", res.Status,
				"latency", float64(time.Since(start).Nanoseconds()/1e4)/100.0,
				"bytes_in", req.Header.Get(echo.HeaderContentLength),
				"bytes_out", strconv.FormatInt(res.Size, 10),
			)

			return nil
		}
	}
}

func (m *Middleware) getReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value("request_id").(string); ok {
		return reqID
	} else {
		return uuid.New().String()
	}
}

func (m *Middleware) compactJSON(data []byte) string {
	var js map[string]interface{}
	if json.Unmarshal(data, &js) != nil {
		return string(data)
	}

	result := new(bytes.Buffer)
	if err := json.Compact(result, data); err != nil {
		fmt.Println(err)
	}
	return result.String()
}
