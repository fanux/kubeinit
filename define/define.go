package define

//EtcdComposeTemp is
var EtcdComposeTemp = `version: '2.1'
services:
{{$element := index .Hosts .Index}}
    etcd{{.Index}}:
        container_name: etcd_infra{{.Index}}
        image: {{.Image}}
        command: |
                etcd --name infra{{.Index}}
                --initial-advertise-peer-urls http://{{$element}}:2380
                --listen-peer-urls http://{{$element}}:2380
                --listen-client-urls http://{{$element}}:2379,http://127.0.0.1:2379
                --advertise-client-urls http://{{$element}}:2379
                --data-dir /etcd-data.etcd
                --initial-cluster-token etcd-cluster-1
                --initial-cluster {{range $index, $element := .Hosts}}{{if $index}},{{end}}infra{{$index}}=http://{{$element}}:2380{{end}}
                --initial-cluster-state new
        restart: always
        volumes:
           - /data/etcd-data.etcd:/etcd-data.etcd
        network_mode: "host"

`

//KubeadmTemp is
var KubeadmTemp = `apiVersion: kubeadm.k8s.io/v1alpha1
kind: MasterConfiguration
apiServerCertSANs:
- 10.1.245.93
- 10.1.245.94
- 10.1.245.95
- 47.52.227.242
etcd:
  endpoints:
  - http://10.1.245.94:2379
networking:
  podSubnet: 192.168.0.0/16
kubernetesVersion: v1.8.2
`

//define is
var (
	EtcdEndPoints     []string
	APIServerCertSANs []string
	PodSubnet         string
	KubernetesVersion string
)

//EtcdTemp is
type EtcdTemp struct {
	ServiceName   string
	ContainerName string
	Image         string
	IfraName      string
	EndPoint      string
	EndPoints     string
}
