package kubecore

import (
	"fmt"
	"html/template"
	"os/exec"
	"strings"

	"github.com/fanux/kubeinit/pkg"
)

//Kubecore is
type Kubecore struct{}

//Info is
func (e *Kubecore) Info() (string, string) {
	return "kubecore", "kubecore include kubelet, apiserver manager scheduler and kubeproxy"
}

func genKubeAdmConfigFile(etcdIPs []string, masterIPs []string, loadbalanceIP string, loadbalancePort string, subnet string, version string, tp string) {
	kubeadm := define.KubeadmTempST{}
	if stringsIn(masterIPs, loadbalanceIP) {
		kubeadm.APIServerCertSANs = masterIPs
	} else {
		kubeadm.APIServerCertSANs = append(masterIPs, loadbalanceIP)
	}
	//add other cert sans ips to APIServerCertSANs list
	for _, ip := range define.KubeFlags.OtherAPIServerCertSANs {
		if stringsIn(kubeadm.APIServerCertSANs, ip) {
		} else {
			kubeadm.APIServerCertSANs = append(kubeadm.APIServerCertSANs, ip)
		}
	}
	kubeadm.EtcdEndPoints = etcdIPs
	kubeadm.PodSubnet = subnet
	kubeadm.KubernetesVersion = version

	t := template.New("kubeadmConfig")
	pkg.Render(t, tp, kubeadm, "out/kube/kubeadm.yaml")
}

func genKubeletSystemdConfig(tp string) {
	driver := "cgroupfs"
	out, err := exec.Command("docker", "info").Output()
	outstr := string(out)
	if err != nil {
		fmt.Println("run docker info error: ", err)
	}
	if strings.Contains(outstr, "cgroupfs") {
	} else if strings.Contains(outstr, "systemd") {
		driver = "systemd"
	}

	t := template.New("systemdConfig")
	pkg.Render(t, tp, driver, "out/kube/10-kubeadm.conf")
}

//Gen is generate kubeadm kubelet config files
func (e *Kubecore) Gen() error {
	//gen kubeadm config file
	genKubeAdmConfigFile(define.KubeFlags.EtcdIPs, define.KubeFlags.MasterIPs, define.KubeFlags.LoadbalanceIP,
		define.KubeFlags.LoadbalancePort, define.KubeFlags.Subnet, define.KubeFlags.Version, define.KubeadmTemp)

	//gen kubelet config file
	genKubeletSystemdConfig(KubeletSystemdTemp)
	pkg.WriteFile("out/kube/kubelet.service", KubeletServiceStr)
	return nil
}

//Run is
func (e *Kubecore) Run() error {
	str := ApplyShellOutput(runSh)
	fmt.Println(str)
	//TODO fetch join command
	return nil
}

//Clean is
func (e *Kubecore) Clean() error {
	//TODO
	return nil
}

//InstallOffline is
func (e *Kubecore) InstallOffline() error {
	str := ApplyShellOutput(installOffline)
	fmt.Println(str)
	return nil
}

//InstallOnline is
func (e *Kubecore) InstallOnline() error {
	url := define.DownloadURL
	sh := pkg.RenderToStr(t, installOnlineWgetShTmpl, url)
	str := ApplyShellOutput(sh)
	fmt.Println(str)
	return nil
}

//Save is save bin files to bin dir, and save kubernetes core images
func (e *Kubecore) Save() error {
	str := ApplyShellOutput(sh)
	fmt.Println(str)
	return nil
}
