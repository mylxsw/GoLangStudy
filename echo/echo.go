package echo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

type Middleware interface{}
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Handler interface{}
type HTTPError struct {
	Code    int
	Message string
	Error   error
}

type HandlerFunc func(*Context) *HTTPError
type HTTPErrorHandler func(*HTTPError, *Context)
type BindFunc func(*http.Request, interface{}) *HTTPError
type Renderer interface {
	Render(w io.Writer, name string, data interface{}) *HTTPError
}

type Echo struct {
	Router           *router
	prefix           string
	middleware       []MiddlewareFunc
	maxParam         byte
	notFoundHandler  HandlerFunc
	httpErrorHandler HTTPErrorHandler
	binder           BindFunc
	renderer         Renderer
	uris             map[Handler]string
	pool             sync.Pool
}

const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"

	MIMEJSON          = "application/json"
	MIMEText          = "text/plain"
	MIMEHtml          = "text/html"
	MIMEForm          = "application/x-www-form-urlencoded"
	MIMEMultipartForm = "multipart/form-data"

	HeaderAccept             = "Accept"
	HeaderContentDisposition = "Content-Disposition"
	HeaderContentLength      = "Content-Length"
	HeaderContentType        = "Content-Type"
)

var methods = [...]string{
	CONNECT,
	DELETE,
	GET,
	POST,
	OPTIONS,
	HEAD,
	PATCH,
	PUT,
	TRACE,
}

var UnsupportedMediaType = errors.New("echo: unsupported media type")
var RendererNotRegistered = errors.New("echo: renderer not registered")

func New() (e *Echo) {
	e = &Echo{
		uris: make(map[Handler]string),
	}

	e.Router = NewRouter(e)
	e.pool.New = func() interface{} {
		return &Context{
			Response: &response{},
			pnames:   make([]string, e.maxParam),
			pvalues:  make([]string, e.maxParam),
			store:    make(store),
		}
	}

	e.MaxParam(5)
	e.NotFoundHandler(func(c *Context) {
		http.Error(c.Response, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	e.HttpErrorHandler(func(he *HTTPError, c *Context) {
		if he.Code == 0 {
			he.Code = http.StatusInternalServerError
		}

		if he.Message == "" {
			if he.Error != nil {
				he.Message = he.Error.Error()
			} else {
				he.Message = http.StatusText(he.Code)
			}
		}

		http.Error(c.Response, he.Message, he.Code)

	})

	e.Binder(func(r *http.Request, v interface{}) *HTTPError {
		ct := r.Header.Get(HeaderContentType)
		err := UnsupportedMediaType

		if strings.HasPrefix(ct, MIMEJSON) {
			err = json.NewDecoder(r.Body).Decode(v)
		} else if strings.HasPrefix(ct, MIMEForm) {
			err = nil
		}

		if err != nil {
			return &HTTPError{Error: err}
		}
		return nil
	})

	return
}

func (e *Echo) Group(pfx string, m ...Middleware) *Echo {
	g := *e
	g.prefix = g.prefix + pfx
	if len(m) > 0 {
		g.middleware = nil
		g.Use(m...)
	}
	return &g
}

func (e *Echo) Use(m ...Middleware) {
	for _, h := range m {
		e.middleware = append(e.middleware, wrapM(h))
	}
}

func (e *Echo) MaxParam(n uint8) {
	e.maxParam = n
}

func (e *Echo) NotFoundHandler(h Handler) {
	e.notFoundHandler = wrapH(h)
}

func (e *Echo) HttpErrorHandler(h HTTPErrorHandler) {
	e.httpErrorHandler = h
}

func (e *Echo) Binder(b BindFunc) {
	e.binder = b
}

func (e *Echo) Renderer(r Renderer) {
	e.renderer = r
}

func (e *Echo) Connect(path string, h Handler) {
	e.add(CONNECT, path, h)
}

func (e *Echo) add(method, path string, h Handler) {
	key := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	e.uris[key] = path
	e.Router.Add(method, e.prefix+path, wrapH(h), e)
}

func wrapM(m Middleware) MiddlewareFunc {
	switch m := m.(type) {
	case func(*Context):
		return func(h HandlerFunc) HandlerFunc {
			return func(c *Context) *HTTPError {
				m(c)
				if !c.Response.committed {
					h(c)
				}
				return nil
			}
		}
	case func(*Context) *HTTPError:
		return func(h HandlerFunc) HandlerFunc {
			return func(c *Context) *HTTPError {
				if he := m(c); he != nil {
					return he
				}
				return h(c)
			}
		}
	case func(HandlerFunc) HandlerFunc:
		return m
	case func(http.Handler) http.Handler:
		return func(h HandlerFunc) HandlerFunc {
			return func(c *Context) (he *HTTPError) {
				m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Response.Writer = w
					c.Request = r
					he = h(c)
				})).ServeHTTP(c.Response.Writer, c.Request)
				return
			}
		}
	case http.Handler, http.HandlerFunc:
		return func(h HandlerFunc) HandlerFunc {
			return func(c *Context) *HTTPError {
				m.(http.Handler).ServeHTTP(c.Response.Writer, c.Request)
				return h(c)
			}
		}
	case func(http.ResponseWriter, *http.Request):
		return func(h HandlerFunc) HandlerFunc {
			return func(c *Context) *HTTPError {
				m(c.Response, c.Request)
				if !c.Response.committed {
					h(c)
				}
				return nil
			}
		}
	case func(http.ResponseWriter, *http.Request) *HTTPError:
		return func(h HandlerFunc) HandlerFunc {
			return func(c *Context) *HTTPError {
				if he := m(c.Response, c.Request); he != nil {
					return he
				}
				if !c.Response.committed {
					h(c)
				}
				return nil
			}
		}
	default:
		panic("echo: unknown middleware")
	}
}

func wrapH(h Handler) HandlerFunc {
	switch h := h.(type) {
	case HandlerFunc:
		return h
	case func(*Context) *HTTPError:
		return h
	case func(*Context):
		return func(c *Context) *HTTPError {
			h(c)
			return nil
		}
	case http.Handler, http.HandlerFunc:
		return func(c *Context) *HTTPError {
			h.(http.Handler).ServeHTTP(c.Response, c.Request)
			return nil
		}
	case func(http.ResponseWriter, *http.Request):
		return func(c *Context) *HTTPError {
			h(c.Response, c.Request)
			return nil
		}
	case func(http.ResponseWriter, *http.Request) *HTTPError:
		return func(c *Context) *HTTPError {
			return h(c.Response, c.Request)
		}
	default:
		panic("echo: unknown handler")
	}

}
