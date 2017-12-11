package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

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

var startKubelet = `
systemctl enable kubelet
systemctl start kubelet
`

var loadDockerImages = `
docker load -i image/images.tar
`

var startEtcdCluster = "docker-compose -H %s:2375 -f out/etcd-docker-compose-%d.yml up -d"

var initKubeadm = `
kubeadm init --config out/kubeadm.yaml
mkdir -p $HOME/.kube
rm -f $HOME/.kube/config
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

func applyShellOutput(sh string) string {
	s, err := exec.Command("bash", "-c", sh).Output()
	if err != nil {
		fmt.Println("exec shell failed: ", sh)
		return ""
	}
	return string(s)
}

func getCurrentIP() string {
	shell := `grep server /etc/kubernetes/admin.conf | awk -F "//" '{print $2}' | awk -F ":" '{print $1}'`
	return applyShellOutput(shell)
}

func changeConfigFileIPs(ip, dip string) {
	sh := fmt.Sprintf("sed -i 's/%s/%s/g' ", ip, dip)
	dir := fmt.Sprintf("/tmp/%s", dip)

	for _, file := range []string{dir + "/manifests/kube-apiserver.yaml", dir + "/kubelet.conf", dir + "./admin.conf", dir + "./controller-manager.conf", dir + "./scheduler.conf"} {
		applyShell(sh + file)
	}
}

func sendFileToDstNode(ip string) {
	sh := fmt.Sprintf("docker -H %s:2375 run --name %s -v /etc/kubernetes:/etc/kubernetes -v /usr/bin:/usr/bin -v /etc/systemd/system:/etc/systemd/system -v /etc/systemd/system/kubelet.service.d:/etc/systemd/system/kubelet.service.d busybox sleep 36000", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp /tmp/%s/etc/kubernetes %s:/etc/kubernetes ", ip, ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp bin/kube* %s:/usr/bin ", ip, ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp out/kubelet.service %s:/etc/systemd/system ", ip, ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp out/10-kubeadm.conf %s:/etc/systemd/system/kubelet.service.d", ip, ip, ip)
	applyShell(sh)

	//load images
	sh = fmt.Sprintf("docker -H %s:2375 load image/images.tar", ip)
	applyShell(sh)

	execSSHCommand(define.User, define.Password, ip, initbasesh)
	execSSHCommand(define.User, define.Password, ip, startKubelet)
}

func distributeFiles() {
	ip := getCurrentIP()
	for _, masterip := range define.KubeFlags.MasterIPs {
		if masterip == ip {
			continue
		}
		dir := fmt.Sprintf("/tmp/%s", masterip)
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}

		sh := fmt.Sprintf("cp -r /etc/kubernetes %s", dir)
		applyShell(sh)

		// change the currentIP to masterip
		changeConfigFileIPs(ip, masterip)

		go sendFileToDstNode(masterip)
	}
}

func execSSHCommand(user, passwd, ip, sh string) {
	cmd := exec.Command("sshpass", "-p", ip, "ssh", user+"@"+passwd, "bash", "-c", sh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

//Apply is
func Apply() {
	LoadKubeinitConfig()

	if define.InitBaseEnvironment {
		applyShell(initbasesh)
		applyShell(cpBinAndConfigs)
		applyShell(loadDockerImages)
	}

	if define.StartEtcdCluster {
		for i, ip := range define.KubeFlags.EtcdIPs {
			sh := fmt.Sprintf(startEtcdCluster, ip, i)
			fmt.Println("apply etcd: ", sh)
			applyShell(sh)
		}
	}

	if define.InitKubeadm {
		s := applyShellOutput(initKubeadm)
		fmt.Println(s)
		i := strings.Index(s, "kubeadm join")
		s1 := s[i:]
		j := strings.Index(s1, "\n")
		joinCmd := s1[:j+1]
		fmt.Println("join Cmd is: ", joinCmd)
		//apply join commands
		for _, ip := range define.KubeFlags.NodeIPs {
			go func(ip string) {
				execSSHCommand(define.User, define.Password, ip, initbasesh)
				execSSHCommand(define.User, define.Password, ip, joinCmd)
			}(ip)
		}
	}

	if define.Distribute {
		distributeFiles()
	}
}
