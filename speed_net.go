package main
import "os"
import "fmt"
import "runtime"
import "time"
func startPerformance() {
  iterations := 10000
  num_workers := 16
  number_nodes := 100000
  job_size := number_nodes / num_workers
  fmt.Println("Initializing...")
  runtime.GOMAXPROCS(num_workers * 2 + 10)
  network := RandomNetwork(number_nodes, job_size)
  runner := NetworkRunner()
  runner.addNetwork(1, network)
  fmt.Println("Now running")
 
  t0 := time.Now()
  for i := 0; i < iterations; i++ {
    network.set_inputs([]byte{43,112,231,195})
    runner.run_network(1)
    network.swap()
    
  }
  t1 := time.Now()
  fmt.Println(network.get_outputs())
  fmt.Printf("All done...\n")
  fmt.Printf("The whole thing took %v to run.\n", t1.Sub(t0))
}


func startController() {
  controller := new(neuralController)
  hostname, envHostname := os.LookupEnv("NEURAL_MASTER")
  if !envHostname {
    fmt.Printf("Please provide the hostname of the master in $NEURAL_MASTER\n")
    os.Exit(1)
  }
  controller.start(hostname)
}

func main(){ 
  if(len(os.Args) != 1 && os.Args[1] == "performance"){
    startPerformance()
  }else{
    startController()
  }
}