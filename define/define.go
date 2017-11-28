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

//KubeletSystemdTemp is
var KubeletSystemdTemp = `[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"
Environment="KUBELET_NETWORK_ARGS=--network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin"
Environment="KUBELET_DNS_ARGS=--cluster-dns=10.96.0.10 --cluster-domain=cluster.local"
Environment="KUBELET_AUTHZ_ARGS=--authorization-mode=Webhook --client-ca-file=/etc/kubernetes/pki/ca.crt"
Environment="KUBELET_CADVISOR_ARGS=--cadvisor-port=0"
Environment="KUBELET_CGROUP_ARGS=--cgroup-driver={{.CgroupDriver}}"
Environment="KUBELET_CERTIFICATE_ARGS=--rotate-certificates=true --cert-dir=/var/lib/kubelet/pki"
ExecStart=
ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_SYSTEM_PODS_ARGS $KUBELET_NETWORK_ARGS $KUBELET_DNS_ARGS $KUBELET_AUTHZ_ARGS $KUBELET_CADVISOR_ARGS $KUBELET_CGROUP_ARGS $KUBELET_CERTIFICATE_ARGS $KUBELET_EXTRA_ARGS
`

//define is
var (
	Nodes             []string
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
