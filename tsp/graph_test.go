package tsp

import (
	"testing"
)

func Test_dist(t *testing.T) {
	type args struct {
		a Point
		b Point
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			args: args{
				a: Point{X: 3, Y: 2},
				b: Point{X: 3, Y: 2},
			},
			want: 0,
		},
		{
			args: args{
				a: Point{X: 0, Y: 0},
				b: Point{X: 3, Y: 4},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dist(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("dist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_path_distance(t *testing.T) {

	tests := []struct {
		name   string
		points []*Point
		want   float64
	}{
		{
			points: []*Point{},
			want:   0,
		},
		{
			points: []*Point{
				NewPoint(0, 0, 0),
			},
			want: 0,
		},
		{
			points: []*Point{
				NewPoint(0, 0, 0),
				NewPoint(1, 1, 0),
				NewPoint(2, 1, 1),
				NewPoint(3, 0, 1),
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Path{
				Points: tt.points,
			}
			if got := p.Distance(); got != tt.want {
				t.Errorf("path.Distance() = %v, want %v", got, tt.want)
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
			if got := IsPrime(tt.n); got != tt.want {
				t.Errorf("isPrime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_path_twoOptSwap(t *testing.T) {

	type args struct {
		i int
		k int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: args{
				i: 1,
				k: 3,
			},
			want: []int{0, 3, 2, 1, 4, 5},
		},
		{
			args: args{
				i: 2,
				k: 4,
			},
			want: []int{0, 1, 4, 3, 2, 5},
		},
		{
			args: args{
				i: 1,
				k: 4,
			},
			want: []int{0, 4, 3, 2, 1, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPath()
			p.AddPoint(NewPoint(0, 0, 0))
			p.AddPoint(NewPoint(1, 0, 0))
			p.AddPoint(NewPoint(2, 0, 0))
			p.AddPoint(NewPoint(3, 0, 0))
			p.AddPoint(NewPoint(4, 0, 0))
			p.AddPoint(NewPoint(5, 0, 0))
			p.Swap(tt.args.i, tt.args.k)

			for i, id := range tt.want {
				if p.Points[i].ID != id {
					t.Errorf("id = %v, want = %v", id, p.Points[i].ID)
				}
			}
		})
	}
}
