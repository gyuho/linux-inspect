package schema

import "strings"

// ToField converts raw YAML key to Go field name.
func ToField(s string) string {
	s = strings.Replace(s, "-", "_", -1)
	s = strings.Replace(s, "/", "", -1)
	cs := strings.Split(s, "_")
	var ss []string
	for _, v := range cs {
		if len(v) > 0 {
			ss = append(ss, strings.Title(v))
		}
	}
	return strings.TrimSpace(strings.Join(ss, ""))
}
