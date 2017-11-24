package dataframe

import (
	"testing"
	"time"
)

func TestTimeValue(t *testing.T) {
	vt := time.Now()
	value := NewTimeValue(vt)
	if v, ok := value.Time(TimeDefaultLayout); !ok || !v.Equal(vt) {
		t.Fatalf("expected time %q, got %q", vt, v)
	}
}
