package main

/*
#include "net.h"
#include "net.c"
*/
import "C"
import "fmt"
import "unsafe"
import "math/rand"
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

func main() {
  number_nodes := 10  
  random_seed := 120101
  iterations := 1000
  //step_size := number_nodes / 1000
  var net_slice, net = create_buffer(number_nodes * C.NEURON_SIZE)
  var input_slice, inputs = create_buffer(number_nodes)
  var output_slice, outputs = create_buffer(number_nodes)
  fmt.Println("Initializing...")
  // randomize the network
  rand.Seed(int64(random_seed))
  randomize_buffer(net_slice)
  randomize_buffer(input_slice)
  fmt.Println("Now running")
  fmt.Println(input_slice)
  for i := 0; i < iterations; i++ {
    C.run_neurons(net, inputs, outputs, 0, C.int(number_nodes), C.int(number_nodes - 1))
    fmt.Println(output_slice)
    temp_ptr, temp_slice := inputs, input_slice
    inputs, input_slice = outputs, output_slice
    outputs, output_slice = temp_ptr, temp_slice
  }
  fmt.Println(output_slice)
 // fmt.Println(output_slice)
  fmt.Printf("All done...\n")
}