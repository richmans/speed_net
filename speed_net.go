package main
import "os"
import "fmt"

func startController() {
  controller := new(neuralController)
  hostname, envHostname := os.LookupEnv("NEURAL_MASTER")
  if !envHostname {
    fmt.Printf("Please provide the hostname of the master in $NEURAL_MASTER\n")
    os.Exit(1)
  }
  controller.start(hostname)
}

func startPerformance() {
  performance()
}

func main(){ 
  if(len(os.Args) != 1 && os.Args[1] == "performance"){
    startPerformance()
  }else{
    startController()
  }
}