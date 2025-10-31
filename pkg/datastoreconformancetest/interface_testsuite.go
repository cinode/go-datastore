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
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/cinode/go-datastore/pkg/blobtypes"
	"github.com/cinode/go-datastore/pkg/common"
	"github.com/cinode/go-datastore/pkg/datastore"
	"github.com/cinode/go-datastore/pkg/datastore/testutils"
	"github.com/cinode/go-datastore/pkg/internal/blobtypes/dynamiclink"
	"github.com/stretchr/testify/suite"
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
	s.Require().NoError(err)
	s.ds = ds
}

func (s *DatastoreTestSuite) TestDatastoreKind() {
	s.Require().Equal(s.expectedKind, s.ds.Kind())
}

func (s *DatastoreTestSuite) TestOpenNonExisting() {
	for _, name := range testutils.EmptyBlobNamesOfAllTypes {
		s.Run(fmt.Sprint(name.Type()), func() {
			r, err := s.ds.Open(context.Background(), name)
			s.Require().ErrorIs(err, datastore.ErrNotFound)
			s.Require().Nil(r)
		})
	}
}

func (s *DatastoreTestSuite) TestOpenInvalidBlobType() {
	bn, err := common.BlobNameFromHashAndType(sha256.New().Sum(nil), common.NewBlobType(0xFF))
	s.Require().NoError(err)

	r, err := s.ds.Open(context.Background(), bn)
	s.Require().ErrorIs(err, blobtypes.ErrUnknownBlobType)
	s.Require().Nil(r)

	err = s.ds.Update(context.Background(), bn, bytes.NewBuffer(nil))
	s.Require().ErrorIs(err, blobtypes.ErrUnknownBlobType)
}

func (s *DatastoreTestSuite) TestBlobValidationFailed() {
	for _, name := range testutils.EmptyBlobNamesOfAllTypes {
		s.Run(fmt.Sprint(name.Type()), func() {
			err := s.ds.Update(context.Background(), name, bytes.NewReader([]byte("test")))
			s.Require().ErrorIs(err, blobtypes.ErrValidationFailed)
		})
	}
}

func (s *DatastoreTestSuite) TestSaveSuccessfulStatic() {
	for _, b := range testutils.TestBlobs {
		exists, err := s.ds.Exists(context.Background(), b.Name)
		s.Require().NoError(err)
		s.Require().False(exists)

		err = s.ds.Update(context.Background(), b.Name, bytes.NewReader(b.Data))
		s.Require().NoError(err)

		exists, err = s.ds.Exists(context.Background(), b.Name)
		s.Require().NoError(err)
		s.Require().True(exists)

		// Overwrite with the same data must be fine
		err = s.ds.Update(context.Background(), b.Name, bytes.NewReader(b.Data))
		s.Require().NoError(err)

		exists, err = s.ds.Exists(context.Background(), b.Name)
		s.Require().NoError(err)
		s.Require().True(exists)

		// Overwrite with wrong data must fail
		err = s.ds.Update(context.Background(), b.Name, bytes.NewReader(append([]byte{0x00}, b.Data...)))
		s.Require().ErrorIs(err, blobtypes.ErrValidationFailed)

		exists, err = s.ds.Exists(context.Background(), b.Name)
		s.Require().NoError(err)
		s.Require().True(exists)

		r, err := s.ds.Open(context.Background(), b.Name)
		s.Require().NoError(err)

		data, err := io.ReadAll(r)
		s.Require().NoError(err)
		s.Require().Equal(b.Data, data)

		err = r.Close()
		s.Require().NoError(err)
	}
}

func (s *DatastoreTestSuite) TestErrorWhileUpdating() {
	for i, b := range testutils.TestBlobs {
		s.Run(fmt.Sprint(i), func() {
			errRet := errors.New("test error")
			err := s.ds.Update(context.Background(), b.Name, testutils.BReader(b.Data, func() error {
				return errRet
			}, nil))
			s.Require().ErrorIs(err, errRet)

			exists, err := s.ds.Exists(context.Background(), b.Name)
			s.Require().NoError(err)
			s.Require().False(exists)
		})
	}
}

func (s *DatastoreTestSuite) TestErrorWhileOverwriting() {
	for i, b := range testutils.TestBlobs {
		s.Run(fmt.Sprint(i), func() {
			err := s.ds.Update(context.Background(), b.Name, bytes.NewReader(b.Data))
			s.Require().NoError(err)

			errRet := errors.New("cancel")

			err = s.ds.Update(context.Background(), b.Name, testutils.BReader(b.Data, func() error {
				exists, err := s.ds.Exists(context.Background(), b.Name)
				s.Require().NoError(err)
				s.Require().True(exists)

				return errRet
			}, nil))

			s.Require().ErrorIs(err, errRet)

			exists, err := s.ds.Exists(context.Background(), b.Name)
			s.Require().NoError(err)
			s.Require().True(exists)

			r, err := s.ds.Open(context.Background(), b.Name)
			s.Require().NoError(err)

			data, err := io.ReadAll(r)
			s.Require().NoError(err)
			s.Require().Equal(b.Data, data)

			err = r.Close()
			s.Require().NoError(err)
		})
	}
}

