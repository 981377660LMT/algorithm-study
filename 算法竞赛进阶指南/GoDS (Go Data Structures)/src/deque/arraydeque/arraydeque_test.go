// !用两个 slice 头对头拼在一起实现(每个slice头部只删除元素，不添加元素，互相弥补劣势)

// https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e

package arraydeque

import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{name: "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
