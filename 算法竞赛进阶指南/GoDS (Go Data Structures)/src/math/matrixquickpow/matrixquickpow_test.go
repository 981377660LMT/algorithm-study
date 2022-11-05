// 矩阵快速幂

// https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/math_matrix.go#L70-L117

package matrixquickpow

import (
	"reflect"
	"testing"
)

func TestMatrix_Pow(t *testing.T) {
	type args struct {
		exp int64
		mod int64
	}
	tests := []struct {
		name   string
		matrix Matrix
		args   args
		want   Matrix
	}{
		// TODO: Add test cases.
		{
			name: "test1", matrix: Matrix{{1, 1, 1}, {1, 0, 0}, {0, 1, 0}}, args: args{876543210987654318, 1e9 + 7},
			want: Matrix{{177180593, 842485564, 442632457}, {442632457, 734548143, 399853107}, {399853107, 42779350, 334695036}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.matrix.Pow(tt.args.exp, tt.args.mod); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Matrix.Pow() = %v, want %v", got, tt.want)
			}
		})
	}
}
