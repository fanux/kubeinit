package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fanux/kubeinit/define"
)

var installComposeAndsshPass = `
chmod +x bin/docker-compose
cp bin/docker-compose /usr/bin
cp bin/sshpass /usr/bin
`

var remoteDockerConfig = `
echo OPTIONS=\"-H 0.0.0.0:2375 -H unix:///var/run/docker.sock --selinux-enabled --log-driver=journald --signature-verification=false\" >> /etc/sysconfig/docker || true
systemctl enable docker
systemctl restart docker
`

var sshEnable = `
cat <<EOF > ~/.ssh/config
Host *
  StrictHostKeyChecking no
EOF
`

var initbasesh = `
cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system

swapoff -a
setenforce 0 || true
systemctl stop firewarlld || true
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

var applyDashboardAndAdmin = `
kubectl apply -f out/kubernetes-dashboard.yaml
kubectl apply -f out/dashboard-admin.yaml
`

var applyHeapsters = `
kubectl create -f out/deploy/kube-config/influxdb/
kubectl create -f out/deploy/kube-config/rbac/heapster-rbac.yaml
`

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
func getCurrentIP() string {
	shell := `grep server /etc/kubernetes/admin.conf | awk -F "//" '{print $2}' | awk -F ":" '{print $1}'`
	return applyShellOutput(shell)
}

func changeConfigFileIPs(ip, dip string) {
	sh := fmt.Sprintf("sed -i 's/%s/%s/g' ", ip, dip)
	dir := fmt.Sprintf("/tmp/%s/kubernetes", dip)

	for _, file := range []string{dir + "/manifests/kube-apiserver.yaml", dir + "/kubelet.conf", dir + "/admin.conf", dir + "/controller-manager.conf", dir + "/scheduler.conf"} {
		applyShell(sh + file)
	}
}

func sendFileToDstMaster(ip string) {
	sh := fmt.Sprintf("docker -H %s:2375 run --name %s -v /etc:/etc -v /usr/bin:/usr/bin -v /etc/systemd/system:/etc/systemd/system -v /etc/systemd/system/kubelet.service.d:/etc/systemd/system/kubelet.service.d busybox /bin/test", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp /tmp/%s/kubernetes %s:/etc", ip, ip, ip)
	applyShell(sh)

	sh = fmt.Sprintf("docker -H %s:2375 cp bin/kubectl %s:/usr/bin ", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp bin/kubelet %s:/usr/bin ", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp bin/kubeadm %s:/usr/bin ", ip, ip)
	applyShell(sh)

	sh = fmt.Sprintf("docker -H %s:2375 cp out/kubelet.service %s:/etc/systemd/system ", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp out/10-kubeadm.conf %s:/etc/systemd/system/kubelet.service.d", ip, ip)
	applyShell(sh)

	execSSHCommand(define.User, define.Password, ip, initbasesh)
	execSSHCommand(define.User, define.Password, ip, startKubelet)
}

func sendFileToDstNode(ip string) {
	sh := fmt.Sprintf("docker -H %s:2375 run --name %s-node  -v /usr/bin:/usr/bin -v /etc/systemd/system:/etc/systemd/system -v /etc/systemd/system/kubelet.service.d:/etc/systemd/system/kubelet.service.d busybox /bin/test", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp bin/kubectl %s-node:/usr/bin ", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp bin/kubelet %s-node:/usr/bin ", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp bin/kubeadm %s-node:/usr/bin ", ip, ip)
	applyShell(sh)

	sh = fmt.Sprintf("docker -H %s:2375 cp out/kubelet.service %s-node:/etc/systemd/system ", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp out/10-kubeadm.conf %s-node:/etc/systemd/system/kubelet.service.d", ip, ip)
	applyShell(sh)
}

func distributeFiles() {
	ip := getCurrentIP()
	i := strings.Index(ip, "\n")
	ip = ip[:i]
	for _, masterip := range define.KubeFlags.MasterIPs {
		loadRemoteImage(masterip)
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

		//go sendFileToDstMaster(masterip)
		sendFileToDstMaster(masterip)
	}
}

func execSSHCommand(user, passwd, ip, sh string) {
	//rsh := fmt.Sprintf("bash -c %s", sh)
	fmt.Println("exec ssh command: ", sh)
	cmd := exec.Command("sshpass", "-p", passwd, "ssh", user+"@"+ip, sh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("exec ssh command error: %s", err)
	}
}

func applyLoadBalance(ip string) {
	//docker cp haproxy.cfg to remote host
	sh := fmt.Sprintf("docker -H %s:2375 run --name %s-ha -v /etc/haproxy:/etc/haproxy busybox /bin/test", ip, ip)
	applyShell(sh)
	sh = fmt.Sprintf("docker -H %s:2375 cp out/haproxy.cfg %s-ha:/etc/haproxy", ip, ip)
	applyShell(sh)

	//start haproxy container
	sh = fmt.Sprintf("docker -H %s:2375 run --restart=always --net=host -v /etc/haproxy:/usr/local/etc/haproxy --name ha -d haproxy:1.7 ", ip)
	applyShell(sh)
}

func changeTOLBIPPort(cmd string) string {
	for _, masterip := range define.KubeFlags.MasterIPs {
		if strings.Contains(cmd, masterip) {
			return strings.Replace(cmd, masterip+":6443", define.KubeFlags.LoadbalanceIP+":"+define.KubeFlags.LoadbalancePort, -1)
		}
	}
	fmt.Println("Error: change LoadbalanceIP failed: ", define.KubeFlags.LoadbalanceIP+":"+define.KubeFlags.LoadbalancePort)
	return cmd
}

func applyDashboard() {
	applyShell(applyDashboardAndAdmin)
}

func applyHeapster() {
	// kubectl create -f deploy/kube-config/influxdb/
	//kubectl create -f deploy/kube-config/rbac/heapster-rbac.yaml
	applyShell(applyHeapsters)
}

func loadRemoteImage(ip string) {
	sh := fmt.Sprintf("docker -H %s:2375 load -i image/images.tar", ip)
	applyShell(sh)
}

//Apply is
func Apply() {
	applyShell(sshEnable)
	applyShell(installComposeAndsshPass)
	LoadKubeinitConfig()

	if define.InitBaseEnvironment {
		applyShell(initbasesh)
		applyShell(cpBinAndConfigs)
		applyShell(loadDockerImages)
	}

	if define.StartEtcdCluster {
		for i, ip := range define.KubeFlags.EtcdIPs {
			execSSHCommand(define.User, define.Password, ip, remoteDockerConfig)
			loadRemoteImage(ip)
			sh := fmt.Sprintf(startEtcdCluster, ip, i)
			applyShell(sh)
		}
	}

	if define.InitKubeadm {
		s := applyShellOutput(initKubeadm)
		fmt.Println(s)
		if define.Distribute {
			distributeFiles()
		}
		i := strings.Index(s, "kubeadm join")
		s1 := s[i:]
		j := strings.Index(s1, "\n")
		joinCmd := s1[:j+1]
		joinCmd = changeTOLBIPPort(joinCmd)
		fmt.Println("join Cmd is: ", joinCmd)

		execSSHCommand(define.User, define.Password, define.KubeFlags.LoadbalanceIP, remoteDockerConfig)
		loadRemoteImage(define.KubeFlags.LoadbalanceIP)
		applyLoadBalance(define.KubeFlags.LoadbalanceIP)
		//apply join commands
		for _, ip := range define.KubeFlags.NodeIPs {
			/*
				go func(ip string) {
					//send files to node
					sendFileToDstNode(ip)
					execSSHCommand(define.User, define.Password, ip, initbasesh)
					execSSHCommand(define.User, define.Password, ip, joinCmd)
				}(ip)
			*/
			execSSHCommand(define.User, define.Password, ip, remoteDockerConfig)
			loadRemoteImage(ip)
			sendFileToDstNode(ip)
			execSSHCommand(define.User, define.Password, ip, initbasesh)
			execSSHCommand(define.User, define.Password, ip, joinCmd)
		}
	}
	applyHeapster()
	applyDashboard()
	//var wait chan int
	//<-wait
	//set .kube/config
	//TODO kubectl config set-cluster kubernetes --server=https://47.52.227.242:6444 --kubeconfig=$HOME/.kube/config
	cmd := fmt.Sprintf("kubectl config set-cluster kubernetes --server=https://%s:%s --kubeconfig=$HOME/.kube/config", define.KubeFlags.LoadbalanceIP, define.KubeFlags.LoadbalancePort)
	applyShell(cmd)
}
