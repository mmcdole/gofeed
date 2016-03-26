package ext

import ()

type Extensions map[string]map[string][]Extension

type Extension struct {
	Name     string                 `json:"name"`
	Value    string                 `json:"value"`
	Attrs    map[string]string      `json:"attrs"`
	Children map[string][]Extension `json:"children"`
}

// Parse a single text value from a given extension key
func ParseTextExtension(name string, extensions map[string][]Extension) (value string) {
	if extensions == nil {
		return
	}

	matches, ok := extensions[name]
	if !ok || len(matches) == 0 {
		return
	}

	match := matches[0]
	return match.Value
}

func ParseTextArrayExtension(name string, extensions map[string][]Extension) (values []string) {
	if extensions == nil {
		return
	}

	matches, ok := extensions[name]
	if !ok || len(matches) == 0 {
		return
	}

	values = []string{}
	for _, m := range matches {
		values = append(values, m.Value)
	}
	return
}
