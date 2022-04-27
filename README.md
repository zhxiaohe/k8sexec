# k8sexec

"kubectl exec -it" webterm

1. 前端使用 xterm.js 简单实现
2. 后端使用 gin + gorilla/websocket + k8s.io/client-go

![]("./doc/1623234872272.jpg")



`tar zxf k8sexec-master.zip && cd k8sexec-master`

`go mod tidy`

`go mod vendor`

`修改 conf/k8s.conf`

`修改 service/view.go  下func ExecShell`

       ` podName := "nginx-569bf69fb7-5gr8w"`

	`podNs := "default"`

	`containerName := "nginx"`

`cd cmd && go run main.go`

![image](https://user-images.githubusercontent.com/11907113/165434889-0ae3d18d-154a-4eee-8866-3c8dacf363ee.png)
