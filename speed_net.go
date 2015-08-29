package main

/*
#include "net.h"
#include "net.c"
*/
import "C"
import "fmt"
import "unsafe"
import "math/rand"
import "runtime"
type Network struct {
  net unsafe.Pointer
  inputs unsafe.Pointer
  outputs unsafe.Pointer
  job_size int
  network_size int
}

func create_buffer(buffer_size int) ([]byte, unsafe.Pointer) {
  var buffer []byte
  buffer = make([]byte, buffer_size)
  var buffer_pointer = unsafe.Pointer(&(buffer[0]))
  return buffer, buffer_pointer
}

func randomize_buffer(slice []byte) {
  for idx, _ := range slice {
    slice[idx] = byte(rand.Intn(256))
  }
}

func worker(requests chan int, 
  responses chan int, 
  network *Network,
  worker_id int) {
  
    for request := range requests {
      C.run_neurons(network.net, network.inputs, network.outputs, C.int(request), C.int(network.job_size) , C.int(network.network_size))
      responses <- request
    }
}

func start_workers(n int, network *Network) (chan int, chan int){
  q_size := (network.network_size / network.job_size) + 10
  requests := make(chan int, q_size)
  responses := make(chan int, q_size)
  for i:=0; i< n; i++ {
    fmt.Printf("Worker %d started\n", i)
    go worker(requests, responses, network, i)
  }
  return requests, responses
}

func run_network(requests chan int, responses chan int, job_size int, network_size int) {
  num_requests := 0
  for i := 0;i < network_size; i+= job_size{
    requests <- i
    num_requests += 1
  }
  
  for ; num_requests > 0; {
    <- responses
    num_requests -= 1
  }
  
}
func main() {
  number_nodes := 1000000
  random_seed := 120101
  iterations := 1000
  job_size := number_nodes / 10
  num_workers := 8
  fmt.Println("Initializing...")
  runtime.GOMAXPROCS(num_workers + 1)
  var net_slice, net = create_buffer(number_nodes * C.NEURON_SIZE)
  var input_slice, inputs = create_buffer(number_nodes)
  var output_slice, outputs = create_buffer(number_nodes)
  
  // randomize the network
  rand.Seed(int64(random_seed))
  randomize_buffer(net_slice)
  randomize_buffer(input_slice)
  
  network := &Network{
    net: net,
    inputs: inputs,
    outputs: outputs,
    job_size: job_size,
    network_size: number_nodes,
  }
  
  requests, responses := start_workers(num_workers, network)
  
  fmt.Println("Now running")
  fmt.Println(input_slice[0:20])
  for i := 0; i < iterations; i++ {
    
    //swap the input and output buffers
    run_network(requests, responses, job_size, number_nodes)
    temp_ptr, temp_slice := inputs, input_slice
    inputs, input_slice = outputs, output_slice
    outputs, output_slice = temp_ptr, temp_slice
  }
  fmt.Println(output_slice[0:20])
  fmt.Printf("All done...\n")
}