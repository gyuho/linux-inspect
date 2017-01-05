package schema

import (
	"fmt"
	"reflect"
	"strings"
)

// ToField converts raw YAML key to Go field name.
func ToField(s string) string {
	s = strings.Replace(s, "-", "_", -1)
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, ">", "", -1)
	cs := strings.Split(s, "_")
	var ss []string
	for _, v := range cs {
		if len(v) > 0 {
			ss = append(ss, strings.Title(v))
		}
	}
	return strings.TrimSpace(strings.Join(ss, ""))
}

// ToField converts raw key to YAML field name.
func ToYAMLField(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, "-", "_", -1)
	s = strings.Replace(s, "/", "", -1)
	return strings.Replace(s, ">", "", -1)
}

// GoType converts to Go type.
func GoType(tp reflect.Kind) string {
	switch tp {
	case reflect.Float64:
		return "float64"
	case reflect.Uint64:
		return "uint64"
	case reflect.Int64:
		return "int64"
	case reflect.String:
		return "string"
	default:
		panic(fmt.Errorf("unknown type %q", tp.String()))
	}
}
