package macnodetest

import (
	"flake"
	"testing"
)


func BenchmarkEnq(b *testing.B) {

	for i:=0;i<b.N;i++ {
		macNodeAssign(b)
	}

}

func macNodeAssign(b *testing.B) {
	
	m := new(flake.MacNodeId)
	_, err:= m.Id()
	if err != nil {
		b.Fatal(err)
	}
}
