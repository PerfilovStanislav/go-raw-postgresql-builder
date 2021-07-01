package ps

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var escaper = strings.NewReplacer("'", "''")

type Sql struct {
	Query string
	Data  interface{}
}

func (s Sql) String() string {
	query := s.Query
	data := getValue(reflect.ValueOf(s.Data))

	re := regexp.MustCompile(`\$\w+`)
	keys := re.FindAllString(query, -1)
	indexes := re.FindAllStringIndex(query, -1)

	if data.Kind() == reflect.Slice {
		var interpolatedStrings []string
		cnt := data.Len()
		for k := 0; k < cnt; k++ {
			interpolated := query
			obj := data.Index(k)
			for i := len(keys) - 1; i >= 0; i-- {
				val := obj.FieldByName(keys[i][1:])
				if val.IsValid() {
					interpolated = interpolated[:indexes[i][0]] + toString(val) + interpolated[indexes[i][1]:]
				}
			}
			interpolatedStrings = append(interpolatedStrings, interpolated)
		}
		return strings.Join(interpolatedStrings, ",")
	} else {
		for i := len(keys) - 1; i >= 0; i-- {
			field := reflect.ValueOf(s.Data).FieldByName(keys[i][1:])
			if field.IsValid() {
				query = query[:indexes[i][0]] + toString(field) + query[indexes[i][1]:]
			}
		}
	}

	return query
}

func getValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func toString(v reflect.Value) string {
	v = getValue(v)

	switch v.Kind() {
	case reflect.Invalid:
		return "NULL"
	case reflect.Bool:
		if v.Bool() {
			return "TRUE"
		}
		return "FALSE"
	case reflect.String:
		return escape(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Slice:
		var s []string
		cnt := v.Len()
		for i := 0; i < cnt; i++ {
			s = append(s, toString(v.Index(i)))
		}
		return strings.Join(s, ",")
	case reflect.Struct:
		switch v.Interface().(type) {
		case Sql:
			return fmt.Sprintf("%s", v)
		default:
			e, _ := json.Marshal(v.Interface())
			return escape(string(e))
		}
	case reflect.Map:
		e, _ := json.Marshal(v.Interface())
		return escape(string(e))
	}

	return ""
}

func escape(str string) string {
	return "'" + escaper.Replace(str) + "'"
}
