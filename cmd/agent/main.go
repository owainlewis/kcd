package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/golang/glog"
	exec "github.com/owainlewis/kcd/pkg/exec"
)

var kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig file")

func main() {

	flag.Parse()

	config, err := getKubeConfig(*kubeconfig)
	if err != nil {
		return
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("Failed to build client: %s", err)
	}

	executor := exec.NewExecutor(config, client)

	stdOut, stdErr, err := executor.Command("default", "kcd-gcr8d", "/bin/bash", "-c", "apt-get install -y emacs")

	fmt.Printf("Result: %s %s", stdOut, stdErr)

	if err != nil {
		return
	}
}

// Build will construct a Kubernetes clientset from a given
// config path. If the path is empty then it will default to use
// an in-cluster configuration
func Build(configpath string) (*kubernetes.Clientset, error) {
	config, err := getKubeConfig(configpath)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getKubeConfig(configpath string) (*rest.Config, error) {
	if configpath != "" {
		return clientcmd.BuildConfigFromFlags("", configpath)
	}

	return rest.InClusterConfig()
}
