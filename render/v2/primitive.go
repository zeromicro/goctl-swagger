package v2

type rawTypeFormatter func() (typ, format string)

func rawTypeFormat(rawType string) (typ, format string) {
	t, ok := rawTypeFormats[rawType]
	if ok {
		return t[0], t[1]
	}
	return "", ""
}

var rawTypeFormats = map[string][2]string{
	"string":   {"string", ""},
	"*string":  {"string", ""},
	"int":      {"integer", "int"},
	"*int":     {"integer", "int"},
	"uint":     {"integer", "uint"},
	"*uint":    {"integer", "uint"},
	"int8":     {"integer", "int8"},
	"*int8":    {"integer", "int8"},
	"uint8":    {"integer", "uint8"},
	"*uint8":   {"integer", "uint8"},
	"int16":    {"integer", "int16"},
	"*int16":   {"integer", "int16"},
	"uint16":   {"integer", "uint16"},
	"*uint16":  {"integer", "uint16"},
	"int32":    {"integer", "int32"},
	"*int32":   {"integer", "int32"},
	"uint32":   {"integer", "uint32"},
	"*uint32":  {"integer", "uint32"},
	"uint64":   {"integer", "uint64"},
	"*uint64":  {"integer", "uint64"},
	"int64":    {"integer", "int64"},
	"*int64":   {"integer", "int64"},
	"bool":     {"boolean", ""},
	"*bool":    {"boolean", ""},
	"float32":  {"number", "float"},
	"*float32": {"number", "float"},
	"float64":  {"number", "double"},
	"*float64": {"number", "double"},
}
