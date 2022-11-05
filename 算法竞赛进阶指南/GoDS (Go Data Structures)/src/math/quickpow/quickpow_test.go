// 快速幂

package quickpow

import "testing"

func TestPow(t *testing.T) {
	type args struct {
		base int64
		exp  int64
		mod  int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: "test", args: args{base: 2, exp: 123456789, mod: 1e9 + 7}, want: 178116276},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pow(tt.args.base, tt.args.exp, tt.args.mod); got != tt.want {
				t.Errorf("Pow() = %v, want %v", got, tt.want)
			}
		})

	}
}
