package singleflight

import "testing"

func TestGroup_Do(t *testing.T) {
	var g Group
	v, err := g.Do("key", func() (i interface{}, err error) {
		return "bar", nil
	})

	if v != "bar" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}
