package main


import "divi"
import "sync"
import "time"
import "fmt"

var enqWg sync.WaitGroup

func BenchmarkEnq() {

	start := time.Now()
	enqInParallel()
	end := time.Now()

	fmt.Println(end.Sub(start))


}

func enqInParallel() {

	for i := 1; i <= 2; i++ {
		enqWg.Add(1)

		go func() {
			e := divi.NewTaskEngine()
			for p := 1; p <= 1000; p++ {
				e.Enqueue("payload","q",p)
			}
			enqWg.Done()
		}()
	}
	enqWg.Wait()
}


func main() {
	BenchmarkEnq()
}
