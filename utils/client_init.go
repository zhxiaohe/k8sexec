package utils

import (
	"io/ioutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func K8sInit() *kubernetes.Clientset {
	//kubectl -n kubernetes-dashboard describe secret $(kubectl -n kubernetes-dashboard get secret | grep kubernetes-dashboard-admin | awk '{print $1}')

	// fmt.Println("server start")
	// config, _ := clientcmd.BuildConfigFromFlags("", "/Users/admin/.kube/config")
	config, _ := clientcmd.BuildConfigFromFlags("", "../conf/k8s.conf")
	// creates the clientset
	clientset, _ := kubernetes.NewForConfig(config)
	return clientset
}

func GetRestConf() (restConf *rest.Config) {
	var (
		kubeconfig []byte
	)

	// 读kubeconfig文件
	kubeconfig, err := ioutil.ReadFile("../conf/k8s.conf")
	if err != nil {

	}
	// 生成rest client配置
	restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {

	}
	return restConf
}
