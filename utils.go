package main

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeHandle() *kubernetes.Clientset {
	var conf *rest.Config
	conf, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		fatal(fmt.Sprintf("error in getting Kubeconfig: %v", err))
	}

	cs, err := kubernetes.NewForConfig(conf)
	if err != nil {
		fatal(fmt.Sprintf("error in getting clientset from Kubeconfig: %v", err))
	}

	return cs
}

func fatal(msg string) {
	os.Stderr.WriteString(msg + "\n")
	os.Exit(1)
}
