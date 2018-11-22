package main

import (
	"testing"
)

func Test_dist(t *testing.T) {
	type args struct {
		a point
		b point
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			args: args{
				a: point{x: 3, y: 2},
				b: point{x: 3, y: 2},
			},
			want: 0,
		},
		{
			args: args{
				a: point{x: 0, y: 0},
				b: point{x: 3, y: 4},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dist(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("dist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_path_distance(t *testing.T) {

	tests := []struct {
		name   string
		points []*point
		want   float64
	}{
		{
			points: []*point{},
			want:   0,
		},
		{
			points: []*point{
				newPoint(0, 0, 0),
			},
			want: 0,
		},
		{
			points: []*point{
				newPoint(0, 0, 0),
				newPoint(1, 1, 0),
				newPoint(2, 1, 1),
				newPoint(3, 0, 1),
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &path{
				points: tt.points,
			}
			if got := p.distance(); got != tt.want {
				t.Errorf("path.distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPrime(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{n: 1, want: true},
		{n: 5, want: true},
		{n: 12, want: false},
		{n: 17, want: true},
		{n: 19, want: true},
		{n: 123, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrime(tt.n); got != tt.want {
				t.Errorf("isPrime() = %v, want %v", got, tt.want)
			}
		})
	}
}
