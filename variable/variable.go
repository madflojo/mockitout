package variable

import (
	"errors"
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
)

const (
	VariableRegexp = "{{}}" // TODO: Add regex for variable
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
