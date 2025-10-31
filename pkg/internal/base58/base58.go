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

package base58

import (
	"errors"
	"fmt"
	"math/big"
)

var ErrInvalidBase58Character = errors.New("invalid base58 character")

var bn2btc, btc2bn = func() (bn2btc, btc2bn [256]byte) {
	const btcAlphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	const bitNumDigits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUV"

	for i := range btcAlphabet {
		bn2btc[bitNumDigits[i]] = btcAlphabet[i]
		btc2bn[btcAlphabet[i]] = bitNumDigits[i]
	}

	return bn2btc, btc2bn
}()

func Encode(data []byte) string {
	// Count leading zero bytes, those are treated especially
	leadingZeros := 0
	for _, b := range data {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	var bi big.Int
	bi.SetBytes(data[leadingZeros:])

	txt := bi.Text(58)

	if txt == "0" {
		txt = ""
	}

	res := make([]byte, leadingZeros+len(txt))

	// Add '1' characters for leading zeros
	for i := 0; i < leadingZeros; i++ {
		res[i] = '1'
	}

	// Convert the rest
	for i, b := range txt {
		res[leadingZeros+i] = bn2btc[b]
	}

	return string(res)
}

func Decode(s string) ([]byte, error) {
	// Split into leading bytes and bit.Int-compatible representation
	leadingZeros := 0
	leadingZerosDone := false
	bnText := make([]byte, 0, len(s))
	for _, b := range s {
		if !leadingZerosDone && b == '1' {
			leadingZeros++
		} else if btc2bn[b] != 0 {
			leadingZerosDone = true
			bnText = append(bnText, btc2bn[b])
		} else {
			return nil, fmt.Errorf("%w: '%v'", ErrInvalidBase58Character, b)
		}
	}

	if len(bnText) == 0 {
		return make([]byte, leadingZeros), nil
	}

	var bn big.Int
	bn.SetString(string(bnText), 58)

	bnBytes := bn.Bytes()

	result := make([]byte, leadingZeros+len(bnBytes))
	copy(result[leadingZeros:], bnBytes)

	return result, nil
}
