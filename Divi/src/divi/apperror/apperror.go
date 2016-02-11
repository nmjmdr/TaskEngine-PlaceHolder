package apperror

import "fmt"

const (

	TaskNotFound = 1
	
	// task deletion related errors
	TaskAlreadyComplete = 100

	NoTaskToClaim
	WrongContextTaskNotYetClaimed
	TaskAllocatedToAnother
	ClaimAlreadyExpied
)

type Err struct {
	Err error
	Code int
}

func (e *Err) Error() string {
	return fmt.Sprintf("%s, error code: %d",e.Err.Error(),e.Code)
}

func (e *Err) ErrorCode() int {
	return e.Code
}

func New(e error,c int) *Err {
	aErr := new(Err)
	aErr.Err = e
	aErr.Code = c
	return aErr
}
