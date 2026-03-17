package search

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    []float32
		b    []float32
		want float32
	}{
		{
			name: "identical vectors",
			a:    []float32{1.0, 0.0, 0.0},
			b:    []float32{1.0, 0.0, 0.0},
			want: 1.0,
		},
		{
			name: "orthogonal vectors",
			a:    []float32{1.0, 0.0, 0.0},
			b:    []float32{0.0, 1.0, 0.0},
			want: 0.0,
		},
		{
			name: "opposite vectors",
			a:    []float32{1.0, 0.0, 0.0},
			b:    []float32{-1.0, 0.0, 0.0},
			want: -1.0,
		},
		{
			name: "45 degree angle",
			a:    []float32{1.0, 0.0, 0.0},
			b:    []float32{1.0, 1.0, 0.0},
			want: float32(1.0 / math.Sqrt(2)),
		},
		{
			name: "different lengths",
			a:    []float32{1.0, 0.0},
			b:    []float32{1.0, 0.0, 0.0},
			want: 0.0,
		},
		{
			name: "zero vectors",
			a:    []float32{0.0, 0.0, 0.0},
			b:    []float32{0.0, 0.0, 0.0},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosineSimilarity(tt.a, tt.b)
			if math.Abs(float64(got-tt.want)) > 0.0001 {
				t.Errorf("cosineSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPositionAwareOrder(t *testing.T) {
	results := []Result{
		{DocID: "a", Score: 0.9},
		{DocID: "b", Score: 0.8},
		{DocID: "c", Score: 0.7},
		{DocID: "d", Score: 0.6},
		{DocID: "e", Score: 0.5},
	}

	ordered := positionAwareOrder(results)

	if len(ordered) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(ordered))
	}

	if len(ordered) >= 2 {
		if ordered[0].Score < ordered[1].Score {
			t.Error("first two results should be highest scoring")
		}
	}
}

func TestPositionAwareOrderSmall(t *testing.T) {
	tests := []struct {
		name    string
		results []Result
		wantLen int
	}{
		{"empty", []Result{}, 0},
		{"single", []Result{{DocID: "a", Score: 0.9}}, 1},
		{"two", []Result{{DocID: "a", Score: 0.9}, {DocID: "b", Score: 0.8}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ordered := positionAwareOrder(tt.results)
			if len(ordered) != tt.wantLen {
				t.Errorf("expected %d results, got %d", tt.wantLen, len(ordered))
			}
		})
	}
}
