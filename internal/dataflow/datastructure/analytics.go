package datastructure

import (
	"fmt"
	"github.com/rs/zerolog"
	"helloworld/internal/utils"
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
}

func (r AnalyticsRecord) MarshalZerologObject(e *zerolog.Event) {
	e.
		Int64("timestamp", r.TimeStamp).
		Str("jobid", r.JobID).
		Str("taskid", r.TaskID).
		Str("tasktag", r.TaskTag).
		Str("taskresult", r.TaskResult)
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
	allFileds := utils.DeepFields(*a)
	fields := []string{}
	for _, valueField := range allFileds {
		var thisVal string
		fmt.Println(valueField.Kind().String() + "  " + valueField.Type().String())

		switch valueField.Type().String() {
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

// AnalyticsFilters defines records should be filtered
type AnalyticsFilters struct {
	//Usernames        []string `json:"usernames"`
	//SkippedUsernames []string `json:"skip_usernames"`
}

// ShouldFilter determine whether a record should to be filtered out.
func (filters AnalyticsFilters) ShouldFilter(record AnalyticsRecord) bool {
	//switch {
	//case len(filters.SkippedUsernames) > 0 && stringInSlice(record.Username, filters.SkippedUsernames):
	//	return true
	//case len(filters.Usernames) > 0 && !stringInSlice(record.Username, filters.Usernames):
	//	return true
	//}

	return false
}

// HasFilter determine whether a record has a filter.
func (filters AnalyticsFilters) HasFilter() bool {
	//if len(filters.SkippedUsernames) == 0 && len(filters.Usernames) == 0 {
	//	return false
	//}
	//
	//return true
	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
