package examples

import "testing"

func TestAlwaysFails(t *testing.T) {
	t.Log("This test always fails to demonstrate running covalyzer-go on a failing test.")
	t.Fail()
}
