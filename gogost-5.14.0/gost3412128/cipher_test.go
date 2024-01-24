// GoGOST -- Pure Go GOST cryptographic functions library
// Copyright (C) 2015-2024 Sergey Matveev <stargrave@stargrave.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy_db of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package gost3412128

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"testing"
	"testing/quick"
)

var (
	key []byte = []byte{
		0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	}
	pt [BlockSize]byte = [BlockSize]byte{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
	}
	ct [BlockSize]byte = [BlockSize]byte{
		0x7f, 0x67, 0x9d, 0x90, 0xbe, 0xbc, 0x24, 0x30,
		0x5a, 0x46, 0x8d, 0x42, 0xb9, 0xd4, 0xed, 0xcd,
	}
)

func TestCipherInterface(t *testing.T) {
	var _ cipher.Block = NewCipher(make([]byte, KeySize))
}

func TestRandom(t *testing.T) {
	data := make([]byte, BlockSize)
	f := func(key [KeySize]byte, pt [BlockSize]byte) bool {
		io.ReadFull(rand.Reader, key[:])
		c := NewCipher(key[:])
		c.Encrypt(data, pt[:])
		c.Decrypt(data, data)
		return bytes.Equal(data, pt[:])
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkEncrypt(b *testing.B) {
	key := make([]byte, KeySize)
	io.ReadFull(rand.Reader, key)
	c := NewCipher(key)
	blk := make([]byte, BlockSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(blk, blk)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	key := make([]byte, KeySize)
	io.ReadFull(rand.Reader, key)
	c := NewCipher(key)
	blk := make([]byte, BlockSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(blk, blk)
	}
}

func TestS(t *testing.T) {
	blk := [BlockSize]byte{
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
	}
	s(&blk)
	if !bytes.Equal(blk[:], []byte{
		0xb6, 0x6c, 0xd8, 0x88, 0x7d, 0x38, 0xe8, 0xd7,
		0x77, 0x65, 0xae, 0xea, 0x0c, 0x9a, 0x7e, 0xfc,
	}) {
		t.FailNow()
	}
	s(&blk)
	if !bytes.Equal(blk[:], []byte{
		0x55, 0x9d, 0x8d, 0xd7, 0xbd, 0x06, 0xcb, 0xfe,
		0x7e, 0x7b, 0x26, 0x25, 0x23, 0x28, 0x0d, 0x39,
	}) {
		t.FailNow()
	}
	s(&blk)
	if !bytes.Equal(blk[:], []byte{
		0x0c, 0x33, 0x22, 0xfe, 0xd5, 0x31, 0xe4, 0x63,
		0x0d, 0x80, 0xef, 0x5c, 0x5a, 0x81, 0xc5, 0x0b,
	}) {
		t.FailNow()
	}
	s(&blk)
	if !bytes.Equal(blk[:], []byte{
		0x23, 0xae, 0x65, 0x63, 0x3f, 0x84, 0x2d, 0x29,
		0xc5, 0xdf, 0x52, 0x9c, 0x13, 0xf5, 0xac, 0xda,
	}) {
		t.FailNow()
	}
}

func R(blk []byte) {
	t := blk[15]
	for i := 0; i < 15; i++ {
		t ^= gfCache[blk[i]][lc[i]]
	}
	copy(blk[1:], blk)
	blk[0] = t
}

func TestR(t *testing.T) {
	blk := [BlockSize]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00,
	}
	R(blk[:])
	if !bytes.Equal(blk[:], []byte{
		0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}) {
		t.FailNow()
	}
	R(blk[:])
	if !bytes.Equal(blk[:], []byte{
		0xa5, 0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}) {
		t.FailNow()
	}
	R(blk[:])
	if !bytes.Equal(blk[:], []byte{
		0x64, 0xa5, 0x94, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}) {
		t.FailNow()
	}
	R(blk[:])
	if !bytes.Equal(blk[:], []byte{
		0x0d, 0x64, 0xa5, 0x94, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}) {
		t.FailNow()
	}
}

func TestL(t *testing.T) {
	blk := [BlockSize]byte{
		0x64, 0xa5, 0x94, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	l(&blk)
	if !bytes.Equal(blk[:], []byte{
		0xd4, 0x56, 0x58, 0x4d, 0xd0, 0xe3, 0xe8, 0x4c,
		0xc3, 0x16, 0x6e, 0x4b, 0x7f, 0xa2, 0x89, 0x0d,
	}) {
		t.FailNow()
	}
	l(&blk)
	if !bytes.Equal(blk[:], []byte{
		0x79, 0xd2, 0x62, 0x21, 0xb8, 0x7b, 0x58, 0x4c,
		0xd4, 0x2f, 0xbc, 0x4f, 0xfe, 0xa5, 0xde, 0x9a,
	}) {
		t.FailNow()
	}
	l(&blk)
	if !bytes.Equal(blk[:], []byte{
		0x0e, 0x93, 0x69, 0x1a, 0x0c, 0xfc, 0x60, 0x40,
		0x8b, 0x7b, 0x68, 0xf6, 0x6b, 0x51, 0x3c, 0x13,
	}) {
		t.FailNow()
	}
	l(&blk)
	if !bytes.Equal(blk[:], []byte{
		0xe6, 0xa8, 0x09, 0x4f, 0xee, 0x0a, 0xa2, 0x04,
		0xfd, 0x97, 0xbc, 0xb0, 0xb4, 0x4b, 0x85, 0x80,
	}) {
		t.FailNow()
	}
}

func TestC(t *testing.T) {
	if !bytes.Equal(cBlk[0][:], []byte{
		0x6e, 0xa2, 0x76, 0x72, 0x6c, 0x48, 0x7a, 0xb8,
		0x5d, 0x27, 0xbd, 0x10, 0xdd, 0x84, 0x94, 0x01,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(cBlk[1][:], []byte{
		0xdc, 0x87, 0xec, 0xe4, 0xd8, 0x90, 0xf4, 0xb3,
		0xba, 0x4e, 0xb9, 0x20, 0x79, 0xcb, 0xeb, 0x02,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(cBlk[2][:], []byte{
		0xb2, 0x25, 0x9a, 0x96, 0xb4, 0xd8, 0x8e, 0x0b,
		0xe7, 0x69, 0x04, 0x30, 0xa4, 0x4f, 0x7f, 0x03,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(cBlk[3][:], []byte{
		0x7b, 0xcd, 0x1b, 0x0b, 0x73, 0xe3, 0x2b, 0xa5,
		0xb7, 0x9c, 0xb1, 0x40, 0xf2, 0x55, 0x15, 0x04,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(cBlk[4][:], []byte{
		0x15, 0x6f, 0x6d, 0x79, 0x1f, 0xab, 0x51, 0x1d,
		0xea, 0xbb, 0x0c, 0x50, 0x2f, 0xd1, 0x81, 0x05,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(cBlk[5][:], []byte{
		0xa7, 0x4a, 0xf7, 0xef, 0xab, 0x73, 0xdf, 0x16,
		0x0d, 0xd2, 0x08, 0x60, 0x8b, 0x9e, 0xfe, 0x06,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(cBlk[6][:], []byte{
		0xc9, 0xe8, 0x81, 0x9d, 0xc7, 0x3b, 0xa5, 0xae,
		0x50, 0xf5, 0xb5, 0x70, 0x56, 0x1a, 0x6a, 0x07,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(cBlk[7][:], []byte{
		0xf6, 0x59, 0x36, 0x16, 0xe6, 0x05, 0x56, 0x89,
		0xad, 0xfb, 0xa1, 0x80, 0x27, 0xaa, 0x2a, 0x08,
	}) {
		t.FailNow()
	}
}

func TestRoundKeys(t *testing.T) {
	c := NewCipher(key)
	if !bytes.Equal(c.ks[0][:], []byte{
		0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[1][:], []byte{
		0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[2][:], []byte{
		0xdb, 0x31, 0x48, 0x53, 0x15, 0x69, 0x43, 0x43,
		0x22, 0x8d, 0x6a, 0xef, 0x8c, 0xc7, 0x8c, 0x44,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[3][:], []byte{
		0x3d, 0x45, 0x53, 0xd8, 0xe9, 0xcf, 0xec, 0x68,
		0x15, 0xeb, 0xad, 0xc4, 0x0a, 0x9f, 0xfd, 0x04,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[4][:], []byte{
		0x57, 0x64, 0x64, 0x68, 0xc4, 0x4a, 0x5e, 0x28,
		0xd3, 0xe5, 0x92, 0x46, 0xf4, 0x29, 0xf1, 0xac,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[5][:], []byte{
		0xbd, 0x07, 0x94, 0x35, 0x16, 0x5c, 0x64, 0x32,
		0xb5, 0x32, 0xe8, 0x28, 0x34, 0xda, 0x58, 0x1b,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[6][:], []byte{
		0x51, 0xe6, 0x40, 0x75, 0x7e, 0x87, 0x45, 0xde,
		0x70, 0x57, 0x27, 0x26, 0x5a, 0x00, 0x98, 0xb1,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[7][:], []byte{
		0x5a, 0x79, 0x25, 0x01, 0x7b, 0x9f, 0xdd, 0x3e,
		0xd7, 0x2a, 0x91, 0xa2, 0x22, 0x86, 0xf9, 0x84,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[8][:], []byte{
		0xbb, 0x44, 0xe2, 0x53, 0x78, 0xc7, 0x31, 0x23,
		0xa5, 0xf3, 0x2f, 0x73, 0xcd, 0xb6, 0xe5, 0x17,
	}) {
		t.FailNow()
	}
	if !bytes.Equal(c.ks[9][:], []byte{
		0x72, 0xe9, 0xdd, 0x74, 0x16, 0xbc, 0xf4, 0x5b,
		0x75, 0x5d, 0xba, 0xa8, 0x8e, 0x4a, 0x40, 0x43,
	}) {
		t.FailNow()
	}
}

func TestVectorEncrypt(t *testing.T) {
	c := NewCipher(key)
	dst := make([]byte, BlockSize)
	c.Encrypt(dst, pt[:])
	if !bytes.Equal(dst, ct[:]) {
		t.FailNow()
	}
}

func TestVectorDecrypt(t *testing.T) {
	c := NewCipher(key)
	dst := make([]byte, BlockSize)
	c.Decrypt(dst, ct[:])
	if !bytes.Equal(dst, pt[:]) {
		t.FailNow()
	}
}
