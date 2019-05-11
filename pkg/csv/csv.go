package csv

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"

	"github.com/rs/zerolog/log"
)

// FieldSetter denotes the transformation from a CSV value string to field type.
type FieldSetter func(reflect.Value, string)

// CsvUnmarshaller is a custom CSV Unmarshaller reading from a CSV reader source to a struct.
type CsvUnmarshaller struct {
	*csv.Reader
	CsvToStructFieldIndex map[int]int
	FieldSetters          map[int]FieldSetter
}

// NewCsvUnmarshaller creates a CSV Unmarshaller by matching the column headers
// to struct fields, creating conversions from CSV values to matching struct
// field types.
func NewCsvUnmarshaller(r *csv.Reader, destType interface{}) *CsvUnmarshaller {
	headers := ReadHeadersStringInt(r)
	CsvToStructFieldIndex := make(map[int]int, len(headers))
	FieldSetters := make(map[int]FieldSetter, len(headers))
	t := reflect.TypeOf(destType).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if v, ok := headers[f.Name]; ok {
			CsvToStructFieldIndex[v] = i
			FieldSetters[i] = NewFieldSetter(f)
		} else if v, ok := headers[f.Tag.Get("csv")]; ok {
			CsvToStructFieldIndex[v] = i
			FieldSetters[i] = NewFieldSetter(f)
		}
	}
	return &CsvUnmarshaller{
		r,
		CsvToStructFieldIndex,
		FieldSetters,
	}
}

func NewCsvUnmarshallerFromFile(f *os.File, destType interface{}) *CsvUnmarshaller {
	return NewCsvUnmarshaller(csv.NewReader(f), destType)
}

// NewFieldSetter creates a field setter based on the struct fields type.
func NewFieldSetter(f reflect.StructField) FieldSetter {
	switch f.Type.Kind() {
	case reflect.String:
		return StringFieldSetter
	case reflect.Int:
		return IntFieldSetter
	case reflect.Bool:
		return BoolFieldSetter
	case reflect.Float64:
		return FloatFieldSetter
	}
	return nil
}

// StringFieldSetter sets CSV value to string.
func StringFieldSetter(v reflect.Value, s string) {
	v.SetString(s)
}

// IntFieldSetter sets CSV value to int.
func IntFieldSetter(v reflect.Value, s string) {
	if len(s) == 0 {
		v.SetInt(0)
		return
	}
	if i, err := strconv.ParseInt(s, 10, 64); err != nil {
		log.Error().Err(err).Str("s", s).Msg("error parsing int64")
	} else {
		v.SetInt(i)
	}
}

// BoolFieldSetter sets CSV value to bool.
func BoolFieldSetter(v reflect.Value, s string) {
	if len(s) == 0 {
		v.SetBool(false)
		return
	}
	if b, err := strconv.ParseBool(s); err != nil {
		log.Error().Err(err).Str("s", s).Msg("error parsing bool")
	} else {
		v.SetBool(b)
	}
}

// FloatFieldSetter sets CSV value to float.
func FloatFieldSetter(v reflect.Value, s string) {
	if len(s) == 0 {
		v.SetFloat(0)
		return
	}
	if f, err := strconv.ParseFloat(s, 64); err != nil {
		log.Error().Err(err).Str("s", s).Msg("error parsing float64")
	} else {
		v.SetFloat(f)
	}
}

// ReadHeadersStringInt creates matching CSV column headers to CSV index.
func ReadHeadersStringInt(r *csv.Reader) map[string]int {
	records, err := r.Read()
	if err != nil {
		log.Fatal().Err(err).Msg("error reading headers")
	}
	headers := make(map[string]int, len(records))
	for i, record := range records {
		headers[record] = i
	}
	return headers
}

// ReadHeadersIntString creates matching CSV index to CSV column headers.
func ReadHeadersIntString(r *csv.Reader) map[int]string {
	records, err := r.Read()
	if err != nil {
		log.Fatal().Err(err).Msg("error reading headers")
	}
	headers := make(map[int]string, len(records))
	for i, record := range records {
		headers[i] = record
	}
	return headers
}

// UnmarshalToStruct reads a row from CSV reader into the dest.
func (u *CsvUnmarshaller) UnmarshalToStruct(dest interface{}) error {
	records, err := u.Read()
	if err != nil {
		return err
	}
	v := reflect.ValueOf(dest).Elem()
	for csvIndex := 0; csvIndex < len(records); csvIndex++ {
		if fieldIndex, ok := u.CsvToStructFieldIndex[csvIndex]; ok {
			if setter, ok := u.FieldSetters[fieldIndex]; ok {
				setter(v.Field(fieldIndex), records[csvIndex])
			}
		}
	}
	return nil
}
