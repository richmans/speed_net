package main

import "math/rand"

const WEIGHT_OFFSET = 10
const NUMBER_INPUTS = 5
const NEURON_SIZE = 1 + 2 * NUMBER_INPUTS

type Network struct {
  network_size int
  inputs []byte
  outputs []byte
  job_size int
  input_nodes []int
  output_nodes []int
  net []byte
}

func RandomNetwork(number_nodes int, job_size int) (*Network) {
  random_seed := 120101
  number_inputs := 4
  number_outputs := 4
  var net = make([]byte, number_nodes * NEURON_SIZE)
  var inputs = make([]byte, number_nodes)
  var outputs = make([]byte, number_nodes)
  var input_nodes = make([]int, number_inputs)
  var output_nodes = make([]int, number_outputs)
  
  // randomize the network
  rand.Seed(int64(random_seed))
  randomize_buffer(net)
  randomize_buffer(inputs)
  randomize_nodes(input_nodes, number_nodes)
  randomize_nodes(output_nodes, number_nodes)
  network := &Network{
    net: net,
    inputs: inputs,
    outputs: outputs,
    job_size: job_size,
    network_size: number_nodes,
    input_nodes: input_nodes,
    output_nodes: output_nodes,
  }
  network.set_inputs([]byte{43,112,231,195})
  return network
}

func (n *Network) get_outputs() ([]byte) {
  result := make([]byte, len(n.output_nodes))
  for idx, _ := range n.output_nodes {
    result[idx] = n.outputs[n.output_nodes[idx]]
  }
  return result
}

func (network *Network) swap() {
  temp_ptr := network.inputs
  network.inputs = network.outputs
  network.outputs = temp_ptr
}

func (n *Network) set_inputs(inputs []byte) {
  for idx, _ := range n.input_nodes {
    n.inputs[n.input_nodes[idx]] = inputs[idx]
  }
}

func (n *Network) calculate_weight(value int, weight int) (int){
  if weight < 127 {
    return value - value * ((127.0 - weight) / 255) 
  } else {
    return value + (255 - value) * ((weight - 127) / 255)
  }
}

func (network *Network) run_neurons(start_index int) {
  end_index := start_index + network.job_size
  if (end_index > network.network_size) {
    end_index = network.network_size
  }
  //printf("Let's do some thinking!\n");
  for idx := start_index; idx < end_index; idx++ {
    //printf("Neuron %u\n", idx);
    node_weight := network.net[idx * NEURON_SIZE + WEIGHT_OFFSET]
    
    sum := 0
    for input := 0; input < NUMBER_INPUTS; input += 1 {
      input_offset := network.net[idx * NEURON_SIZE + input]
      input_index := (idx + int(input_offset)) % network.network_size
      input_value := network.inputs[input_index]
      input_weight := network.net[idx * NEURON_SIZE + input + NUMBER_INPUTS];
      //printf("Input %u, weight %u, offset %u\n", input, input_weight, input_offset);
      sum += network.calculate_weight(int(input_value), int(input_weight));
    }
    total := sum / NUMBER_INPUTS;
    output := network.calculate_weight(total, int(node_weight));
    network.outputs[idx] = byte(output);
  }
}


func randomize_buffer(slice []byte) {
  for idx, _ := range slice {
    slice[idx] = byte(rand.Intn(256))
  }
}

func randomize_nodes(slice []int, max int) {
  for idx, _ := range slice {
    slice[idx] = rand.Intn(max)
  }
}


