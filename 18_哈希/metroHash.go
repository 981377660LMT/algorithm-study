// https://github.com/dgryski/go-metro/blob/master/metro64.go
// Hash64、Hash64Str、Hash128

package metro

import (
	"encoding/binary"
	"math/bits"
)

func Hash64(buffer []byte, seed uint64) uint64 {

	const (
		k0 = 0xD6D018F5
		k1 = 0xA2AA033B
		k2 = 0x62992FC1
		k3 = 0x30BC5B29
	)

	ptr := buffer

	hash := (seed + k2) * k0

	if len(ptr) >= 32 {
		v0, v1, v2, v3 := hash, hash, hash, hash

		for len(ptr) >= 32 {
			v0 += binary.LittleEndian.Uint64(ptr[:8]) * k0
			v0 = bits.RotateLeft64(v0, -29) + v2
			v1 += binary.LittleEndian.Uint64(ptr[8:16]) * k1
			v1 = bits.RotateLeft64(v1, -29) + v3
			v2 += binary.LittleEndian.Uint64(ptr[16:24]) * k2
			v2 = bits.RotateLeft64(v2, -29) + v0
			v3 += binary.LittleEndian.Uint64(ptr[24:32]) * k3
			v3 = bits.RotateLeft64(v3, -29) + v1
			ptr = ptr[32:]
		}

		v2 ^= bits.RotateLeft64(((v0+v3)*k0)+v1, -37) * k1
		v3 ^= bits.RotateLeft64(((v1+v2)*k1)+v0, -37) * k0
		v0 ^= bits.RotateLeft64(((v0+v2)*k0)+v3, -37) * k1
		v1 ^= bits.RotateLeft64(((v1+v3)*k1)+v2, -37) * k0
		hash += v0 ^ v1
	}

	if len(ptr) >= 16 {
		v0 := hash + (binary.LittleEndian.Uint64(ptr[:8]) * k2)
		v0 = bits.RotateLeft64(v0, -29) * k3
		v1 := hash + (binary.LittleEndian.Uint64(ptr[8:16]) * k2)
		v1 = bits.RotateLeft64(v1, -29) * k3
		v0 ^= bits.RotateLeft64(v0*k0, -21) + v1
		v1 ^= bits.RotateLeft64(v1*k3, -21) + v0
		hash += v1
		ptr = ptr[16:]
	}

	if len(ptr) >= 8 {
		hash += binary.LittleEndian.Uint64(ptr[:8]) * k3
		ptr = ptr[8:]
		hash ^= bits.RotateLeft64(hash, -55) * k1
	}

	if len(ptr) >= 4 {
		hash += uint64(binary.LittleEndian.Uint32(ptr[:4])) * k3
		hash ^= bits.RotateLeft64(hash, -26) * k1
		ptr = ptr[4:]
	}

	if len(ptr) >= 2 {
		hash += uint64(binary.LittleEndian.Uint16(ptr[:2])) * k3
		ptr = ptr[2:]
		hash ^= bits.RotateLeft64(hash, -48) * k1
	}

	if len(ptr) >= 1 {
		hash += uint64(ptr[0]) * k3
		hash ^= bits.RotateLeft64(hash, -37) * k1
	}

	hash ^= bits.RotateLeft64(hash, -28)
	hash *= k0
	hash ^= bits.RotateLeft64(hash, -29)

	return hash
}

func Hash64Str(buffer string, seed uint64) uint64 {
	return Hash64([]byte(buffer), seed)
}

func Hash128(buffer []byte, seed uint64) (uint64, uint64) {

	const (
		k0 = 0xC83A91E1
		k1 = 0x8648DBDB
		k2 = 0x7BDEC03B
		k3 = 0x2F5870A5
	)

	ptr := buffer

	var v [4]uint64

	v[0] = (seed - k0) * k3
	v[1] = (seed + k1) * k2

	if len(ptr) >= 32 {
		v[2] = (seed + k0) * k2
		v[3] = (seed - k1) * k3

		for len(ptr) >= 32 {
			v[0] += binary.LittleEndian.Uint64(ptr) * k0
			ptr = ptr[8:]
			v[0] = rotate_right(v[0], 29) + v[2]
			v[1] += binary.LittleEndian.Uint64(ptr) * k1
			ptr = ptr[8:]
			v[1] = rotate_right(v[1], 29) + v[3]
			v[2] += binary.LittleEndian.Uint64(ptr) * k2
			ptr = ptr[8:]
			v[2] = rotate_right(v[2], 29) + v[0]
			v[3] += binary.LittleEndian.Uint64(ptr) * k3
			ptr = ptr[8:]
			v[3] = rotate_right(v[3], 29) + v[1]
		}

		v[2] ^= rotate_right(((v[0]+v[3])*k0)+v[1], 21) * k1
		v[3] ^= rotate_right(((v[1]+v[2])*k1)+v[0], 21) * k0
		v[0] ^= rotate_right(((v[0]+v[2])*k0)+v[3], 21) * k1
		v[1] ^= rotate_right(((v[1]+v[3])*k1)+v[2], 21) * k0
	}

	if len(ptr) >= 16 {
		v[0] += binary.LittleEndian.Uint64(ptr) * k2
		ptr = ptr[8:]
		v[0] = rotate_right(v[0], 33) * k3
		v[1] += binary.LittleEndian.Uint64(ptr) * k2
		ptr = ptr[8:]
		v[1] = rotate_right(v[1], 33) * k3
		v[0] ^= rotate_right((v[0]*k2)+v[1], 45) * k1
		v[1] ^= rotate_right((v[1]*k3)+v[0], 45) * k0
	}

	if len(ptr) >= 8 {
		v[0] += binary.LittleEndian.Uint64(ptr) * k2
		ptr = ptr[8:]
		v[0] = rotate_right(v[0], 33) * k3
		v[0] ^= rotate_right((v[0]*k2)+v[1], 27) * k1
	}

	if len(ptr) >= 4 {
		v[1] += uint64(binary.LittleEndian.Uint32(ptr)) * k2
		ptr = ptr[4:]
		v[1] = rotate_right(v[1], 33) * k3
		v[1] ^= rotate_right((v[1]*k3)+v[0], 46) * k0
	}

	if len(ptr) >= 2 {
		v[0] += uint64(binary.LittleEndian.Uint16(ptr)) * k2
		ptr = ptr[2:]
		v[0] = rotate_right(v[0], 33) * k3
		v[0] ^= rotate_right((v[0]*k2)+v[1], 22) * k1
	}

	if len(ptr) >= 1 {
		v[1] += uint64(ptr[0]) * k2
		v[1] = rotate_right(v[1], 33) * k3
		v[1] ^= rotate_right((v[1]*k3)+v[0], 58) * k0
	}

	v[0] += rotate_right((v[0]*k0)+v[1], 13)
	v[1] += rotate_right((v[1]*k1)+v[0], 37)
	v[0] += rotate_right((v[0]*k2)+v[1], 13)
	v[1] += rotate_right((v[1]*k3)+v[0], 37)

	return v[0], v[1]
}

func rotate_right(v uint64, k uint) uint64 {
	return (v >> k) | (v << (64 - k))
}
