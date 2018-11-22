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
		points []point
		want   float64
	}{
		{
			points: []point{},
			want:   0,
		},
		{
			points: []point{
				point{0, 0, 0},
			},
			want: 0,
		},
		{
			points: []point{
				point{0, 0, 0},
				point{0, 1, 0},
				point{0, 1, 1},
				point{0, 4, 5},
			},
			want: 7,
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
