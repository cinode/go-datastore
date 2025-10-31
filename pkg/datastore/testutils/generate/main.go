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

package main

//go:generate go run .
import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/cinode/go-common/blob"
	"github.com/cinode/go-common/cutl"
	"github.com/cinode/go-datastore/pkg/blobtypes"
	"github.com/cinode/go-datastore/pkg/internal/blobtypes/dynamiclink"
)

type blobData struct {
	Name     string
	Data     string
	Expected string
}

func static(data string) blobData {
	content := []byte(data)
	hash := sha256.Sum256(content)

	n, err := blob.NameFromHashAndType(hash[:], blobtypes.Static)
	cutl.PanicIfError(err)

	return blobData{
		Name:     n.String(),
		Data:     hex.EncodeToString(content),
		Expected: hex.EncodeToString(content),
	}
}

func dynamicLink(data string, version uint64, seed int) blobData {
	baseSeed := []byte{
		byte(seed >> 8), byte(seed),
		// That's some random byte sequence, fixed here to ensure we're regenerating the same dataset
		0x7b, 0x7c, 0x99, 0x35, 0x7f, 0xed, 0x93, 0xd3,
	}

	pseudoRandBuffer := []byte{}
	for i := byte(0); i < 5; i++ {
		h := sha256.Sum256(append([]byte{i}, baseSeed...))
		pseudoRandBuffer = append(pseudoRandBuffer, h[:]...)
	}

	dl, err := dynamiclink.Create(bytes.NewReader(pseudoRandBuffer))
	cutl.PanicIfError(err)

	pr, _, err := dl.UpdateLinkData(bytes.NewBufferString(data), version)
	cutl.PanicIfError(err)

	buf, err := io.ReadAll(pr.GetPublicDataReader())
	cutl.PanicIfError(err)

	pr, err = dynamiclink.FromPublicData(dl.BlobName(), bytes.NewReader(buf))
	cutl.PanicIfError(err)

	elink, err := io.ReadAll(pr.GetEncryptedLinkReader())
	cutl.PanicIfError(err)

	return blobData{
		Name:     dl.BlobName().String(),
		Data:     hex.EncodeToString(buf),
		Expected: hex.EncodeToString(elink),
	}
}

//go:embed tesblobs.go.tpl
var templateString string
var tmpl = cutl.Must(template.New("testblobs").Parse(templateString))

func main() {
	fl := cutl.Must(os.Create("../testblobs.go"))
	defer fl.Close()

	err := tmpl.Execute(fl, map[string]any{
		"TestBlobs": []blobData{
			static("Test"),
			static("Test1"),
			static(""),
			dynamicLink("Test", 0, 0),
			dynamicLink("Test1", 1, 1),
			dynamicLink("", 2, 2),
		},
		"DynamicLinkPropagationData": []blobData{
			dynamicLink("Test1", 10000, 999),
			dynamicLink("Test2", 20000, 999),
			dynamicLink("Test3", 20000, 999),
		},
	})
	cutl.PanicIfError(err)

	fmt.Println("Successfully generated testblobs.go")
}
