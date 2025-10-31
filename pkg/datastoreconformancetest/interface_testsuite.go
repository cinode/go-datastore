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
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/cinode/go-common/blob"
	"github.com/cinode/go-common/picotestify/require"
	"github.com/cinode/go-common/picotestify/suite"
	"github.com/cinode/go-datastore/pkg/blobtypes"
	"github.com/cinode/go-datastore/pkg/datastore"
	"github.com/cinode/go-datastore/pkg/datastore/testutils"
	"github.com/cinode/go-datastore/pkg/internal/blobtypes/dynamiclink"
)

type DatastoreTestSuite struct {
	suite.Suite

	// ds is the current datastore instance to test.
	ds datastore.DS

	// createDS is a function that creates a new datastore for each test.
	createDS func() (datastore.DS, error)

	// expectedKind is the expected kind of the datastore
	expectedKind string
}

func NewDatastoreTestSuite(createDS func() (datastore.DS, error), expectedKind string) *DatastoreTestSuite {
	return &DatastoreTestSuite{
		createDS:     createDS,
		expectedKind: expectedKind,
	}
}

func (s *DatastoreTestSuite) SetupTest() {
	ds, err := s.createDS()
	require.NoError(s.T(), err)
	s.ds = ds
}

func (s *DatastoreTestSuite) TestDatastoreKind() {
	require.Equal(s.T(), s.expectedKind, s.ds.Kind())
}

func (s *DatastoreTestSuite) TestOpenNonExisting() {
	t := s.T()

	for _, name := range testutils.EmptyBlobNamesOfAllTypes {
		t.Run(fmt.Sprint(name.Type()), func(t *testing.T) {
			r, err := s.ds.Open(t.Context(), name)
			require.ErrorIs(t, err, datastore.ErrNotFound)
			require.Nil(t, r)
		})
	}
}

func (s *DatastoreTestSuite) TestOpenInvalidBlobType() {
	t := s.T()

	bn, err := blob.NameFromHashAndType(sha256.New().Sum(nil), blob.NewType(0xFF))
	require.NoError(t, err)

	r, err := s.ds.Open(t.Context(), bn)
	require.ErrorIs(t, err, blobtypes.ErrUnknownBlobType)
	require.Nil(t, r)

	err = s.ds.Update(t.Context(), bn, bytes.NewBuffer(nil))
	require.ErrorIs(t, err, blobtypes.ErrUnknownBlobType)
}

func (s *DatastoreTestSuite) TestBlobValidationFailed() {
	t := s.T()

	for _, name := range testutils.EmptyBlobNamesOfAllTypes {
		t.Run(fmt.Sprint(name.Type()), func(t *testing.T) {
			err := s.ds.Update(t.Context(), name, bytes.NewReader([]byte("test")))
			require.ErrorIs(t, err, blobtypes.ErrValidationFailed)
		})
	}
}

func (s *DatastoreTestSuite) TestSaveSuccessfulStatic() {
	t := s.T()

	for _, b := range testutils.TestBlobs {
		exists, err := s.ds.Exists(t.Context(), b.Name)
		require.NoError(t, err)
		require.False(t, exists)

		err = s.ds.Update(t.Context(), b.Name, bytes.NewReader(b.Data))
		require.NoError(t, err)

		exists, err = s.ds.Exists(t.Context(), b.Name)
		require.NoError(t, err)
		require.True(t, exists)

		// Overwrite with the same data must be fine
		err = s.ds.Update(t.Context(), b.Name, bytes.NewReader(b.Data))
		require.NoError(t, err)

		exists, err = s.ds.Exists(t.Context(), b.Name)
		require.NoError(t, err)
		require.True(t, exists)

		// Overwrite with wrong data must fail
		err = s.ds.Update(t.Context(), b.Name, bytes.NewReader(append([]byte{0x00}, b.Data...)))
		require.ErrorIs(t, err, blobtypes.ErrValidationFailed)

		exists, err = s.ds.Exists(t.Context(), b.Name)
		require.NoError(t, err)
		require.True(t, exists)

		r, err := s.ds.Open(t.Context(), b.Name)
		require.NoError(t, err)

		data, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, b.Data, data)

		err = r.Close()
		require.NoError(t, err)
	}
}

func (s *DatastoreTestSuite) TestErrorWhileUpdating() {
	t := s.T()

	for i, b := range testutils.TestBlobs {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			errRet := errors.New("test error")
			err := s.ds.Update(t.Context(), b.Name, testutils.BReader(b.Data, func() error {
				return errRet
			}, nil))
			require.ErrorIs(t, err, errRet)

			exists, err := s.ds.Exists(t.Context(), b.Name)
			require.NoError(t, err)
			require.False(t, exists)
		})
	}
}

func (s *DatastoreTestSuite) TestErrorWhileOverwriting() {
	t := s.T()

	for i, b := range testutils.TestBlobs {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			err := s.ds.Update(t.Context(), b.Name, bytes.NewReader(b.Data))
			require.NoError(t, err)

			errRet := errors.New("cancel")

			err = s.ds.Update(t.Context(), b.Name, testutils.BReader(b.Data, func() error {
				exists, err := s.ds.Exists(t.Context(), b.Name)
				require.NoError(t, err)
				require.True(t, exists)

				return errRet
			}, nil))

			require.ErrorIs(t, err, errRet)

			exists, err := s.ds.Exists(t.Context(), b.Name)
			require.NoError(t, err)
			require.True(t, exists)

			r, err := s.ds.Open(t.Context(), b.Name)
			require.NoError(t, err)

			data, err := io.ReadAll(r)
			require.NoError(t, err)
			require.Equal(t, b.Data, data)

			err = r.Close()
			require.NoError(t, err)
		})
	}
}

