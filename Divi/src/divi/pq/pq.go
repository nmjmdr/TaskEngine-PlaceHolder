package pq

import "divi/defs"


type Queue interface {
	Enq(t *defs.Task)
	Deq() (*defs.Task,bool)
}





