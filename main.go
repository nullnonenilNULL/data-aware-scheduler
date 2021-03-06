package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/johscheuer/data-aware-scheduler/databackend"
	"github.com/johscheuer/data-aware-scheduler/databackend/quobyte"
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/rest"
	"k8s.io/client-go/1.5/tools/clientcmd"
)

var (
	schedulerConfigPath = flag.String("config", "./config.yaml", "absolute path to the scheduler config file")
)

const schedulerName = "data-aware-scheduler"

func main() {
	log.Println("Starting data-aware-scheduler...")

	flag.Parse()
	schedulerConfig := readConfig(*schedulerConfigPath)

	var cfg *rest.Config
	var err error
	if schedulerConfig.InCluster {
		log.Println("Starting data-aware-scheduler in Cluster")
		cfg, err = rest.InClusterConfig()
	} else {
		log.Println("Starting data-aware-scheduler out of Cluster")
		cfg, err = clientcmd.BuildConfigFromFlags("", schedulerConfig.Kubeconfig)
	}
	if err != nil {
		panic(err.Error())
	}

	clientset := kubernetes.NewForConfigOrDie(cfg)
	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	var backend databackend.DataBackend
	if schedulerConfig.Backend == "quobyte" {
		log.Println("Starting data-aware-scheduler with Quobyte backend")
		backend = quobyte.NewQuobyteBackend(
			schedulerConfig.Opts,
			clientset,
		)
	}

	processor := newProcessor(clientset, doneChan, &wg, backend)
	wg.Add(1)
	go processor.monitorUnscheduledPods()

	wg.Add(1)
	go processor.reconcileUnscheduledPods(30)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}
}
