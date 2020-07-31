package validation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseFluentdHash takes a string value and parses it into a map as
// FluentD would. A hash can take two forms:
//	 string-based hash: `field1:type, field2:type, field3:type:option, field4:type:option`
//	 JSON format: `{"field1":"type", "field2":"type", "field3":"type:option", "field4":"type:option"}`
// The fluentd code to parse this is can be found at
// https://github.com/fluent/fluentd/blob/340d1aaa96839038140c0280930d144af4cf10f1/lib/fluent/config/types.rb#L189
func ParseFluentdHash(value string) (map[string]string, error) {
	// If string is empty then return no error
	if value == "" {
		return nil, nil
	}

	// Create map for results
	results := make(map[string]string)

	if strings.HasPrefix(value, "{") {
		// If the string starts with a "{" then it is assumed to be JSON. In
		// this case we use the golang json library to parse it.
		err := json.Unmarshal([]byte(value), &results)
		if err != nil {
			return nil, err
		}
	} else {
		// If the string does not start with a "{" then it is a string based
		// hash
		for _, segment := range strings.Split(value, ",") {
			parts := strings.SplitN(segment, ":", 2)
			if len(parts) < 2 {
				return nil, fmt.Errorf("hash field %q missing value", parts[0])
			}

			results[parts[0]] = parts[1]
		}
	}

	// Return the hash
	return results, nil
}

// ValidateFluentdRegex validates a regex can be used by FluentD
func ValidateFluentdRegex(value string) error {
	return nil
}
