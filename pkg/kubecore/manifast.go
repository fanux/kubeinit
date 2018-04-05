package kubecore

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

//KubeletServiceStr  is
var KubeletServiceStr = `
[Unit]
Description=kubelet: The Kubernetes Node Agent
Documentation=http://kubernetes.io/docs/

[Service]
ExecStart=/usr/bin/kubelet
Restart=always
StartLimitInterval=0
RestartSec=10

[Install]
WantedBy=multi-user.target
`

var installOnlineYumSh = `
yum install -y git
yum install -y docker
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
yum install -y kubelet kubeadm kubectl
`

// also must gen config first
var installOnlineWgetShTmpl = `
wget {{.DownloadURL}}
tar zxvf kubernetes-node-linux-amd64.tar.gz
cp kubernetes/node/bin/kube* /usr/bin

cp ../out/kube/kubelet.service /etc/systemd/system/
mkdir /etc/systemd/system/kubelet.service.d
cp ../out/kube/10-kubeadm.conf /etc/systemd/system/kubelet.service.d
`

var installOffline = `
cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system

sysctl -w net.ipv4.ip_forward=1
systemctl stop firewalld && systemctl disable firewalld

swapoff -a
setenforce 0
docker load -i ../image/kube-core-images.tar
cp ../bin/kube* /usr/bin
cp ../out/kube/kubelet.service /etc/systemd/system/
mkdir /etc/systemd/system/kubelet.service.d
cp ../out/kube/10-kubeadm.conf /etc/systemd/system/kubelet.service.d
systemctl enable kubelet
systemctl enable docker
`

var runSh = `
# kubeadm init
kubeadm init --config ../out/kube/config
mkdir ~/.kube
cp /etc/kubernetes/admin.conf ~/.kube/config
`

var saveSh = `
cp /usr/bin/kube* bin
# TODO this will save all the images
docker save $(docker images|grep ago|awk '{print $1":"$2}') -o kube-core-images.tar
mv kube-core-images.tar image/
`
