package main

import (
	"golang.org/x/xerrors"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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
	dsc, err := config.DefaultSchedulerConfig()
	if err != nil {
		return xerrors.Errorf("get defaultScheduler Config: %w", err)
	}
	// must be in cluster
	restCfg, err := clientcmd.BuildConfigFromFlags("", "")
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

	return nil
}
