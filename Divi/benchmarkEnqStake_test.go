package divitest

import "divi"
import "testing"
import "sync"
import "divi/apperror"
import "fmt"
import "divi/defs"
import "time"
import "math/rand"

func BenchmarkEnq(b *testing.B) {

	for i:=0;i<b.N;i++ {
		enqInParallel(b)
	}

}

func enqInParallel(b *testing.B) {

	var wg sync.WaitGroup

	for i := 1; i <= 200; i++ {
		wg.Add(1)

		func() {
			e := divi.NewTaskEngine()
			rand.Seed(time.Now().UnixNano())
			for p := 1; p <= 100; p++ {
				_,err := e.Enqueue("payload","q",rand.Intn(100))
				if err != nil {
					b.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func enqAndClaimInParallel(b *testing.B) {

	var ewg sync.WaitGroup
	var cwg sync.WaitGroup

	for i := 1; i <= 200; i++ {
		ewg.Add(1)

		func() {
			e := divi.NewTaskEngine()
			rand.Seed(time.Now().UnixNano())
			for p := 1; p <= 100; p++ {
				_,err := e.Enqueue("payload","q",rand.Intn(100))
				if err != nil {
					b.Fatal(err)
				}
			}
			ewg.Done()
		}()
	}
	

	for i := 1; i <= 10; i++ {
		cwg.Add(1)

		func(i int) {
			e := divi.NewTaskEngine()
			for p := 1; p <= 200; p++ {
				for {
					_,_,err := e.Stake(defs.Id(fmt.Sprintf("%d",i)),"q")
					if err != nil {
						aerr,ok := err.(*apperror.Err)
						if !ok {
							b.Fatal("Got a mistyped error")
						}
						if aerr.Code == apperror.NoTaskToClaim {
							continue
						} else {
							b.Fatal(err)
						}
					} else {
						break
					}

				}
			}
			cwg.Done()
		}(i)
	}



	ewg.Wait()
	cwg.Wait()
}

func BenchmarkEnqAndClaim(b *testing.B) {

	for i:=0;i<b.N;i++ {
		enqAndClaimInParallel(b)
	}	
}

