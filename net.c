
#define neuron_output(inputs, index) (*(net + index))
#define char_product(x,y) ((x < y) ? ((x + 1) * y >> 8) : ((y + 1) * x >> 8 ))
#define pull_down(i,w) i - char_product(i, (127 - w))
#define pull_up(i,w) i + char_product((255 - i), (w - 127)) 
#define weight(i,w) ((w < 127) ? pull_down(i,w) : pull_up(i,w))
#include <stdio.h>
#include "net.h"
void run_neurons(void* net_ptr, void* inputs_ptr, void* outputs_ptr, int start_index, int number, int network_size) {
  unsigned char* net = ((unsigned char*)net_ptr);
  unsigned char* inputs = ((unsigned char*)inputs_ptr);
  unsigned char* outputs = ((unsigned char*)outputs_ptr);
  int end_index = start_index + number;
  //printf("Let's do some thinking!\n");
  for(int idx = start_index; idx < end_index; idx++) {
    //printf("Neuron %u\n", idx);
    unsigned char node_weight = net[idx * NEURON_SIZE + WEIGHT_OFFSET];
    
    int sum = 0;
    for(int input = 0; input < NUMBER_INPUTS; input += 1) {
      signed char input_offset = net[idx * NEURON_SIZE + input];
      int input_index = (idx + input_offset) % (network_size);
      unsigned char input_value = inputs[input_index];
      //printf("Input index %d", input_index);
      unsigned char input_weight = net[idx * NEURON_SIZE + input + NUMBER_INPUTS];
      //printf("Input %u, weight %u, offset %u\n", input, input_weight, input_offset);
      sum += weight(input_value, input_weight);
    }
    unsigned char avg_input = ((unsigned char)(sum / NUMBER_INPUTS));
    unsigned char total = avg_input;
    unsigned char output = weight(total, node_weight);
    //printf("Neuron output %u\n", (unsigned char)sum);
    outputs[idx] = (unsigned char)output;
    
  }
}