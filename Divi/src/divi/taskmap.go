package divi

import "divi/defs"
import "sync"

type TaskMap interface {
	Put(t *defs.Task)
	Get(taskId defs.Id) (*defs.Task,bool)
	Delete(taskId defs.Id)
}



type SimpleMap struct {
	m map[string](*defs.Task)
	rwLock *sync.RWMutex
}

func GetTaskMap() TaskMap {
	var tm TaskMap
	tm = newSimpleMap()
	return tm
}


func newSimpleMap() *SimpleMap {
	s := new(SimpleMap)
	s.m = make(map[string](*defs.Task))
	s.rwLock = &sync.RWMutex{}
	return s
}

func (s *SimpleMap) Put(t *defs.Task) {
	s.rwLock.Lock()
	s.m[string(t.Id)] = t
	s.rwLock.Unlock()
}

func (s *SimpleMap) Get(taskId defs.Id) (*defs.Task,bool) {
	s.rwLock.RLock()
	t,ok := s.m[string(taskId)]
	s.rwLock.RUnlock()
	return t,ok
}

func (s *SimpleMap) Delete(taskId defs.Id) {
	delete(s.m,string(taskId))
}
