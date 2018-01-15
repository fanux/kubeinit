# install docker
yum install -y docker
systemctl enable docker
systemctl start docker

# install kubeadm
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
setenforce 0
yum install -y kubelet kubeadm kubectl
systemctl enable kubelet && systemctl start kubelet

# disable selinux
setenforce 0

# iptable set
cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system

# gen config files
TODO

# init
kubeadm init --config out/kubeadm.yaml
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config

# install calico
wget https://docs.projectcalico.org/v2.6/getting-started/kubernetes/installation/hosted/kubeadm/1.6/calico.yaml
rm out/calico.yaml
mv calico.yaml out
kubectl apply -f out/calico.yaml

# taint master
kubectl taint nodes --all node-role.kubernetes.io/master-

# install heapster
git clone https://github.com/kubernetes/heapster
mv deploy out/
rm -rf heapster
kubectl create -f out/deploy/kube-config/influxdb/
kubectl create -f out/deploy/kube-config/rbac/heapster-rbac.yaml

# install dashboard
wget https://raw.githubusercontent.com/kubernetes/dashboard/v1.8.1/src/deploy/recommended/kubernetes-dashboard.yaml
rm out/kubernetes-dashboard.yaml
mv kubernetes-dashboard.yaml out
kubectl apply -f out/kubernetes-dashboard.yaml
