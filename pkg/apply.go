package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/fanux/kubeinit/define"
)

var initbasesh = `
cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system

swapoff -a
setenforce 0
systemctl stop firewarlld
`

var cpBinAndConfigs = `
cp bin/kube* /usr/bin
cp out/kubelet.service /etc/systemd/system/
mkdir /etc/systemd/system/kubelet.service.d
cp out/10-kubeadm.conf /etc/systemd/system/kubelet.service.d
systemctl enable kubelet
systemctl enable docker
`

var loadDockerImages = `
docker load -i images.tar
`

var startEtcdCluster = "docker-compose -H %s:2375 -f out/etcd-docker-compose-%d.yml up -d"

var initKubeadm = `
kubeadm init --config out/kubeadm.yaml
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
kubectl apply -f out/calico.yaml
`

//WriteFile is
func WriteFile(fileName string, content string) {
	b := []byte(content)
	err := ioutil.WriteFile(fileName, b, 0644)
	if err != nil {
		fmt.Println("write file error", err)
	}
}

//LoadKubeinitConfig is
func LoadKubeinitConfig() {
	fileName := fmt.Sprintf("out/kubeinit.json")
	input, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return // exit the function on error
	}
	defer input.Close()
	json.NewDecoder(input).Decode(&define.KubeFlags)
}

func applyShell(sh string) {
	cmd := exec.Command("bash", "-c", sh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

//Apply is
func Apply() {
	LoadKubeinitConfig()

	if define.KubeFlags.InitBaseEnvironment {
		applyShell(initbasesh)
		applyShell(cpBinAndConfigs)
		applyShell(loadDockerImages)
	}

	if define.KubeFlags.StartEtcdCluster {
		for i, ip := range define.KubeFlags.EtcdIPs {
			sh := fmt.Sprintf(startEtcdCluster, ip, i)
			fmt.Println("apply etcd: ", sh)
			applyShell(sh)
		}
	}

	if define.KubeFlags.InitKubeadm {
		applyShell(initKubeadm)
	}
}
