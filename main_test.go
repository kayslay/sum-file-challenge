package main

import (
	"testing"
)

func Benchmark_main(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run()
	}
}

func Test_concatenate(t *testing.T) {
	type args struct {
		last int
		b    []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"last",
			args{0, []byte("120")},
			120,
		},
		{
			"last",
			args{0, []byte("120")},
			120,
		},
		{
			"last",
			args{1, []byte("120")},
			1120,
		},
		{
			"last",
			args{3, []byte("33")},
			333,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := concatenate(tt.args.last, tt.args.b); got != tt.want {
				t.Errorf("concatenate() = %v, want %v", got, tt.want)
			}
		})
	}
}
