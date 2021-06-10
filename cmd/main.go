package main

import (
	"context"
	"encoding/json"
	"fmt"
	"k8sexec/service"
	"k8sexec/utils"
	"runtime"

	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// v1 "k8s.io/client-go/1.5/pkg/api/v1"
// "k8s.io/client-go/kubernetes"
// "k8s.io/client-go/tools/clientcmd"
// "github.com/gin-gonic/gin"

func main() {
	fmt.Println("emmm.")
	service.Svc()
	maxProces := runtime.NumCPU()
	fmt.Println(maxProces)
}

func main2() {

	var (
		clientset *kubernetes.Clientset
		podsList  *core_v1.PodList
		err       error
	)
	clientset = utils.K8sInit()
	podsList, err = clientset.CoreV1().Pods("ingress-nginx1").List(context.TODO(), meta_v1.ListOptions{})

	if err != nil {
		fmt.Println("err:", err)
	}
	pod, err := json.Marshal(podsList)
	fmt.Println(string(pod))

}

func main3() {
	clientset := utils.K8sInit()
	// podsList, err := clientset.CoreV1().Pods("ingress-nginx1").List(context.TODO(), meta_v1.ListOptions{})
	// restclient = clientset.CoreV1().RESTClient()
	// clientSet.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec")
	// req := restclient.Post().
	// 		Resource("pods").
	// 		Name(t.pod).
	// 		Namespace(t.namespace).
	// 		SubResource("exec").
	// 		Param("container", t.container).
	// 		Param("stdin", "true").
	// 		Param("stdout", "true").
	// 		Param("stderr", "true").
	// 		Param("command", cmd).Param("tty", "true")
	// 	req.VersionedParams(
	// 		&v1.PodExecOptions{
	// 			Container: t.container,
	// 			Command:   []string{},
	// 			Stdin:     true,
	// 			Stdout:    true,
	// 			Stderr:    true,
	// 			TTY:       true,
	// 		},
	// 		scheme.ParameterCodec,
	// 	)

	// req := clientSet.CoreV1().RESTClient().Post().
	// 	Resource("pods").
	// 	Name(podName).
	// 	Namespace(namespace).
	// 	SubResource("exec")

	// req.VersionedParams(&v1.PodExecOptions{
	// 	Container: containerName,
	// 	Command:   cmd,
	// 	Stdin:     true,
	// 	Stdout:    true,
	// 	Stderr:    true,
	// 	TTY:       true,
	// }, scheme.ParameterCodec)

	podName := "ingress-nginx1-controller-b6xcs"
	podNs := "ingress-nginx1"
	containerName := "ingress-nginx1-controller-b6xcs"
	var restConf *rest.Config
	restConf = utils.GetRestConf()
	Req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(podNs).
		SubResource("exec").
		VersionedParams(&core_v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	fmt.Println(Req.URL())
	executor, err := remotecommand.NewSPDYExecutor(restConf, "POST", Req.URL())
	if err != nil {
		fmt.Println("END")
	}
	fmt.Println("executor")
	return

	fmt.Println(executor)
}
