package split

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	got := Split("我爱你", "爱")
	want := []string{"我", "你"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want:%v,got:%v", want, got)
	}
}

func TestSplit_Group(t *testing.T) {
	type test []struct {
		input string
		sep   string
		want  []string
	}
	tests := test{
		{input: "a:b:c", sep: ":", want: []string{"a", "b", "c"}},
		{input: "a:b:c", sep: ",", want: []string{"a:b:c"}},
		{input: "abcdefg", sep: "cde", want: []string{"ab", "fg"}},
		{input: "沙河有沙又有河", sep: "沙", want: []string{"", "河有", "又有河"}},
	}
	for _, tc := range tests {
		got := Split(tc.input, tc.sep)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("excepted:%v,got:%v", tc.want, got)
		}
	}
}

func TestSplit_Sub(t *testing.T) {
	type test struct {
		input string
		sep   string
		want  []string
	}
	tests := map[string]test{
		"simple":     {input: "我爱你", sep: "爱", want: []string{"我", "你"}},
		"multi sep":  {input: "a:b:c", sep: ",", want: []string{"a:b:c"}},
		"multi sep2": {input: "abcdefg", sep: "cde", want: []string{"ab", "fg"}},
		"chinese":    {input: "沙河有沙又有河", sep: "沙", want: []string{"", "河有", "又有河"}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := Split(tc.input, tc.sep)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("excepted:%v,got:%v", tc.want, got)
			}
		})

	}
}

func BenchmarkSplit(b *testing.B) {
	//b.N不是固定的数
	for i := 0; i < b.N; i++ {
		Split("沙河有沙又有河", "沙")
	}
}
