package pq

import "divi/defs"
import "sync"

// a simple but inefficient version of heap based priority queue

type Heap struct {
	arr [](*defs.Task)
	n int
	lock *sync.Mutex	
}



func NewHeap(initCapacity int) *Heap {
	b := new(Heap)
	b.arr = make([](*defs.Task),initCapacity)
	b.lock = &sync.Mutex{}
	return b
}


func (a *Heap) Deq() (*defs.Task,bool) {
	a.lock.Lock()
	if a.n == 0 {
		return nil,false
	}
	
	head := a.arr[1]
	// move the last element to head and re-heapify

	a.arr[1] = a.arr[a.n]
	a.sink(1)
	a.lock.Unlock()
	return head,true
}


func (a *Heap) sink(index int) {
	
	var child int

	for 2*index <=a.n {
		
		child = 2*index
		// is there a right child and is its priority greater than left?
		if (child+1) <= a.n && a.arr[child+1].Priority > a.arr[child].Priority {
			child = child + 1
		}

		// is the current index's priority greater than its children's priorities?
		if a.arr[index].Priority > a.arr[child].Priority {
			return
		} else {
			// swap and continue
			a.swap(index,child)
			index = child
		}		
	}
}

func (a *Heap) Enq(t *defs.Task) {
	a.lock.Lock()
	a.n++
	if a.n == len(a.arr) {
		// double
		a.doubleArr()
	}
	a.arr[a.n] = t

	// swim up
	a.swim(a.n)
	a.lock.Unlock()
}

func (a *Heap) doubleArr() {
	
	newarr := make([](*defs.Task),(len(a.arr)*2))

	// copy elements
	for i:=0;i<len(a.arr);i++ {
		newarr[i] = a.arr[i]
	}

	a.arr = newarr
}


func (a *Heap) swim(index int) {

	// swim up the element insert at arr[n]
	
	parent := index/2

	for parent >= 1 && a.arr[parent].Priority < a.arr[index].Priority {
		// swap the elements
		a.swap(parent,index)
		index = parent
		parent = index/2
	} 
}
func (a *Heap) swap(i int,j int) {
	hold := a.arr[i]
	a.arr[i] = a.arr[j]
	a.arr[j] = hold
}
