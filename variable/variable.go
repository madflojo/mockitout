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
	// ErrInvalidVariableFormat is returned when the variable format is invalid (e.g. prefix present but no named variable)
	ErrInvalidVariableFormat = errors.New("error invalid variable format")

	// ErrInvalidVariablePrefix is returned when the variable prefix is invalid
	ErrInvalidVariablePrefix = errors.New("error invalid variable prefix")

	// ErrRandomVariable is returned when the random variable is not found
	ErrInvalidRandomVariable = errors.New("error random variable not found")

	// ErrEnvironmentVariable is returned when the environment variable is not found
	ErrEnvironmentVariable = errors.New("environment variable not found")

	// ErrNoBody is returned when no body is found
	ErrNoBody = errors.New("no body found")

	// ErrInvalidJsonBody is returned when the json body is unable to be parsed
	ErrInvalidJsonBody = errors.New("unable to parse json body")

	// ErrInvalidJsonVar is returned when the json variable is not found
	ErrInvalidJsonVar = errors.New("unable to find json variable")
)

const (
	// VariableRegexp is the regular expression to match variables with the format {{ variable }}
	VariableRegexp = `\{\{[\w\-\s._$]+\}\}`

	// Variable prefixes
	RandomPrefix = "$"
	HeaderPrefix = "header."
	QueryPrefix  = "query."
	ParamPrefix  = "param."
	EnvPrefix    = "environment."
	BodyPrefix   = "body."
)

// variableInstance is a struct that holds the request, response writer, and params for a request instance
type variableInstance struct {
	r *http.Request
	w http.ResponseWriter
	p httprouter.Params
}

// NewVariableInstance creates a new variable instance with the request, response writer and params
func NewVariableInstance(r *http.Request, w http.ResponseWriter, p httprouter.Params) *variableInstance {
	return &variableInstance{
		r: r,
		w: w,
		p: p,
	}
}

// ParseVariable takes the variable (with prefix e.g. $, header., ...) and returns the corresponding value or an error
func (r *variableInstance) ParseVariable(variable string) (string, error) {
	if len(variable) == 0 {
		return "", ErrInvalidVariablePrefix
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
		return r.getEnvironmentVariable(cutVariable)
	}

	cutVariable, ok = strings.CutPrefix(variable, BodyPrefix)
	if ok {
		return r.getBodyJsonVariable(cutVariable)
	}

	if strings.Compare(variable, "body") == 0 {
		return r.getTextBody(variable)
	}

	return "", ErrInvalidVariablePrefix
}

func getRandomVariable(variable string) (string, error) {
	randVariableFunction, ok := RandomMap[variable]
	if !ok {
		return "", ErrInvalidRandomVariable
	}

	return randVariableFunction(), nil
}

func (r *variableInstance) getHeaderVariable(variable string) (string, error) {
	header := r.r.Header.Get(variable)
	if len(header) == 0 {
		return "", ErrInvalidVariableFormat
	}
	return header, nil
}

func (r *variableInstance) getQueryVariable(variable string) (string, error) {
	query := r.r.URL.Query().Get(variable)
	if len(query) == 0 {
		return "", ErrInvalidVariableFormat
	}
	return query, nil
}

func (r *variableInstance) getParamVariable(variable string) (string, error) {
	param := r.p.ByName(variable)
	if len(param) == 0 {
		return "", ErrInvalidVariableFormat
	}
	return param, nil
}

func (r *variableInstance) getEnvironmentVariable(variable string) (string, error) {
	val := os.Getenv(variable)
	if len(val) == 0 {
		return "", ErrEnvironmentVariable
	}
	return val, nil
}

func (r *variableInstance) getTextBody(variable string) (string, error) {
	if r.r.Body == nil || r.r.ContentLength < 1 {
		return "", ErrNoBody
	}

	reader := io.Reader(r.r.Body)
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (r *variableInstance) getBodyJsonVariable(variable string) (string, error) {
	if r.r.Body == nil || r.r.ContentLength < 1 {
		return "", ErrNoBody
	}

	jsonPath := strings.Split(variable, ".")
	if len(jsonPath) == 0 {
		return "", ErrInvalidJsonVar
	}

	var data interface{}
	err := json.NewDecoder(r.r.Body).Decode(&data)
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
