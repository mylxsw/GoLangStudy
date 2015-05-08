package echo

import (
	"encoding/json"
	"net/http"
)

type store map[string]interface{}

type Context struct {
	Request  *http.Request
	Response *response
	pnames   []string
	pvalues  []string
	store    store
	echo     *Echo
}
