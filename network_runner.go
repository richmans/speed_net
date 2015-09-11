package main
import "fmt"
const NUM_WORKERS = 32
const QUEUE_SIZE = 40

type runRequest struct {
  network *Network
  start int
}

type networkRunner struct{
  networks map[int]*Network
  requests chan *runRequest
  responses chan int
}

func NetworkRunner() (*networkRunner){
  runner :=  &networkRunner{}
  runner.networks = make(map[int]*Network)
  runner.start_workers(NUM_WORKERS)
  return runner
}

func (n *networkRunner) addNetwork(networkId int, network *Network) {
  n.networks[networkId] = network
}

func (n *networkRunner) worker(requests chan *runRequest, 
  responses chan int, 
  worker_id int) {
  
  for request := range requests {

    request.network.run_neurons(request.start)
    responses <- request.start
  }
}

func (n *networkRunner) start_workers(numWorkers int) {

  n.requests = make(chan *runRequest, QUEUE_SIZE)
  n.responses = make(chan int, QUEUE_SIZE)
  for i:=0; i< numWorkers; i++ {
    fmt.Printf("Worker %d started\n", i)
    go n.worker(n.requests, n.responses, i)
  }
}

func (n *networkRunner) run_network(networkId int) {
  network, ok := n.networks[networkId]
  if !ok { return }
  network_size := network.network_size
  job_size := network.job_size
  num_requests := 0
  for i := 0;i < network_size; i += job_size{
    n.requests <- &runRequest{network, i}
    num_requests += 1
  }
  
  for ; num_requests > 0; {
    <- n.responses
    num_requests -= 1
  }
  
}