package main

import (
	"flag"
	"golang.org/x/xerrors"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"syscall"

	"sigs.k8s.io/kube-scheduler-simulator/simulator/scheduler"
	"sigs.k8s.io/kube-scheduler-simulator/simulator/scheduler/config"
)

// entry point.
func main() {
	if err := startScheduler(); err != nil {
		panic(err)
	}
}

func startScheduler() error {
	path := flag.String("kubeconfig", "/etc/kubernetes/scheduler.conf", "kubeconfig")
	flag.Parse()
	dsc, err := config.DefaultSchedulerConfig()
	if err != nil {
		return xerrors.Errorf("get defaultScheduler Config: %w", err)
	}
	// must be in cluster
	restCfg, err := clientcmd.BuildConfigFromFlags("", *path)
	if err != nil {
		panic(err.Error())
	}

	client, err := clientset.NewForConfig(restCfg)
	if err != nil {
		panic(err.Error())
	}

	svc := scheduler.NewSchedulerService(client, restCfg, dsc, false)
	if err = svc.StartScheduler(dsc); err != nil {
		return xerrors.Errorf("start scheduler: %w", err)
	}

	// wait the signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit

	return nil
}
