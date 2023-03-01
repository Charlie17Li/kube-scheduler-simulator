package main

import (
	"golang.org/x/xerrors"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kube-scheduler-simulator/simulator/config"
	"sigs.k8s.io/kube-scheduler-simulator/simulator/scheduler"
)

// entry point.
func main() {
	if err := startScheduler(); err != nil {
		panic(err)
	}
}

func startScheduler() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return xerrors.Errorf("get config: %w", err)
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

	svc := scheduler.NewSchedulerService(client, restCfg, cfg.InitialSchedulerCfg, false)
	if err = svc.StartScheduler(cfg.InitialSchedulerCfg); err != nil {
		return xerrors.Errorf("start scheduler: %w", err)
	}

	return nil
}
