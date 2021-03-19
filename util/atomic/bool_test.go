package atomic

import "testing"

func TestBool(t *testing.T) {
	var isOk Bool
	if isOk.Get() != false {
		t.Fatal("expect false")
	}

	isOk.Set(true)
	if isOk.Get() != true {
		t.Fatal("expect true")
	}

	isOk.Set(false)
	if isOk.Get() != false {
		t.Fatal("expect false")
	}
}
