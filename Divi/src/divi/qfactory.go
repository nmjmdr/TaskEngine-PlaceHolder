package divi

import "divi/pq"

type QueueType int

const (
	ArrayHeap QueueType = iota
)


func GetQueue(qType QueueType) pq.Queue {
	

	a :=  pq.NewHeap(100)
	
	return a
}



