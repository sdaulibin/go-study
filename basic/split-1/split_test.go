package split

import "testing"
import "reflect"

func TestSplit(t *testing.T) {
	got := Split("我爱你", "爱")
	want := []string{"我", "你"}
}
