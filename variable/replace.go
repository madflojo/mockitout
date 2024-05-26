package variable

import (
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// varRegex holds the compiled regular expression for matching variables
var varRegex = regexp.MustCompile(VariableRegexp)

// log is used within the variable package for logging.
var log *logrus.Logger

// InitLogger takes in a logrus.Logger instance and initialises the log variable
func InitLogger(logger *logrus.Logger) {
	log = logger
}

// ReplaceVariables replaces all variables with the pattern {{ variable }} in the data string with their corresponding values
func (r *variableInstance) ReplaceVariables(data string) (string, error) {
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
