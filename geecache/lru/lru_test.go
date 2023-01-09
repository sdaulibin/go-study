package lru

import (
	"fmt"
	"reflect"
	"testing"
)

type String string

func (d String)Len() int {
	return len(d)
}

func TestGet(t *testing.T)  {
	lru := New(int64(0),nil)
	lru.Add("key1",String("1234"))
	if v,ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("geecache hit key1=1234 failed")
	}
	if _,ok:= lru.Get("key2");ok {
		t.Fatalf("geecache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1,k2,k3 := "key1","key2","key3"
	v1,v2,v3 := "value1","value2","value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap),nil)
	lru.Add(k1,String(v1))
	lru.Add(k2,String(v2))
	lru.Add(k3,String(v3))

	fmt.Println(lru.Len(),cap)

	if _,ok :=  lru.Get(k1);ok || lru.Len() !=2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string,0)
	lru := New(int64(10),func(key string,value Value) {
		keys = append(keys,key)
	})
	lru.Add("key1",String("1"))
	lru.Add("key2",String("2"))
	lru.Add("key3",String("3"))
	lru.Add("key4",String("4"))

	expact := []string{"key1","key2"}

	if !reflect.DeepEqual(expact,keys) {
		t.Fatalf("call OnEvicted failed,expect keys equals to %s",expact)
	}
}