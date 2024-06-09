package variable

import (
	"fmt"
	"regexp"
	"strings"
)

// varRegex holds the compiled regular expression for matching variables
var varRegex = regexp.MustCompile(VariableRegexp)

// ReplaceVariables replaces all variables with the pattern {{ variable }} in the data string with their corresponding values
func (r *variableInstance) ReplaceVariables(data string) (string, error) {
	varInstances := varRegex.FindAllString(data, -1)

	errs := []error{}
	for _, v := range varInstances {
		replacement, err := r.ParseVariable(removeBraces(v))
		if err != nil {
			err = fmt.Errorf("%s: %w", v, err)
			errs = append(errs, err)

			// on error replace variable instance with blank string
			replacement = ""
		}
		data = strings.Replace(data, v, replacement, 1)
	}

	return data, combineErrors(errs)
}

func removeBraces(data string) string {
	return strings.Trim(data, "{} ")
}

func combineErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	fmtErr := errs[0]
	for i := 1; i < len(errs); i++ {
		fmtErr = fmt.Errorf("%w, %w", fmtErr, errs[i])
	}
	return fmtErr
}