func (s *DatastoreTestSuite) TestDeleteNonExisting() {
	b := testutils.TestBlobs[0]

	err := s.ds.Update(context.Background(), b.Name, bytes.NewReader(b.Data))
	s.Require().NoError(err)

	err = s.ds.Delete(context.Background(), testutils.TestBlobs[1].Name)
	s.Require().ErrorIs(err, datastore.ErrNotFound)

	exists, err := s.ds.Exists(context.Background(), b.Name)
	s.Require().NoError(err)
	s.Require().True(exists)
}

func (s *DatastoreTestSuite) TestDeleteExisting() {
	b := testutils.TestBlobs[0]
	err := s.ds.Update(context.Background(), b.Name, bytes.NewReader(b.Data))
	s.Require().NoError(err)

	exists, err := s.ds.Exists(context.Background(), b.Name)
	s.Require().NoError(err)
	s.Require().True(exists)

	err = s.ds.Delete(context.Background(), b.Name)
	s.Require().NoError(err)

	exists, err = s.ds.Exists(context.Background(), b.Name)
	s.Require().NoError(err)
	s.Require().False(exists)

	r, err := s.ds.Open(context.Background(), b.Name)
	s.Require().ErrorIs(err, datastore.ErrNotFound)
	s.Require().Nil(r)
}

func (s *DatastoreTestSuite) TestGetKind() {
	k := s.ds.Kind()
	s.Require().NotEmpty(k)
}

func (s *DatastoreTestSuite) TestAddress() {
	address := s.ds.Address()
	s.Require().Regexp(`^[a-zA-Z0-9_-]+://`, address)
}

func (s *DatastoreTestSuite) TestSimultaneousReads() {
	const threadCnt = 10
	const readCnt = 200

	// Prepare data
	for _, b := range testutils.TestBlobs {
		err := s.ds.Update(context.Background(), b.Name, bytes.NewReader(b.Data))
		s.Require().NoError(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(threadCnt)

	for i := 0; i < threadCnt; i++ {
		go func(i int) {
			defer wg.Done()
			for n := 0; n < readCnt; n++ {
				b := testutils.TestBlobs[(i+n)%len(testutils.TestBlobs)]

				r, err := s.ds.Open(context.Background(), b.Name)
				s.Require().NoError(err)

				data, err := io.ReadAll(r)
				s.Require().NoError(err)
				s.Require().Equal(b.Data, data)

				err = r.Close()
				s.Require().NoError(err)
			}
		}(i)
	}

	wg.Wait()
}

func (s *DatastoreTestSuite) TestSimultaneousUpdates() {
	const threadCnt = 3

	b := testutils.TestBlobs[0]
	wg := sync.WaitGroup{}

	for range threadCnt {
		wg.Go(func() {
			err := s.ds.Update(context.Background(), b.Name, bytes.NewReader(b.Data))
			if errors.Is(err, datastore.ErrUploadInProgress) {
				// TODO: We should be able to handle this case
				return
			}

			s.Require().NoError(err)

			exists, err := s.ds.Exists(context.Background(), b.Name)
			s.Require().NoError(err)
			s.Require().True(exists)
		})
	}

	wg.Wait()

	exists, err := s.ds.Exists(context.Background(), b.Name)
	s.Require().NoError(err)
	s.Require().True(exists)

	r, err := s.ds.Open(context.Background(), b.Name)
	s.Require().NoError(err)

	data, err := io.ReadAll(r)
	s.Require().NoError(err)
	s.Require().Equal(b.Data, data)

	err = r.Close()
	s.Require().NoError(err)
}

func (s *DatastoreTestSuite) updateDynamicLink(num int) {
	err := s.ds.Update(
		context.Background(),
		testutils.DynamicLinkPropagationData[num].Name,
		bytes.NewReader(testutils.DynamicLinkPropagationData[num].Data),
	)
	s.Require().NoError(err)
}

func (s *DatastoreTestSuite) readDynamicLinkData() []byte {
	r, err := s.ds.Open(context.Background(), testutils.DynamicLinkPropagationData[0].Name)
	s.Require().NoError(err)

	dl, err := dynamiclink.FromPublicData(testutils.DynamicLinkPropagationData[0].Name, r)
	s.Require().NoError(err)

	elink, err := io.ReadAll(dl.GetEncryptedLinkReader())
	s.Require().NoError(err)

	err = r.Close()
	s.Require().NoError(err)

	return elink
}

func (s *DatastoreTestSuite) expectDynamicLinkData(num int) {
	s.Require().Equal(
		testutils.DynamicLinkPropagationData[num].Expected,
		s.readDynamicLinkData(),
	)
}

func (s *DatastoreTestSuite) TestDynamicLinkPropagation() {
	s.updateDynamicLink(0)
	s.expectDynamicLinkData(0)

	s.updateDynamicLink(1)
	s.expectDynamicLinkData(1)

	s.updateDynamicLink(0)
	s.expectDynamicLinkData(1)

	s.updateDynamicLink(2)
	s.expectDynamicLinkData(2)

	s.updateDynamicLink(1)
	s.expectDynamicLinkData(2)

	s.updateDynamicLink(0)
	s.expectDynamicLinkData(2)
}
