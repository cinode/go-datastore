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

package datastoreconformancetest

import (
	"crypto/sha256"
	"fmt"
	"io"
	"testing"

	"github.com/cinode/go-datastore/pkg/blobtypes"
	"github.com/cinode/go-datastore/pkg/common"
	"github.com/cinode/go-datastore/pkg/datastore"
	"github.com/cinode/go-datastore/pkg/datastore/testutils"
	"github.com/cinode/go-datastore/pkg/internal/picotestify/require"
	"github.com/cinode/go-datastore/pkg/internal/picotestify/suite"
)

type StorageBackendTestSuite struct {
	suite.Suite

	// st is the current storage backend instance to test.
	st datastore.StorageBackend

	// CreateStorage is a function that creates a new storage backend for each test.
	CreateStorage func() (datastore.StorageBackend, error)

	// ExpectedKind is the expected kind of the storage backend
	ExpectedKind string
}

func NewStorageBackendTestSuite(
	createStorage func() (
		datastore.StorageBackend,
		error,
	),
	expectedKind string,
) *StorageBackendTestSuite {
	return &StorageBackendTestSuite{
		CreateStorage: createStorage,
		ExpectedKind:  expectedKind,
	}
}

func (s *StorageBackendTestSuite) SetupTest() {
	t := s.T()

	st, err := s.CreateStorage()
	require.NoError(t, err)
	s.st = st
}

func (s *StorageBackendTestSuite) TestStorageKind() {
	t := s.T()

	kind := s.st.Kind()
	require.Equal(t, s.ExpectedKind, kind)
}

func (s *StorageBackendTestSuite) TestStorageOpenFailureNotFound() {
	t := s.T()

	r, err := s.st.OpenReadStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.ErrorIs(t, err, datastore.ErrNotFound)
	require.Nil(t, r)
}

func (s *StorageBackendTestSuite) TestStorageSaveOpenSuccess() {
	t := s.T()

	exists, err := s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.False(t, exists)

	w, err := s.st.OpenWriteStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	exists, err = s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.False(t, exists)

	n, err := w.Write([]byte("Hello world!"))
	require.NoError(t, err)
	require.Equal(t, 12, n)

	err = w.Close()
	require.NoError(t, err)

	exists, err = s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.True(t, exists)

	r, err := s.st.OpenReadStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	b, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, []byte("Hello world!"), b)

	err = r.Close()
	require.NoError(t, err)
}

func (s *StorageBackendTestSuite) TestStorageSaveOpenCancelSuccess() {
	t := s.T()

	exists, err := s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.False(t, exists)

	w, err := s.st.OpenWriteStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	exists, err = s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.False(t, exists)

	n, err := w.Write([]byte("Hello world!"))
	require.NoError(t, err)
	require.Equal(t, 12, n)

	exists, err = s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.False(t, exists)

	w.Cancel()

	exists, err = s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.False(t, exists)

	r, err := s.st.OpenReadStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.ErrorIs(t, err, datastore.ErrNotFound)
	require.Nil(t, r)
}

func (s *StorageBackendTestSuite) TestStorageDelete() {
	t := s.T()

	blobNames := []*common.BlobName{}
	blobDatas := [][]byte{}

	t.Run("generate test data", func(t *testing.T) {
		for _, d := range []string{
			"first",
			"second",
			"third",
		} {
			h := sha256.Sum256([]byte(d))
			bn, err := common.BlobNameFromHashAndType(h[:], blobtypes.Static)
			require.NoError(t, err)

			blobNames = append(blobNames, bn)
			blobDatas = append(blobDatas, []byte(d))

			err = s.st.Delete(t.Context(), bn)
			require.ErrorIs(t, err, datastore.ErrNotFound)

			w, err := s.st.OpenWriteStream(t.Context(), bn)
			require.NoError(t, err)

			exists, err := s.st.Exists(t.Context(), bn)
			require.NoError(t, err)
			require.False(t, exists)

			n, err := w.Write([]byte(d))
			require.NoError(t, err)
			require.Equal(t, len(d), n)

			err = w.Close()
			require.NoError(t, err)

			exists, err = s.st.Exists(t.Context(), bn)
			require.NoError(t, err)
			require.True(t, exists)
		}
	})

	t.Run("delete blob", func(t *testing.T) {
		const toDelete = 1

		err := s.st.Delete(t.Context(), blobNames[toDelete])
		require.NoError(t, err)

		err = s.st.Delete(t.Context(), blobNames[toDelete])
		require.ErrorIs(t, err, datastore.ErrNotFound)

		for i := range blobNames {
			t.Run(fmt.Sprintf("exists test %d", i), func(t *testing.T) {
				exists, err := s.st.Exists(t.Context(), blobNames[i])
				require.NoError(t, err)
				require.Equal(t, i != toDelete, exists)
			})
		}
	})
}

func (s *StorageBackendTestSuite) TestStorageTooManySimultaneousSaves() {
	t := s.T()

	// Start the first writer
	w1, err := s.st.OpenWriteStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	// Any attempt to update while the update is in progress should fail now
	w2, err := s.st.OpenWriteStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.ErrorIs(t, err, datastore.ErrUploadInProgress)
	require.Nil(t, w2)

	// Finish the original ingestion
	err = w1.Close()
	require.NoError(t, err)

	// We should be able to successfully read the ingested data
	r, err := s.st.OpenReadStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	b, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, []byte{}, b)

	err = r.Close()
	require.NoError(t, err)
}

func (s *StorageBackendTestSuite) TestStorageSaveWhileDeleting() {
	t := s.T()

	w, err := s.st.OpenWriteStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	exists, err := s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.True(t, exists)

	w, err = s.st.OpenWriteStream(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	err = s.st.Delete(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	exists, err = s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.False(t, exists)

	err = w.Close()
	require.NoError(t, err)

	exists, err = s.st.Exists(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)
	require.True(t, exists)
}
