/*
Copyright © 2025 Bartłomiej Święcki (byo)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package assert

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/exp/constraints"
)

type TestingT interface {
	Helper()
	Error(msgAndArgs ...any)
}

func fail(t TestingT, msgAndArgs []any, fmtString string, fmtArgs ...any) {
	t.Error(
		append(
			[]any{fmt.Sprintf(fmtString, fmtArgs...)},
			msgAndArgs...,
		)...,
	)
}

func Equal[T any](t TestingT, expected, actual T, msgAndArgs ...any) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Values not equal, expected: %+v, actual: %+v", expected, actual)

	return false
}

func NotEqual[T any](t TestingT, expected, actual T, msgAndArgs ...any) bool {
	if !reflect.DeepEqual(expected, actual) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Values are equal: %v", actual)

	return false
}

func True(t TestingT, condition bool, msgAndArgs ...any) bool {
	if condition {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Condition is not true")

	return false
}

func False(t TestingT, condition bool, msgAndArgs ...any) bool {
	if !condition {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Condition is not false")

	return false
}

func Nil(t TestingT, object any, msgAndArgs ...any) bool {
	if isNil(object) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Object is not nil")

	return false
}

func NotNil(t TestingT, object any, msgAndArgs ...any) bool {
	if !isNil(object) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Object is nil")

	return false
}

func isNil(object any) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice,
		reflect.UnsafePointer:
		if value.IsNil() {
			return true
		}
	}

	return false
}

func NoError(t TestingT, err error, msgAndArgs ...any) bool {
	if err == nil {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Received unexpected error: %v", err)

	return false
}

func ErrorIs(t TestingT, err error, target error, msgAndArgs ...any) bool {
	if errors.Is(err, target) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Error is not %T: %v", target, err)

	return false
}

func ErrorContains(t TestingT, err error, contains string, msgAndArgs ...any) bool {
	if err != nil && strings.Contains(err.Error(), contains) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Error %q does not contain %q", err.Error(), contains)

	return false
}

func NotEmpty(t TestingT, object any, msgAndArgs ...any) bool {
	if !isEmpty(object) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Object is empty")

	return false
}

func isEmpty(object any) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return value.Len() == 0
	}

	return false
}

func Greater[T constraints.Ordered](t TestingT, a, b T, msgAndArgs ...any) bool {
	if a > b {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Expected %v to be greater than %v", a, b)

	return false
}

func GreaterOrEqual[T constraints.Ordered](t TestingT, a, b T, msgAndArgs ...any) bool {
	if a >= b {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Expected %v to be greater or equal than %v", a, b)

	return false
}

func isZero(value any) bool {
	return value == nil || reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}

func Zero(t TestingT, value any, msgAndArgs ...any) bool {
	if isZero(value) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Expected %v to be zero", value)

	return false
}

func NotZero(t TestingT, value any, msgAndArgs ...any) bool {
	if !isZero(value) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Expected %v to not be zero", value)

	return false
}

func didPanic(f func()) (didPanic bool) {
	didPanic = true

	defer func() { recover() }()

	f()
	didPanic = false

	return
}

func Panics(t TestingT, f func(), msgAndArgs ...any) bool {
	if didPanic(f) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Expected function to panic")

	return false
}

func NotPanics(t TestingT, f func(), msgAndArgs ...any) bool {
	if !didPanic(f) {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Expected function not to panic")

	return false
}

func Regexp(t TestingT, pattern string, text string, msgAndArgs ...any) bool {
	if matches, err := regexp.MatchString(pattern, text); err != nil {
		t.Helper()
		fail(t, msgAndArgs, "Invalid regexp pattern: %v", err)
		t.Error(msgAndArgs...)

		return false
	} else if matches {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Expected %v to match %v", text, pattern)

	return false
}

func Len(t TestingT, obj any, len int, msgAndArgs ...any) bool {
	r := reflect.ValueOf(obj)

	if r.Len() == len {
		return true
	}

	t.Helper()
	fail(t, msgAndArgs, "Expected length of %d, found: %d", len, r.Len())

	return false
}
