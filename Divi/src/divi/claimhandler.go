package divi



import "divi/defs"
import "divi/apperror"
import "time"
import "errors"
import "reflect"


const milli = 1000000
const sec = 1000 * milli
const min = 1000 * 60
const ClaimExpiryMin = 2 * min

type ClaimHandler struct {	
}

func NewClaimHandler() *ClaimHandler {
	ch := new(ClaimHandler)
	return ch
}

func newClaim(workerId defs.Id,complete bool) *defs.Claim {
	// create a new claim
	var err error
	claim := new(defs.Claim)
	claim.Id,err = defs.GetFlake().NextId()
	if err != nil {
		panic(err)
	}
	claim.Start = time.Now().UnixNano()
	// now here is the problem - the classic problem
	// of computer clock being reset - or time moving back
	// how do we address it?
	claim.End = claim.Start + ClaimExpiryMin
	claim.WorkerId = workerId
	claim.Complete = complete
	return claim
}


func (c *ClaimHandler) Stake(task *defs.Task,workerId defs.Id) *defs.Claim {
	
	// staking a new claim
	task.Claims = append(task.Claims,newClaim(workerId,false))
	return task.Claims[len(task.Claims)-1]
}

func hasExpired(claim *defs.Claim) bool {
	
	// Major isue - this means the servers (leader and follower clocks now need to 
	// be in sync, because we are now measuring time!!!
	// can we do it differently?? - anything in RAFT??

	// now here is the problem - the classic problem
	// of computer clock being reset - or time moving back
	// how do we address it?

	return (claim.End - time.Now().UnixNano() < 0) 
}

func (c *ClaimHandler) Claim(task *defs.Task,holdingClaimId defs.Id,complete bool) (*defs.Claim,error) {

	// check if the previous claim has expired
	// if yes then it cannot be claimed
	// there are two scenarios here:
	// 1. The claim has expired and it has already been claimed by another worker
	// 2. The claim has expired and it has not yet been claimed by any other worker
	if task.Claims == nil {
		return nil,apperror.New(errors.New("Task has not been claimed at all"),apperror.WrongContextTaskNotYetClaimed)
	}

	lastClaim := task.Claims[len(task.Claims)-1]

	if reflect.DeepEqual(lastClaim.Id,holdingClaimId) {
		// somebody else hs been allocated the task
		return nil,apperror.New(errors.New("Task allocated to another work - cannot be claimed"),apperror.TaskAllocatedToAnother)
	}


	if hasExpired(lastClaim) {
		return nil,apperror.New(errors.New("Claim has already expired, cannot be renewed"),apperror.ClaimAlreadyExpied)
	}

	// ok, it can be renewed
	// at some point we will have to merge the claims
	// remove older claims
	task.Claims = append(task.Claims,newClaim(lastClaim.WorkerId,complete))
	return task.Claims[len(task.Claims)-1],nil

}
