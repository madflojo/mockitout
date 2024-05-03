package variable

import (
	"log"
	"regexp"
	"strings"
)

var varRegex *regexp.Regexp

func init() {
	varRegex = regexp.MustCompile(VariableRegexp)
}

func (r *RequestContext) ReplaceVariables(data string) (string, error) {
	varInstances := varRegex.FindAllString(data, -1)

	for _, v := range varInstances {
		replacement, err := r.ParseVariable(removeBraces(v))
		if err != nil {
			log.Printf("Error parsing variable %s: %s", v, err)
			continue
		}
		data = strings.Replace(data, v, replacement, 1)
	}

	return data, nil
}

func removeBraces(data string) string {
	return strings.Trim(data, "{} ")
}
