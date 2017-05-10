// Package schema represents linux-inspect schema.
package schema

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// Generate generates go struct text from given RawData.
func Generate(raw RawData) string {
	tagstr := "yaml"
	if !raw.IsYAML {
		tagstr = "column"
	}

	buf := new(bytes.Buffer)
	for _, col := range raw.Columns {
		goFieldName := ToField(col.Name)
		goFieldTagName := col.Name
		if !raw.IsYAML {
			goFieldTagName = ToFieldTag(goFieldTagName)
		}
		if col.Godoc != "" {
			buf.WriteString(fmt.Sprintf("\t// %s is %s.\n", goFieldName, col.Godoc))
		}
		buf.WriteString(fmt.Sprintf("\t%s\t%s\t`%s:\"%s\"`\n",
			goFieldName,
			GoType(col.Kind),
			tagstr,
			goFieldTagName,
		))

		// additional parsed column
		if v, ok := raw.ColumnsToParse[col.Name]; ok {
			switch v {
			case TypeInt64, TypeFloat64:
				// need no additional columns

			case TypeBytes:
				ntstr := "uint64"
				if col.Kind == reflect.Int64 {
					ntstr = "int64"
				}
				buf.WriteString(fmt.Sprintf("\t%sBytesN\t%s\t`%s:\"%s_bytes_n\"`\n",
					goFieldName,
					ntstr,
					tagstr,
					goFieldTagName,
				))
				buf.WriteString(fmt.Sprintf("\t%sParsedBytes\tstring\t`%s:\"%s_parsed_bytes\"`\n",
					goFieldName,
					tagstr,
					goFieldTagName,
				))

			case TypeTimeMicroseconds, TypeTimeSeconds:
				buf.WriteString(fmt.Sprintf("\t%sParsedTime\tstring\t`%s:\"%s_parsed_time\"`\n",
					goFieldName,
					tagstr,
					goFieldTagName,
				))

			case TypeIPAddress:
				buf.WriteString(fmt.Sprintf("\t%sParsedIPHost\tstring\t`%s:\"%s_parsed_ip_host\"`\n",
					goFieldName,
					tagstr,
					goFieldTagName,
				))
				buf.WriteString(fmt.Sprintf("\t%sParsedIPPort\tint64\t`%s:\"%s_parsed_ip_port\"`\n",
					goFieldName,
					tagstr,
					goFieldTagName,
				))

			case TypeStatus:
				buf.WriteString(fmt.Sprintf("\t%sParsedStatus\tstring\t`%s:\"%s_parsed_status\"`\n",
					goFieldName,
					tagstr,
					goFieldTagName,
				))

			default:
				panic(fmt.Errorf("unknown parse type %d", raw.ColumnsToParse[col.Name]))
			}
		}
	}

	return buf.String()
}

// RawDataType defines how the raw data bytes are defined.
type RawDataType int

const (
	TypeBytes RawDataType = iota
	TypeInt64
	TypeFloat64
	TypeTimeMicroseconds
	TypeTimeSeconds
	TypeIPAddress
	TypeStatus
)

// RawData defines 'proc' raw data.
type RawData struct {
	// IsYAML is true if raw data is parsable in YAML.
	IsYAML bool

	Columns        []Column
	ColumnsToParse map[string]RawDataType
}

// Column represents the schema column.
type Column struct {
	Name  string
	Godoc string
	Kind  reflect.Kind
}

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

// ToFieldTag converts raw key to field name.
func ToFieldTag(s string) string {
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
	case reflect.Int:
		return "int"
	case reflect.Int64:
		return "int64"
	case reflect.String:
		return "string"
	default:
		panic(fmt.Errorf("unknown type %q", tp.String()))
	}
}
