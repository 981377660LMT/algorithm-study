// template <typename T>
// struct Index_Compression_DISTINCT_SMALL {
//   static_assert(is_same_v<T, int>);
//   int mi, ma;
//   vc<int> dat;
//   vc<int> build(vc<int> X) {
//     mi = 0, ma = -1;
//     if (!X.empty()) mi = MIN(X), ma = MAX(X);
//     dat.assign(ma - mi + 2, 0);
//     for (auto& x: X) dat[x - mi + 1]++;
//     FOR(i, len(dat) - 1) dat[i + 1] += dat[i];
//     for (auto& x: X) { x = dat[x - mi]++; }
//     FOR_R(i, 1, len(dat)) dat[i] = dat[i - 1];
//     dat[0] = 0;
//     return X;
//   }
//   int operator()(ll x) { return dat[clamp<ll>(x - mi, 0, ma - mi + 1)]; }
// };

// template <typename T>
// struct Index_Compression_SAME_SMALL {
//   static_assert(is_same_v<T, int>);
//   int mi, ma;
//   vc<int> dat;
//   vc<int> build(vc<int> X) {
//     mi = 0, ma = -1;
//     if (!X.empty()) mi = MIN(X), ma = MAX(X);
//     dat.assign(ma - mi + 2, 0);
//     for (auto& x: X) dat[x - mi + 1] = 1;
//     FOR(i, len(dat) - 1) dat[i + 1] += dat[i];
//     for (auto& x: X) { x = dat[x - mi]; }
//     return X;
//   }
//   int operator()(ll x) { return dat[clamp<ll>(x - mi, 0, ma - mi + 1)]; }
// };

// template <typename T>
// struct Index_Compression_SAME_LARGE {
//   vc<T> dat;
//   vc<int> build(vc<T> X) {
//     vc<int> I = argsort(X);
//     vc<int> res(len(X));
//     for (auto& i: I) {
//       if (!dat.empty() && dat.back() == X[i]) {
//         res[i] = len(dat) - 1;
//       } else {
//         res[i] = len(dat);
//         dat.eb(X[i]);
//       }
//     }
//     dat.shrink_to_fit();
//     return res;
//   }
//   int operator()(T x) { return LB(dat, x); }
// };

// template <typename T>
// struct Index_Compression_DISTINCT_LARGE {
//   vc<T> dat;
//   vc<int> build(vc<T> X) {
//     vc<int> I = argsort(X);
//     vc<int> res(len(X));
//     for (auto& i: I) { res[i] = len(dat), dat.eb(X[i]); }
//     dat.shrink_to_fit();
//     return res;
//   }
//   int operator()(T x) { return LB(dat, x); }
// };

// template <typename T, bool SMALL>
// using Index_Compression_DISTINCT =
//     typename std::conditional<SMALL, Index_Compression_DISTINCT_SMALL<T>,
//                               Index_Compression_DISTINCT_LARGE<T>>::type;
// template <typename T, bool SMALL>
// using Index_Compression_SAME =
//     typename std::conditional<SMALL, Index_Compression_SAME_SMALL<T>,
//                               Index_Compression_SAME_LARGE<T>>::type;

// // SAME: [2,3,2] -> [0,1,0]
// // DISTINCT: [2,2,3] -> [0,2,1]
// // (x): lower_bound(X,x) をかえす
// template <typename T, bool SAME, bool SMALL>
// using Index_Compression =
//     typename std::conditional<SAME, Index_Compression_SAME<T, SMALL>,
//                               Index_Compression_DISTINCT<T, SMALL>>::type;

package main

import (
	"fmt"
	"sort"
	"unsafe"
)

func main() {
	{
		arr := []int32{1449, 12, 12, 987}
		newArr, index := indexCompressionDistinctSmall(arr)
		fmt.Println(newArr)
		fmt.Println(index(1000000))
	}

	{
		arr := []int32{1449, 12, 12, 987}
		newArr, index := indexCompressionSameSmall(arr)
		fmt.Println(newArr)
		fmt.Println(index(1000000))
	}

	{
		arr := []int32{1449, 12, 12, 987, 1e9 + 7}
		newArr, index := indexCompressionDistinctLarge(arr)
		fmt.Println(newArr)
		fmt.Println(index(1000000))
		fmt.Println(index(1e9 + 7))
	}

	{
		arr := []int32{1449, 12, 12, 987, 1e9 + 7}
		newArr, index := indexCompressionSameLarge(arr)
		fmt.Println(newArr)
		fmt.Println(index(1000000))
		fmt.Println(index(1e9 + 7))
	}
}

