package unionfindarray

import (
	"reflect"
	"testing"
)

func TestNewUnionFindArray(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want *UnionFindArray
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				n: 10,
			},
			want: &UnionFindArray{
				Part:   10,
				Rank:   []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				size:   10,
				parent: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnionFindArray(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnionFindArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_demo(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			demo()
		})
	}
}
