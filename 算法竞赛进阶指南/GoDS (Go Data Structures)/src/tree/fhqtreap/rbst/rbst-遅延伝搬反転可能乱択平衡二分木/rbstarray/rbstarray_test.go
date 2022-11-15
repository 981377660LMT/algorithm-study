// An effective arraylist implemented by FHQTreap.
//
// Author:

// https://github.com/981377660LMT/algorithm-study

//

// Reference:

// https://baobaobear.github.io/post/20191215-fhq-treap/

// https://nyaannyaan.github.io/library/rbst/treap.hpp

// https://github.com/EndlessCheng/codeforces-go/blob/f9d97465d8b351af7536b5b6dac30b220ba1b913/copypasta/treap.go

package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_Update(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"Test_Update"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := make([]int, 1e5)
			for i := range nums {
				if i <= 5e4 {
					nums[i] = i + 5e4
				} else {
					nums[i] = 5e4
				}
			}

			k := int(5e4)
			initNums := make([]int, 1e5)
			// 更新:单点更新,查询:区间最大值
			treap := NewFHQTreap(initNums)
			time1 := time.Now()
			for _, num := range nums {
				preMax := treap.Query(max(0, num-k), num)
				treap.Update(num, num+1, preMax+1)
			}
			treap.QueryAll()
			fmt.Println(time.Since(time1))

		})
	}
}
