/*
Copyright © 2025 Bartłomiej Święcki (byo)

Licensed under the Apache License,
        Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
        software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
        either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package base58_test

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cinode/go-datastore/pkg/internal/base58"
	"github.com/cinode/go-datastore/pkg/internal/picotestify/require"
)

// Test cases from (MIT License):
//
//	https://github.com/bitcoin/bitcoin/blob/master/src/test/data/base58_encode_decode.json
//
//go:embed base58_encode_decode.json
var testJSONData string

type testCase struct {
	data     []byte
	expected string
}

func validTestCases(t require.TestingT) []testCase {
	testCases := []testCase{
		{
			[]byte{},
			"",
		},
		{
			[]byte("Hello World!"),
			"2NEpo7TZRRrLZSi2U",
		},
		{
			[]byte("The quick brown fox jumps over the lazy dog."),
			"USm3fpXnKG5EUBx2ndxBDMPVciP5hGey2Jh4NDv6gmeo1LkMeiKrLJUUBk6Z",
		},
		{
			[]byte{0x00, 0x00, 0x28, 0x7f, 0xb4, 0xcd},
			"11233QC4",
		},
	}

	jsonTests := [][]string{}
	require.NoError(t, json.Unmarshal([]byte(testJSONData), &jsonTests))

	for _, tc := range jsonTests {
		require.Len(t, tc, 2)

		data, err := hex.DecodeString(tc[0])
		require.NoError(t, err)
		expected := tc[1]

		testCases = append(testCases, testCase{
			data:     data,
			expected: expected,
		})
	}
	return testCases
}

func TestEncodeDecode(t *testing.T) {
	for _, test := range validTestCases(t) {
		t.Run(fmt.Sprintf("data=%v", test.data), func(t *testing.T) {
			encoded := base58.Encode(test.data)
			require.Equal(t, test.expected, encoded)

			decodedBack, err := base58.Decode(encoded)
			require.NoError(t, err)

			require.Equal(t, test.data, decodedBack)
		})
	}
}

func TestErrorOnInvalidDecode(t *testing.T) {
	for _, test := range []string{
		"@",
		"3SEo3LWLoPntC@",
		"@3SEo3LWLoPntC",
		"3SEo3@LWLoPntC",
	} {
		t.Run(test, func(t *testing.T) {
			decoded, err := base58.Decode(test)
			require.Nil(t, decoded)
			require.ErrorIs(t, err, base58.ErrInvalidBase58Character)
		})
	}
}

func FuzzEncodeDecode(f *testing.F) {
	for _, test := range validTestCases(f) {
		f.Add(test.data)
	}

	f.Fuzz(func(t *testing.T, a []byte) {
		if len(a) > 256 {
			// No point in testing those
			t.SkipNow()
		}

		str := base58.Encode(a)
		back, err := base58.Decode(str)
		require.NoError(t, err)
		require.Equal(t, a, back)
	})
}
