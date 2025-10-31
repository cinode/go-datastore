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

package testvectors

import (
	"embed"
	"encoding/json"
	"io/fs"
	"strings"

	"github.com/cinode/go-datastore/pkg/utilities/golang"
	"github.com/cinode/go-datastore/testvectors/internal"
)

//go:embed dynamic/*
var testVectorsData embed.FS

type TestCase = internal.TestCase

var parsedTestCases = func() []*TestCase {
	var testCases []*TestCase

	fs.WalkDir(testVectorsData, "dynamic", func(path string, d fs.DirEntry, err error) error {
		golang.Must(err, err)

		if d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		data := golang.Must(fs.ReadFile(testVectorsData, path))

		testCase := TestCase{}

		err = json.Unmarshal(data, &testCase)
		golang.Must(err, err)

		testCase.Details = strings.Join(testCase.DetailsLines, "\n")

		testCases = append(testCases, &testCase)

		return nil
	})

	return testCases
}()

func AllTestCases(yield func(tc *TestCase) bool) {
	for _, tc := range parsedTestCases {
		if !yield(tc) {
			return
		}
	}
}
