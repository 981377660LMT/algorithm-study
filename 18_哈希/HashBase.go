// HashBase

package main

import "math/rand"

func main() {

}

func GetHashBase1D() []uint64 {}

func GetHashBase2D() [][]uint64 {}

func randUint64(min, max uint64) uint64 {
	return min + uint64(rand.Int63n(int64(max-min+1)))
}
