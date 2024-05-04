package variable

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidVariableFormat = errors.New("error invalid variable format")
	ErrRandomVariable        = errors.New("error random variable not found")
	ErrEnvironmentVariable   = errors.New("environment variable not found")
	ErrVariableNotFound      = errors.New("variable not found")
	ErrNoBody                = errors.New("no body found")
	ErrInvalidJsonBody       = errors.New("unable to parse json body")
	ErrInvalidJsonVar        = errors.New("unable to find json variable")
)

const (
	VariableRegexp = `\{\{[\w\-\s._$]+\}\}` // TODO: Add regex for variable
	RandomPrefix   = "$"
	HeaderPrefix   = "header."
	QueryPrefix    = "query."
	ParamPrefix    = "param."
	EnvPrefix      = "environment."
	BodyPrefix     = "body."
)

type RequestContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   httprouter.Params
}

func NewRequestContext(r *http.Request, w http.ResponseWriter, p httprouter.Params) *RequestContext {
	return &RequestContext{
		Request:  r,
		Response: w,
		Params:   p,
	}
}

func (r *RequestContext) ParseVariable(variable string) (string, error) {
	if len(variable) == 0 {
		return "", ErrVariableNotFound
	}

	cutVariable, ok := strings.CutPrefix(variable, RandomPrefix)
	if ok {
		return getRandomVariable(cutVariable)
	}

	cutVariable, ok = strings.CutPrefix(variable, HeaderPrefix)
	if ok {
		return r.getHeaderVariable(cutVariable)
	}

	cutVariable, ok = strings.CutPrefix(variable, QueryPrefix)
	if ok {
		return r.getQueryVariable(cutVariable)
	}

	cutVariable, ok = strings.CutPrefix(variable, ParamPrefix)
	if ok {
		return r.getParamVariable(cutVariable)
	}

	cutVariable, ok = strings.CutPrefix(variable, EnvPrefix)
	if ok {
		return getEnvironmentVariable(cutVariable)
	}

	cutVariable, ok = strings.CutPrefix(variable, BodyPrefix)
	if ok {
		return r.getBodyJsonVariable(cutVariable)
	}

	if strings.Compare(variable, "body") == 0 {
		return r.getTextBody(variable)
	}

	return "", ErrVariableNotFound
}

func getRandomVariable(variable string) (string, error) {
	randVariableFunction, ok := RandomMap[variable]
	if !ok {
		return "", ErrRandomVariable
	}

	return randVariableFunction(), nil
}

func (r *RequestContext) getHeaderVariable(variable string) (string, error) {
	header := r.Request.Header.Get(variable)
	if len(header) == 0 {
		return "", ErrInvalidVariableFormat
	}
	return header, nil
}

func (r *RequestContext) getQueryVariable(variable string) (string, error) {
	query := r.Request.URL.Query().Get(variable)
	if len(query) == 0 {
		return "", ErrInvalidVariableFormat
	}
	return query, nil
}

func (r *RequestContext) getParamVariable(variable string) (string, error) {
	param := r.Params.ByName(variable)
	if len(param) == 0 {
		return "", ErrInvalidVariableFormat
	}
	return param, nil
}

func getEnvironmentVariable(variable string) (string, error) {
	val := os.Getenv(variable)
	if len(val) == 0 {
		return "", ErrEnvironmentVariable
	}
	return val, nil
}

func (r *RequestContext) getTextBody(variable string) (string, error) {
	if r.Request.Body == nil || r.Request.ContentLength < 1 {
		return "", ErrNoBody
	}

	reader := io.Reader(r.Request.Body)
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (r *RequestContext) getBodyJsonVariable(variable string) (string, error) {
	if r.Request.Body == nil || r.Request.ContentLength < 1 {
		return "", ErrNoBody
	}

	jsonPath := strings.Split(variable, ".")
	if len(jsonPath) == 0 {
		return "", ErrInvalidJsonVar
	}

	var data interface{}
	err := json.NewDecoder(r.Request.Body).Decode(&data)
	if err != nil {
		return "", ErrInvalidJsonBody
	}

	for _, path := range jsonPath {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return "", ErrInvalidJsonBody
		}

		val, ok := dataMap[path]
		if !ok {
			return "", ErrInvalidJsonVar
		}

		data = val
	}

	// try return as string
	if val, ok := data.(string); ok {
		return val, nil
	}

	// try return as json (or default to string)
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("%v", data), nil
	} else {
		return string(jsonValue), nil
	}
}
