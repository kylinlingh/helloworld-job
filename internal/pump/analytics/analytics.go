package analytics

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AnalyticsRecord struct {
	TimeStamp  int64 `json:"timestamp"`
	JobID      string
	TaskID     string
	TaskTag    string
	TaskResult string
	ExpireAt   time.Time `json:"expireAt"   bson:"expireAt"`
}

func (r *AnalyticsRecord) GetFieldNames() []string {
	val := reflect.ValueOf(r).Elem()
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		fields = append(fields, typeField.Name)
	}

	return fields
}

// GetLineValues returns all the line values.
func (a *AnalyticsRecord) GetLineValues() []string {
	val := reflect.ValueOf(a).Elem()
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		var thisVal string
		switch typeField.Type.String() {
		case "int":
			thisVal = strconv.Itoa(int(valueField.Int()))
		case "int64":
			thisVal = strconv.Itoa(int(valueField.Int()))
		case "[]string":
			tmpVal, _ := valueField.Interface().([]string)
			thisVal = strings.Join(tmpVal, ";")
		case "time.Time":
			tmpVal, _ := valueField.Interface().(time.Time)
			thisVal = tmpVal.String()
		case "time.Month":
			tmpVal, _ := valueField.Interface().(time.Month)
			thisVal = tmpVal.String()
		default:
			thisVal = valueField.String()
		}

		fields = append(fields, thisVal)
	}

	return fields
}

type AnalyticsFilters struct {
}
