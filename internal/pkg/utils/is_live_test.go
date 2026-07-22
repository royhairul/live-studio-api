package utils

import (
	"testing"
	"time"
)

func TestIsLive(t *testing.T) {
	now := time.Now()
	ms := func(t time.Time) int64 { return t.UnixMilli() }

	tests := []struct {
		name      string
		startTime int64
		duration  int64
		want      bool
	}{
		{
			// A running stream: Shopee's duration advances with it, so
			// start+duration stays level with now.
			name:      "ongoing session ending now",
			startTime: ms(now.Add(-2 * time.Hour)),
			duration:  (2 * time.Hour).Milliseconds(),
			want:      true,
		},
		{
			name:      "just started",
			startTime: ms(now.Add(-10 * time.Second)),
			duration:  (10 * time.Second).Milliseconds(),
			want:      true,
		},
		{
			// Poll lag: reported end is slightly stale but within tolerance.
			name:      "ended 1 minute ago, within tolerance",
			startTime: ms(now.Add(-61 * time.Minute)),
			duration:  time.Hour.Milliseconds(),
			want:      true,
		},
		{
			name:      "ended 10 minutes ago, past tolerance",
			startTime: ms(now.Add(-70 * time.Minute)),
			duration:  time.Hour.Milliseconds(),
			want:      false,
		},
		{
			// The real case observed from Shopee: a 10-minute stream that
			// started 5 hours ago still appears in realtime/sessionList.
			name:      "finished hours ago",
			startTime: ms(now.Add(-5 * time.Hour)),
			duration:  (10 * time.Minute).Milliseconds(),
			want:      false,
		},
		{
			name:      "missing start time",
			startTime: 0,
			duration:  time.Hour.Milliseconds(),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLive(tt.startTime, tt.duration); got != tt.want {
				t.Errorf("IsLive(%d, %d) = %v, want %v", tt.startTime, tt.duration, got, tt.want)
			}
		})
	}
}
