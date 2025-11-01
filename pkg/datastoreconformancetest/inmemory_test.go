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

package datastoreconformancetest_test

import (
	"testing"

	"github.com/cinode/go-common/picotestify/suite"
	"github.com/cinode/go-datastore/pkg/datastore"
	"github.com/cinode/go-datastore/pkg/datastoreconformancetest"
)

func TestInMemoryDatastoreTestSuite(t *testing.T) {
	suite.Run(t, datastoreconformancetest.NewDatastoreTestSuite(
		func() (datastore.DS, error) { return datastore.InMemory(), nil },
		"Memory",
	))
}

func TestInMemoryStorageBackendTestsuite(t *testing.T) {
	suite.Run(t, datastoreconformancetest.NewStorageBackendTestSuite(
		func() (datastore.StorageBackend, error) { return datastore.NewInMemoryStorageBackend(), nil },
		"Memory",
	))
}