type Num interface {
	int | int32 | int64 | uint | uint32 | uint64
}

// IndexCompression 用于对数组进行压缩.
//
//	 allowSame: 相同大小的元素编号是否能相同.
//		true -> [2,3,2] -> [0,1,0]
//		false -> [2,3,2] -> [0,2,0]
//	 smallRange: 数据极差较小(不超过1e7).
func IndexCompression[T Num](arr []T, allowSame bool, smallRange bool) (build func([]T) []int32, index func(T) int32) {
	return
}

func indexCompressionSameLarge[T Num](arr []T) (compressedArr []int32, index func(T) int32) {
	var data []T
	order := ArgSort(arr)
	compressedArr = make([]int32, len(arr))
	for _, v := range order {
		if len(data) == 0 || data[len(data)-1] != arr[v] {
			data = append(data, arr[v])
		}
		compressedArr[v] = int32(len(data) - 1)
	}
	data = data[:len(data):len(data)]

	index = func(x T) int32 {
		return int32(sort.Search(len(data), func(i int) bool { return data[i] >= x }))
	}

	return
}

func indexCompressionDistinctLarge[T Num](arr []T) (compressedArr []int32, index func(T) int32) {
	var data []T
	order := ArgSort(arr)
	compressedArr = make([]int32, len(arr))
	for _, v := range order {
		compressedArr[v] = int32(len(data))
		data = append(data, arr[v])
	}
	data = data[:len(data):len(data)]

	index = func(x T) int32 {
		return int32(sort.Search(len(data), func(i int) bool { return data[i] >= x }))
	}

	return
}

func indexCompressionSameSmall(arr []int32) (compressedArr []int32, index func(int) int32) {
	var min_, max_ int
	var data []int32

	compressedArr = append(arr[:0:0], arr...)
	min32, max32 := int32(0), int32(-1)
	if len(compressedArr) > 0 {
		for _, x := range compressedArr {
			if x < min32 {
				min32 = x
			}
			if x > max32 {
				max32 = x
			}
		}
	}
	data = make([]int32, max32-min32+2)
	for _, x := range compressedArr {
		data[x-min32+1] = 1
	}
	for i := 0; i < len(data)-1; i++ {
		data[i+1] += data[i]
	}
	for i, v := range compressedArr {
		compressedArr[i] = data[v-min32]
	}
	min_, max_ = int(min32), int(max32)

	index = func(x int) int32 {
		return data[clamp(x-min_, 0, max_-min_+1)]
	}

	return
}

func indexCompressionDistinctSmall(arr []int32) (compressedArr []int32, index func(int) int32) {
	var min_, max_ int
	var data []int32

	compressedArr = append(arr[:0:0], arr...)
	min32, max32 := int32(0), int32(-1)
	if len(compressedArr) > 0 {
		for _, x := range compressedArr {
			if x < min32 {
				min32 = x
			}
			if x > max32 {
				max32 = x
			}
		}
	}
	data = make([]int32, max32-min32+2)
	for _, x := range compressedArr {
		data[x-min32+1]++
	}
	for i := 0; i < len(data)-1; i++ {
		data[i+1] += data[i]
	}
	for i, x := range compressedArr {
		compressedArr[i] = data[x-min32]
		data[x-min32]++
	}
	copy(data[1:], data)
	data[0] = 0
	min_, max_ = int(min32), int(max32)

	index = func(x int) int32 {
		return data[clamp(x-min_, 0, max_-min_+1)]
	}

	return
}

func clamp[T Num](x, min, max T) T {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func ArgSort[T Num](nums []T) []int32 {
	order := make([]int32, len(nums))
	for i := int32(0); i < int32(len(order)); i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func cast[To, From any](v From) To {
	return *(*To)(unsafe.Pointer(&v))
}
