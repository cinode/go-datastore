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

package cipherfactory

import (
	"crypto/sha256"
	"hash"
	"io"

	"github.com/cinode/go-common/blob"
	"golang.org/x/crypto/chacha20"
)

const (
	preambleHashKey       = 0x01
	preambleHashIV        = 0x02
	preambleHashDefaultIV = 0x03
)

type KeyGenerator interface {
	io.Writer
	Generate() *blob.Key
}

type keyGenerator struct {
	h hash.Hash
}

func (g keyGenerator) Write(b []byte) (int, error) { return g.h.Write(b) }

func (g keyGenerator) Generate() *blob.Key {
	return blob.KeyFromBytes(append(
		[]byte{reservedByteForKeyType},
		g.h.Sum(nil)[:chacha20.KeySize]...,
	))
}

type IVGenerator interface {
	io.Writer
	Generate() *blob.IV
}

type ivGenerator struct {
	h hash.Hash
}

func (g ivGenerator) Write(b []byte) (int, error) { return g.h.Write(b) }

func (g ivGenerator) Generate() *blob.IV {
	return blob.IVFromBytes(g.h.Sum(nil)[:chacha20.NonceSizeX])
}

func NewKeyGenerator(t blob.Type) KeyGenerator {
	h := sha256.New()
	h.Write([]byte{preambleHashKey, reservedByteForKeyType, t.IDByte()})
	return keyGenerator{h: h}
}

func NewIVGenerator(t blob.Type) IVGenerator {
	h := sha256.New()
	h.Write([]byte{preambleHashIV, reservedByteForKeyType, t.IDByte()})
	return ivGenerator{h: h}
}

var defaultXChaCha20IV = func() *blob.IV {
	h := sha256.New()
	h.Write([]byte{preambleHashDefaultIV, reservedByteForKeyType})
	return blob.IVFromBytes(h.Sum(nil)[:chacha20.NonceSizeX])
}()

func DefaultIV(k *blob.Key) *blob.IV {
	return defaultXChaCha20IV
}
