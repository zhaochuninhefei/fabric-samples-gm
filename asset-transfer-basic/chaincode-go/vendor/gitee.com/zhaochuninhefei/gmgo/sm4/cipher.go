// Copyright (c) 2022 zhaochun
// gmgo is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package sm4

import (
	"crypto/cipher"
	"fmt"

	"gitee.com/zhaochuninhefei/gmgo/internal/subtle"
)

// BlockSize the sm4 block size in bytes.
const BlockSize = 16

const rounds = 32

// A cipher is an instance of SM4 encryption using a particular key.
type sm4Cipher struct {
	enc []uint32
	dec []uint32
}

// NewCipher creates and returns a new cipher.Block.
// The key argument should be the SM4 key,
func NewCipher(key []byte) (cipher.Block, error) {
	k := len(key)
	switch k {
	default:
		return nil, fmt.Errorf("sm4: invalid key size %d", k)
	case 16:
		break
	}
	return newCipher(key)
}

// newCipher creates and returns a new cipher.Block
// implemented in pure Go.
func newCipherGeneric(key []byte) (cipher.Block, error) {
	c := sm4Cipher{make([]uint32, rounds), make([]uint32, rounds)}
	expandKeyGo(key, c.enc, c.dec)
	return &c, nil
}

func (c *sm4Cipher) BlockSize() int { return BlockSize }

func (c *sm4Cipher) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("sm4: input not full block")
	}
	if len(dst) < BlockSize {
		panic("sm4: output not full block")
	}
	if subtle.InexactOverlap(dst[:BlockSize], src[:BlockSize]) {
		panic("sm4: invalid buffer overlap")
	}
	encryptBlockGo(c.enc, dst, src)
}

func (c *sm4Cipher) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("sm4: input not full block")
	}
	if len(dst) < BlockSize {
		panic("sm4: output not full block")
	}
	if subtle.InexactOverlap(dst[:BlockSize], src[:BlockSize]) {
		panic("sm4: invalid buffer overlap")
	}
	decryptBlockGo(c.dec, dst, src)
}
