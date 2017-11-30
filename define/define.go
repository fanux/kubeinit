package define

//EtcdComposeTemp is
var EtcdComposeTemp = `version: '2.1'
services:
    etcd{{.Index}}:
        container_name: etcd_infra{{.Index}}
        image: {{.Image}}
        command: |
                etcd --name infra{{.Index}}
                --initial-advertise-peer-urls http://{{.EndPoint}}:2380
                --listen-peer-urls http://{{.EndPoint}}:2380
                --listen-client-urls http://{{.EndPoint}}:2379,http://127.0.0.1:2379
                --advertise-client-urls http://{{.EndPoint}}:2379
                --data-dir /etcd-data.etcd
                --initial-cluster-token etcd-cluster-1
                --initial-cluster {{.EndPoints}}
                --initial-cluster-state new
        restart: always
        volumes:
           - /data/etcd-data.etcd:/etcd-data.etcd
        network_mode: "host"

`

//EtcdComposeTempST is
type EtcdComposeTempST struct {
	Index     string
	Image     string
	EndPoint  string
	EndPoints string
}

//KubeadmTemp is
var KubeadmTemp = `apiVersion: kubeadm.k8s.io/v1alpha1
kind: MasterConfiguration
apiServerCertSANs:
{{range .APIServerCertSANs}}- {{.}} {{"\n"}}{{end}}
etcd:
  endpoints:
{{range .EtcdEndPoints}}  - http://{{.}}:2379  {{"\n"}}{{end}}
networking:
  podSubnet: {{.PodSubnet}}
kubernetesVersion: {{.KubernetesVersion}}
`

//KubeadmTempST is
type KubeadmTempST struct {
	APIServerCertSANs []string
	EtcdEndPoints     []string
	PodSubnet         string
	KubernetesVersion string
}

//KubeletSystemdTemp is
var KubeletSystemdTemp = `[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"
Environment="KUBELET_NETWORK_ARGS=--network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin"
Environment="KUBELET_DNS_ARGS=--cluster-dns=10.96.0.10 --cluster-domain=cluster.local"
Environment="KUBELET_AUTHZ_ARGS=--authorization-mode=Webhook --client-ca-file=/etc/kubernetes/pki/ca.crt"
Environment="KUBELET_CADVISOR_ARGS=--cadvisor-port=0"
Environment="KUBELET_CGROUP_ARGS=--cgroup-driver={{.}}"
Environment="KUBELET_CERTIFICATE_ARGS=--rotate-certificates=true --cert-dir=/var/lib/kubelet/pki"
ExecStart=
ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_SYSTEM_PODS_ARGS $KUBELET_NETWORK_ARGS $KUBELET_DNS_ARGS $KUBELET_AUTHZ_ARGS $KUBELET_CADVISOR_ARGS $KUBELET_CGROUP_ARGS $KUBELET_CERTIFICATE_ARGS $KUBELET_EXTRA_ARGS
`

//HaproxyTemp  is
var HaproxyTemp = `global
  daemon
  log 127.0.0.1 local0
  log 127.0.0.1 local1 notice
  maxconn 4096

defaults
  log               global
  retries           3
  maxconn           2000
  timeout connect   5s
  timeout client    50s
  timeout server    50s

frontend k8s
  bind *:{{.LoadbalancePort}}
  mode tcp
  default_backend k8s-backend

backend k8s-backend
  balance roundrobin
  mode tcp
{{range $index, $element := .BackendEndPoint}}  server k8s-{{$index}} {{$element}}:6443 check {{"\n"}}{{end}}
`

//HaproxyTempST is
type HaproxyTempST struct {
	LoadbalancePort string
	BackendEndPoint []string
}

//Flags is
type Flags struct {
	MasterIPs         []string
	EtcdIPs           []string
	MasterEndPoints   []string
	APIServerCertSANs []string
	PodSubnet         string
	KubernetesVersion string
	LoadbalancePort   string
	LoadbalanceIP     string
	Apply             bool
	EtcdImage         string
	Version           string
	Subnet            string

	// render out compose file and kubeadm config
	ConfigOutDir string
	// /etc/kubernetes files, this need copy to other nodes, and change ips
	KubernetesDir string

	InitBaseEnvironment bool
	InitKubeadm         bool
	InitOtherMasters    bool
	StartEtcdCluster    bool
}

//KubeFlags is
var KubeFlags Flags
