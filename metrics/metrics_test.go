package metrics

// utilities for metrics tests
import "regexp"

func ContainsKeyRegexp(values Values, reg string) bool {
	for name, _ := range values {
		if matches := regexp.MustCompile(reg).FindStringSubmatch(name); matches != nil {
			return true
		}
	}
	return false
}
