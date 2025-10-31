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

package require

import (
	"github.com/cinode/go-datastore/pkg/internal/picotestify/assert"
	"golang.org/x/exp/constraints"
)

type TestingT interface {
	assert.TestingT
	FailNow()
}

func Equal[T any](t TestingT, expected, actual T, msgAndArgs ...any) {
	t.Helper()
	if !assert.Equal(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

func NotEqual[T any](t TestingT, expected, actual T, msgAndArgs ...any) {
	t.Helper()
	if !assert.NotEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

func Greater[T constraints.Ordered](t TestingT, a, b T, msgAndArgs ...any) {
	t.Helper()
	if !assert.Greater(t, a, b, msgAndArgs...) {
		t.FailNow()
	}
}

func GreaterOrEqual[T constraints.Ordered](t TestingT, a, b T, msgAndArgs ...any) {
	t.Helper()
	if !assert.GreaterOrEqual(t, a, b, msgAndArgs...) {
		t.FailNow()
	}
}

func True(t TestingT, condition bool, msgAndArgs ...any) {
	t.Helper()
	if !assert.True(t, condition, msgAndArgs...) {
		t.FailNow()
	}
}

func False(t TestingT, condition bool, msgAndArgs ...any) {
	t.Helper()
	if !assert.False(t, condition, msgAndArgs...) {
		t.FailNow()
	}
}

func Nil(t TestingT, object any, msgAndArgs ...any) {
	t.Helper()
	if !assert.Nil(t, object, msgAndArgs...) {
		t.FailNow()
	}
}

func NotNil(t TestingT, object any, msgAndArgs ...any) {
	t.Helper()
	if !assert.NotNil(t, object, msgAndArgs...) {
		t.FailNow()
	}
}

func NoError(t TestingT, err error, msgAndArgs ...any) {
	t.Helper()
	if !assert.NoError(t, err, msgAndArgs...) {
		t.FailNow()
	}
}

func ErrorIs(t TestingT, err, target error, msgAndArgs ...any) {
	t.Helper()
	if !assert.ErrorIs(t, err, target, msgAndArgs...) {
		t.FailNow()
	}
}

func ErrorContains(t TestingT, err error, contains string, msgAndArgs ...any) {
	t.Helper()
	if !assert.ErrorContains(t, err, contains, msgAndArgs...) {
		t.FailNow()
	}
}

func NotEmpty(t TestingT, object any, msgAndArgs ...any) {
	t.Helper()
	if !assert.NotEmpty(t, object, msgAndArgs...) {
		t.FailNow()
	}
}

func Zero(t TestingT, value any, msgAndArgs ...any) {
	t.Helper()
	if !assert.Zero(t, value, msgAndArgs...) {
		t.FailNow()
	}
}

func NotZero(t TestingT, value any, msgAndArgs ...any) {
	t.Helper()
	if !assert.NotZero(t, value, msgAndArgs...) {
		t.FailNow()
	}
}

func Panics(t TestingT, f func(), msgAndArgs ...any) {
	t.Helper()
	if !assert.Panics(t, f, msgAndArgs...) {
		t.FailNow()
	}
}

func NotPanics(t TestingT, f func(), msgAndArgs ...any) {
	t.Helper()
	if !assert.NotPanics(t, f, msgAndArgs...) {
		t.FailNow()
	}
}

func Regexp(t TestingT, pattern string, text string, msgAndArgs ...any) {
	t.Helper()
	if !assert.Regexp(t, pattern, text, msgAndArgs...) {
		t.FailNow()
	}
}

func Len(t TestingT, obj any, len int) {
	t.Helper()
	if !assert.Len(t, obj, len) {
		t.FailNow()
	}
}
