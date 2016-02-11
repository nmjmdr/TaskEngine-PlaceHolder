package divi

//import "defs"
import "divi/pq"
import "sync"

type Store struct {
	qMap map[string]pq.Queue
	mutex *sync.Mutex
}

var qs *Store = nil
var mutex = &sync.Mutex{}

func GetStore() * Store {
	
	mutex.Lock()
	if qs == nil {
		qs = new(Store)
		qs.qMap = make(map[string]pq.Queue)
		qs.mutex = &sync.Mutex{}
		// start collecting here
	}
	mutex.Unlock()
	return qs
}

func (qs *Store) collect() {
	
}

func (qs *Store) GetQueue(namedq string) pq.Queue {
	
	qs.mutex.Lock()
	q,ok := qs.qMap[namedq]	
	if !ok {
		q = GetQueue(ArrayHeap)
		qs.qMap[namedq] = q
	}
	qs.mutex.Unlock()
	
	return q
}
