package divi

import "divi/defs"
import "divi/apperror"
import "errors"


// Enqueue does not cuurently support - task dependencies
// but the idea is that you can pass a list of task ids as parameters
// The current task will only picked up for execution - if all of its
// dependencies are marked as "HasEnded" = true

// we need to support streaming instead of clients polling later
// think about how it can be done later




type TaskServer interface {
	Enqueue(payload string,namedq string,priority int) (defs.Id, error) 
	Delete(taskId defs.Id) error
	Poll(taskId defs.Id) (bool,error)
	Get(taskId defs.Id) (defs.Task,error)
}


type Coordinator interface {
	Stake(workerId defs.Id,namedq string) (defs.Task,defs.Claim,error)
	Claim(taskId defs.Id,claimId defs.Id,complete bool) (defs.Claim,error)
	Checkpoint(taskId defs.Id,payload string) error	
}


// what is most efficient way to store the task map?
// do we do a two level task?
type TaskEngine struct {
	qStore *Store
	taskMap TaskMap
	ch *ClaimHandler
}


func NewTaskEngine() *TaskEngine {
	t := new(TaskEngine)
	t.qStore = GetStore()
	t.taskMap = GetTaskMap()
	t.ch = NewClaimHandler()
	return t
}




func (t *TaskEngine) Stake(workerId defs.Id,namedq string) (defs.Task,defs.Claim,error) {
	
	// look for the next available task in namedq to be allocated to the worker
	q := t.qStore.GetQueue(namedq)
	task,ok := q.Deq()

	if !ok {
		return defs.Task{},defs.Claim{},apperror.New(errors.New("No task to allocate"),apperror.NoTaskToClaim)
	}

	// mark the task as allocated
	// the way we allocate a task is - when we create a cliam for it
	// if the claim expires, the task needs to go back to the priority queue to be reclaimed
	
	claim := t.ch.Stake(task,workerId)

	return *task,*claim,nil
}


func (t *TaskEngine) Claim(taskId defs.Id,holdingClaimId defs.Id,complete bool) (defs.Claim,error) {
	
	task,ok := t.taskMap.Get(taskId)
	if !ok {
		return defs.Claim{},apperror.New(errors.New("Task not found"),apperror.TaskNotFound)
	}

	
	newClaim,err :=  t.ch.Claim(task,holdingClaimId,complete)
	if err != nil {
		return defs.Claim{},err
	} else {
		return *newClaim,err
	}
}

func (t *TaskEngine) Enqueue(payload string,namedq string,priority int) (defs.Id,error) {

		
	task := defs.NewTask(payload,namedq,priority)
	// get the queue from store
	q := t.qStore.GetQueue(namedq)
	
	q.Enq(task)
	t.taskMap.Put(task)
	
	return task.Id,nil
}

func (t *TaskEngine) Delete(taskId defs.Id) error {
	panic("Not implemented")	
}


func (t *TaskEngine) Poll(taskId defs.Id) (bool,error) {
	
	task,ok := t.taskMap.Get(taskId)

	if !ok {
		return false,apperror.New(errors.New("Task not found"),apperror.TaskNotFound)
	}

	if task.Claims == nil {
		return false,nil
	} else {
		return task.Claims[len(task.Claims)-1].Complete,nil
	}
}

func (t *TaskEngine) Get(taskId defs.Id) (defs.Task,error) {
	task,ok := t.taskMap.Get(taskId)

	if !ok {
		return defs.Task{},apperror.New(errors.New("Task not found"),apperror.TaskNotFound)
	}

	return (*task),nil
}


