package defs

import (
	"flake"
	"sync"
)


type Id []byte

var flkInstance *flake.Flk
var lck = &sync.RWMutex{}

func GetFlake() *flake.Flk {
	var err error
	lck.Lock()
	if flkInstance == nil {	
		flkInstance,err = flake.FlakeNode()
		if err != nil {
			lck.Unlock()
			panic(err)
		}
	}
	lck.Unlock()
	
	return flkInstance
}

type Checkpoint struct {
	Id Id
	Value string
}



// a read only struct - should not be set
 // technically a crdt
type Claim struct {
	Id Id
	Start int64
	End int64
	WorkerId Id
	Complete bool
}


type Task struct {
	Id Id
	Payload string
	Namedq string
	Priority int
	Checkpoints []Checkpoint
	Claims [](*Claim)
}


func NewTask(payload string,namedq string,priority int) *Task {
	var err error
	t := new(Task)
	t.Payload = payload
	t.Id,err = GetFlake().NextId()

	if err != nil {
		panic(err)
	}
	
	t.Namedq = namedq
	t.Priority = priority
	return t
}

