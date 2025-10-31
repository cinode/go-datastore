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

package dynamiclink

import (
	"bytes"
	"io"
	"testing"

	"github.com/cinode/go-datastore/pkg/common"
	"github.com/cinode/go-datastore/testvectors"
	"github.com/stretchr/testify/require"
)

func TestVectors(t *testing.T) {
	for testCase := range testvectors.AllTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Run("validate public scope", func(t *testing.T) {
				err := func() error {
					bn, err := common.BlobNameFromBytes(testCase.BlobName)
					if err != nil {
						return err
					}
					pr, err := FromPublicData(
						bn,
						bytes.NewReader(testCase.UpdateDataset),
					)
					if err != nil {
						return err
					}

					dr := pr.GetEncryptedLinkReader()

					_, err = io.ReadAll(dr)
					if err != nil {
						return err
					}

					return nil
				}()

				if testCase.ValidPublicly {
					require.NoError(t, err, testCase.Details)
				} else {
					require.ErrorContains(t, err, testCase.GoErrorContains, testCase.Details)
				}
			})

			t.Run("validate private scope", func(t *testing.T) {
				err := func() error {
					bn, err := common.BlobNameFromBytes(testCase.BlobName)
					if err != nil {
						return err
					}

					pr, err := FromPublicData(
						bn,
						bytes.NewReader(testCase.UpdateDataset),
					)
					if err != nil {
						return err
					}

					dr, err := pr.GetLinkDataReader(
						common.BlobKeyFromBytes(testCase.EncryptionKey),
					)
					if err != nil {
						return err
					}

					data, err := io.ReadAll(dr)
					if err != nil {
						return err
					}

					// If we've got here - validation passed, the dataset must be correct.
					// If it is not it may indicate failure to detect attack
					require.Equal(t, testCase.DecryptedDataset, data)
					return nil
				}()

				if testCase.ValidPrivately {
					require.NoError(t, err, testCase.Details)
				} else {
					require.ErrorContains(t, err, testCase.GoErrorContains, testCase.Details)
				}
			})
		})
	}
}
