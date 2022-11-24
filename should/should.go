package should

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"
)

type assertion func(actual any, expected ...any) error

func So(t t, actual any, assertion assertion, expected ...any) {
	if t == nil {
		t = Fmt{}
	}
	if err := assertion(actual, expected...); err != nil {
		t.Helper()
		t.Error(err)
	}
}

type t interface {
	Helper()
	Error(...any)
}

type Log struct{}
type Fmt struct{}

func (Log) Helper() {}
func (Fmt) Helper() {}

func (Log) Error(a ...any) { log.Println(a...) }
func (Fmt) Error(a ...any) {
	_, file, line, _ := runtime.Caller(2)
	fmt.Printf("%s:%d", file, line)
	fmt.Println(a...)
}

var NOT negated

type negated struct{}

func Equal(actual any, expected ...any) error {
	if equalTimes(actual, expected[0]) {
		return nil
	}
	if reflect.DeepEqual(actual, expected[0]) {
		return nil
	}
	return fmt.Errorf("\nExpected: %s\nActual:   %s", format(expected[0]), format(actual))
}
func (negated) Equal(actual any, expected ...any) error {
	if Equal(actual, expected...) != nil {
		return nil
	}
	return fmt.Errorf("\nExpected:     %s\nto not equal: %s\n(but it did)", format(expected[0]), format(actual))
}
func BeTrue(actual any, _ ...any) error          { return Equal(actual, true) }
func BeFalse(actual any, _ ...any) error         { return Equal(actual, false) }
func BeNil(actual any, _ ...any) error           { return Equal(actual, nil) }
func (negated) BeNil(actual any, _ ...any) error { return NOT.Equal(actual, nil) }

func format(v any) string {
	if isTime(v) {
		return fmt.Sprintf("%v", v)
	}
	if v == nil {
		return fmt.Sprintf("%v", v)
	}
	if t := reflect.TypeOf(v); t.Kind() <= reflect.Float64 {
		return fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("%#v", v)
}

func equalTimes(a, b any) bool {
	return isTime(a) && isTime(b) && a.(time.Time).Equal(b.(time.Time))
}
func isTime(v any) bool {
	_, ok := v.(time.Time)
	return ok
}
