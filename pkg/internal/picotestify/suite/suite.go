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

package suite

import (
	"reflect"
	"regexp"
	"testing"
)

type Suite struct {
	t *testing.T
}

type SuiteInterface interface {
	T() *testing.T
	SetT(t *testing.T)
}

func (s *Suite) T() *testing.T     { return s.t }
func (s *Suite) SetT(t *testing.T) { s.t = t }

var testMethodRe = regexp.MustCompile("^Test")

func listTests(suite SuiteInterface) func(yield func(string, func()) bool) {
	return func(yield func(string, func()) bool) {
		st := reflect.TypeOf(suite)

		for i := 0; i < st.NumMethod(); i++ {
			method := st.Method(i)
			if !testMethodRe.MatchString(method.Name) {
				continue
			}

			if !yield(
				method.Name,
				func() { method.Func.Call([]reflect.Value{reflect.ValueOf(suite)}) },
			) {
				return
			}
		}
	}
}

func Run(parentT *testing.T, suite SuiteInterface) {
	suite.SetT(parentT)

	if setupSuite, isSetupSuite := suite.(interface{ SetupTestSuite() }); isSetupSuite {
		setupSuite.SetupTestSuite()
	}

	if tearDownSuite, isTearDownSuite := suite.(interface{ TearDownTestSuite() }); isTearDownSuite {
		parentT.Cleanup(tearDownSuite.TearDownTestSuite)
	}

	for name, testFunc := range listTests(suite) {
		parentT.Run(name, func(childT *testing.T) {
			suite.SetT(childT)

			if setupTest, isSetupTest := suite.(interface{ SetupTest() }); isSetupTest {
				setupTest.SetupTest()
			}

			if tearDownTest, isTearDownTest := suite.(interface{ TearDownTest() }); isTearDownTest {
				childT.Cleanup(tearDownTest.TearDownTest)
			}

			testFunc()

			suite.SetT(parentT)
		})
	}
}
