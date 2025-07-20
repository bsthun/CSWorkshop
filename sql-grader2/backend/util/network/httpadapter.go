package network

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

type ResponseWriter struct {
	StatusCode int
	Headers    http.Header
	Body       *bytes.Buffer
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		StatusCode: http.StatusOK,
		Headers:    make(http.Header),
		Body:       bytes.NewBuffer(nil),
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.Headers
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func FiberRequestAdapter(c *fiber.Ctx, prefix string) (*http.Request, error) {
	url := string(c.Request().URI().Path())
	url = strings.TrimPrefix(url, prefix) + "?" + string(c.Request().URI().QueryString())
	method := string(c.Request().Header.Method())
	body := bytes.NewReader(c.Body())

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})

	return req, nil
}

func FiberResponseAdapter(w *ResponseWriter, c *fiber.Ctx) error {
	for key, values := range w.Headers {
		for _, value := range values {
			c.Response().Header.Add(key, value)
		}
	}

	c.Status(w.StatusCode)
	return c.Send(w.Body.Bytes())
}

func FiberAdapter(handler http.Handler, prefix string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := FiberRequestAdapter(c, prefix)
		if err != nil {
			return err
		}

		w := NewResponseWriter()
		handler.ServeHTTP(w, req)

		return FiberResponseAdapter(w, c)
	}
}