func (s *DatastoreTestSuite) TestDeleteNonExisting() {
	t := s.T()

	b := testutils.TestBlobs[0]

	err := s.ds.Update(t.Context(), b.Name, bytes.NewReader(b.Data))
	require.NoError(t, err)

	err = s.ds.Delete(t.Context(), testutils.TestBlobs[1].Name)
	require.ErrorIs(t, err, datastore.ErrNotFound)

	exists, err := s.ds.Exists(t.Context(), b.Name)
	require.NoError(t, err)
	require.True(t, exists)
}

func (s *DatastoreTestSuite) TestDeleteExisting() {
	t := s.T()

	b := testutils.TestBlobs[0]
	err := s.ds.Update(t.Context(), b.Name, bytes.NewReader(b.Data))
	require.NoError(t, err)

	exists, err := s.ds.Exists(t.Context(), b.Name)
	require.NoError(t, err)
	require.True(t, exists)

	err = s.ds.Delete(t.Context(), b.Name)
	require.NoError(t, err)

	exists, err = s.ds.Exists(t.Context(), b.Name)
	require.NoError(t, err)
	require.False(t, exists)

	r, err := s.ds.Open(t.Context(), b.Name)
	require.ErrorIs(t, err, datastore.ErrNotFound)
	require.Nil(t, r)
}

func (s *DatastoreTestSuite) TestGetKind() {
	t := s.T()

	k := s.ds.Kind()
	require.NotEmpty(t, k)
}

func (s *DatastoreTestSuite) TestAddress() {
	t := s.T()

	address := s.ds.Address()
	require.Regexp(t, `^[a-zA-Z0-9_-]+://`, address)
}

func (s *DatastoreTestSuite) TestSimultaneousReads() {
	t := s.T()

	const threadCnt = 10
	const readCnt = 200

	// Prepare data
	for _, b := range testutils.TestBlobs {
		err := s.ds.Update(t.Context(), b.Name, bytes.NewReader(b.Data))
		require.NoError(t, err)
	}

	wg := sync.WaitGroup{}
	wg.Add(threadCnt)

	for i := 0; i < threadCnt; i++ {
		go func(i int) {
			defer wg.Done()
			for n := 0; n < readCnt; n++ {
				b := testutils.TestBlobs[(i+n)%len(testutils.TestBlobs)]

				r, err := s.ds.Open(t.Context(), b.Name)
				require.NoError(t, err)

				data, err := io.ReadAll(r)
				require.NoError(t, err)
				require.Equal(t, b.Data, data)

				err = r.Close()
				require.NoError(t, err)
			}
		}(i)
	}

	wg.Wait()
}

func (s *DatastoreTestSuite) TestSimultaneousUpdates() {
	t := s.T()

	const threadCnt = 3

	b := testutils.TestBlobs[0]
	wg := sync.WaitGroup{}

	for range threadCnt {
		wg.Go(func() {
			err := s.ds.Update(t.Context(), b.Name, bytes.NewReader(b.Data))
			if errors.Is(err, datastore.ErrUploadInProgress) {
				// TODO: We should be able to handle this case
				return
			}

			require.NoError(t, err)

			exists, err := s.ds.Exists(t.Context(), b.Name)
			require.NoError(t, err)
			require.True(t, exists)
		})
	}

	wg.Wait()

	exists, err := s.ds.Exists(t.Context(), b.Name)
	require.NoError(t, err)
	require.True(t, exists)

	r, err := s.ds.Open(t.Context(), b.Name)
	require.NoError(t, err)

	data, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, b.Data, data)

	err = r.Close()
	require.NoError(t, err)
}

func (s *DatastoreTestSuite) updateDynamicLink(t *testing.T, num int) {
	err := s.ds.Update(
		t.Context(),
		testutils.DynamicLinkPropagationData[num].Name,
		bytes.NewReader(testutils.DynamicLinkPropagationData[num].Data),
	)
	require.NoError(t, err)
}

func (s *DatastoreTestSuite) readDynamicLinkData(t *testing.T) []byte {
	r, err := s.ds.Open(t.Context(), testutils.DynamicLinkPropagationData[0].Name)
	require.NoError(t, err)

	dl, err := dynamiclink.FromPublicData(testutils.DynamicLinkPropagationData[0].Name, r)
	require.NoError(t, err)

	elink, err := io.ReadAll(dl.GetEncryptedLinkReader())
	require.NoError(t, err)

	err = r.Close()
	require.NoError(t, err)

	return elink
}

func (s *DatastoreTestSuite) expectDynamicLinkData(t *testing.T, num int) {
	require.Equal(t,
		testutils.DynamicLinkPropagationData[num].Expected,
		s.readDynamicLinkData(t),
	)
}

func (s *DatastoreTestSuite) TestDynamicLinkPropagation() {
	t := s.T()

	s.updateDynamicLink(t, 0)
	s.expectDynamicLinkData(t, 0)

	s.updateDynamicLink(t, 1)
	s.expectDynamicLinkData(t, 1)

	s.updateDynamicLink(t, 0)
	s.expectDynamicLinkData(t, 1)

	s.updateDynamicLink(t, 2)
	s.expectDynamicLinkData(t, 2)

	s.updateDynamicLink(t, 1)
	s.expectDynamicLinkData(t, 2)

	s.updateDynamicLink(t, 0)
	s.expectDynamicLinkData(t, 2)
}
