package pqtest

import "divi"
import "divi/defs"
import "testing"
import "math/rand"
import "time"

func TestEnqAndDeq(t *testing.T) {

	q := divi.GetQueue(divi.ArrayHeap)

	rand.Seed(time.Now().UnixNano())
	for i:=0;i<100;i++ {
		t := defs.NewTask("payload","q",rand.Intn(100))
		q.Enq(t)
	}

	for i:=0;i<100;i++ {	
		_,flag := q.Deq()
		if flag != true {
			t.Fatal("No value dequeued")
		}
	}

}
