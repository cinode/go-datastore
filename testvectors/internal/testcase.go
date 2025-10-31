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

package internal

type TestCase struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Details          string   `json:"-"`
	DetailsLines     []string `json:"details,omitempty"`
	WhenAdded        string   `json:"added_at"`
	BlobName         []byte   `json:"blob_name"`
	EncryptionKey    []byte   `json:"encryption_key"`
	UpdateDataset    []byte   `json:"update_dataset"`
	DecryptedDataset []byte   `json:"decrypted_dataset"`
	ValidPublicly    bool     `json:"valid_publicly"`
	ValidPrivately   bool     `json:"valid_privately"`
	GoErrorContains  string   `json:"go_error_contains,omitempty"`
}
