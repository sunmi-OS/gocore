package smap

import (
	"log"
	"testing"
)

type SmapTest struct {
	Name string
	Age  int
}

func TestSmap(t *testing.T) {
	sm := Map[string, *SmapTest]{}

	actual, loaded := sm.LoadOrStore("test", &SmapTest{Name: "test", Age: 10})
	if !loaded {
		log.Printf("LoadOrStore not have and store: %v", actual)
	}
	log.Printf("actual: %v", actual)
	value, ok := sm.Load("test")
	if !ok {
		log.Printf("Load not have")
	}
	log.Printf("value: %v", value)
	sm.Store("test2", &SmapTest{Name: "test2", Age: 11})
	l := sm.Len()
	log.Printf("len = 2 ?: %d", l)
	sm.Store("test2", &SmapTest{Name: "test2", Age: 12})
	l = sm.Len()
	log.Printf("len = 2 ?: %d", l)

	// ==================================

	andDelete, ok := sm.LoadAndDelete("test")
	if !ok {
		log.Printf("LoadAndDelete not have")
		return
	}
	log.Printf("andDelete: %v", andDelete)
	l = sm.Len()
	log.Printf("len = 1 ?: %d", l)
	sm.Store("test3", &SmapTest{Name: "test3", Age: 13})
	l = sm.Len()
	log.Printf("len = 2 ?: %d", l)
	_, ok = sm.Load("test")
	log.Printf("after load and delete load sm[test] is %v", ok)

	sm.Store("test2", &SmapTest{Name: "test2", Age: 20})

	v2, ok := sm.Load("test2")
	if !ok {
		log.Printf("sm[test2] Load not have")
		return
	}
	log.Printf("sm[test2] value: %v", v2)
	l = sm.Len()
	log.Printf("len = 2 ?: %d", l)
	sm.Delete("test2")
	l = sm.Len()
	log.Printf("len = 1 ?: %d", l)
	sm.Delete("test3")
	l = sm.Len()
	log.Printf("len = 0 ?: %d", l)

}
