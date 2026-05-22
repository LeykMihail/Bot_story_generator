package router

import "testing"

func TestOutboundQueueIndex(t *testing.T) {
	t.Parallel()
	const numWorkers = 16

	tests := []struct {
		name   string
		userID int64
		want   int
	}{
		{"positive", 64, 0},
		{"positive non-zero", 65, 1},
		{"negative user id", -1, 15},
		{"large negative", -17, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := outboundQueueIndex(tt.userID, numWorkers)
			if got != tt.want {
				t.Fatalf("outboundQueueIndex(%d, %d) = %d, want %d", tt.userID, numWorkers, got, tt.want)
			}
			if got < 0 || got >= numWorkers {
				t.Fatalf("index %d out of range [0, %d)", got, numWorkers)
			}
		})
	}
}
