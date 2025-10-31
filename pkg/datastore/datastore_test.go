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

package datastore

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/cinode/go-datastore/pkg/blobtypes"
	"github.com/cinode/go-datastore/pkg/common"
	"github.com/cinode/go-datastore/pkg/datastore/testutils"
	"github.com/cinode/go-datastore/pkg/internal/picotestify/require"
)

type mockStore struct {
	fKind            func() string
	fAddress         func() string
	fOpenReadStream  func(ctx context.Context, name *common.BlobName) (io.ReadCloser, error)
	fOpenWriteStream func(ctx context.Context, name *common.BlobName) (WriteCloseCanceller, error)
	fExists          func(ctx context.Context, name *common.BlobName) (bool, error)
	fDelete          func(ctx context.Context, name *common.BlobName) error
}

func (s *mockStore) Kind() string {
	return s.fKind()
}
func (s *mockStore) Address() string {
	return s.fAddress()
}
func (s *mockStore) OpenReadStream(ctx context.Context, name *common.BlobName) (io.ReadCloser, error) {
	return s.fOpenReadStream(ctx, name)
}
func (s *mockStore) OpenWriteStream(ctx context.Context, name *common.BlobName) (WriteCloseCanceller, error) {
	return s.fOpenWriteStream(ctx, name)
}
func (s *mockStore) Exists(ctx context.Context, name *common.BlobName) (bool, error) {
	return s.fExists(ctx, name)
}
func (s *mockStore) Delete(ctx context.Context, name *common.BlobName) error {
	return s.fDelete(ctx, name)
}

type mockWriteCloseCanceller struct {
	fWrite  func([]byte) (int, error)
	fClose  func() error
	fCancel func()
}

func (m *mockWriteCloseCanceller) Write(b []byte) (int, error) {
	return m.fWrite(b)
}
func (m *mockWriteCloseCanceller) Close() error {
	return m.fClose()
}
func (m *mockWriteCloseCanceller) Cancel() {
	m.fCancel()
}

func TestDatastoreWriteFailure(t *testing.T) {
	t.Run("error on opening write stream", func(t *testing.T) {
		errRet := errors.New("error")
		ds := &datastore{s: &mockStore{
			fOpenWriteStream: func(ctx context.Context, name *common.BlobName) (WriteCloseCanceller, error) {
				return nil, errRet
			},
		}}

		err := ds.Update(t.Context(), testutils.EmptyBlobNameStatic, bytes.NewBuffer(nil))
		require.ErrorIs(t, err, errRet)
	})

	t.Run("error on closing write stream", func(t *testing.T) {
		errRet := errors.New("error")

		closeCalled := false
		cancelCalled := false
		ds := &datastore{s: &mockStore{
			fOpenWriteStream: func(ctx context.Context, name *common.BlobName) (WriteCloseCanceller, error) {
				return &mockWriteCloseCanceller{
					fWrite: func(b []byte) (int, error) {
						require.False(t, closeCalled)
						return len(b), nil
					},
					fClose: func() error {
						require.False(t, closeCalled)
						require.False(t, cancelCalled)
						closeCalled = true
						return errRet
					},
					fCancel: func() {
						require.True(t, closeCalled)
						require.False(t, cancelCalled)
						cancelCalled = true
					},
				}, nil
			},
			fOpenReadStream: func(ctx context.Context, name *common.BlobName) (io.ReadCloser, error) {
				return nil, ErrNotFound
			},
		}}

		err := ds.Update(t.Context(), testutils.EmptyBlobNameStatic, bytes.NewBuffer(nil))
		require.ErrorIs(t, err, errRet)

		// Failed Close call will be followed by a Cancel call
		require.True(t, closeCalled)
		require.True(t, cancelCalled)
	})
}

func TestDatastoreDetectCorruptedRead(t *testing.T) {
	ds := InMemory()
	mem := ds.(*datastore).s.(*memory)
	mem.bmap[testutils.EmptyBlobNameStatic.String()] = []byte("I should not be here")

	r, err := ds.Open(t.Context(), testutils.EmptyBlobNameStatic)
	require.NoError(t, err)

	_, err = io.ReadAll(r)
	require.ErrorIs(t, err, blobtypes.ErrValidationFailed)

	err = r.Close()
	require.NoError(t, err)
}
