package etcd

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/fanux/kubeinit/define"
	"github.com/fanux/kubeinit/pkg"
)

//Etcd is
type Etcd struct {
	tp ComposeTempST
}

//Info is
func (e *Etcd) Info() (string, string) {
	return "etcd", nil
}

//Gen is
func (e *Etcd) Gen() error {
	genEtcdyamls(define.KubeFlags.EtcdIPs, EtcdComposeTemp)
	return nil
}

//Run is
func (e *Etcd) Run() error {
	//TODO
	return nil
}

//Clean is
func (e *Etcd) Clean() error {
	//TODO
	return nil
}

// infra0=http://10.1.245.93:2380,infra1=http://10.1.245.94:2380,infra2=http://10.1.245.95:2380
func getEtcdEndpoints(etcdIPs []string) (out string) {
	for i, ip := range etcdIPs {
		var temp string
		temp = fmt.Sprintf("infra%d=http://%s:2380,", i, ip)
		out = out + temp
	}

	out = out[:len(out)-1]

	fmt.Println("etcd endpoints: ", out)
	return
}

func genEtcdyamls(etcdIPs []string, tp string) {
	var etcdComposeFileNmae string

	etcd := define.ComposeTempST{}
	etcd.EndPoints = getEtcdEndpoints(etcdIPs)
	etcd.Image = define.KubeFlags.EtcdImage

	for i, ip := range etcdIPs {
		etcd.EndPoint = ip
		etcd.Index = strconv.Itoa(i)

		etcdComposeFileNmae = fmt.Sprintf("%s/etcd-docker-compose-%d.yml", outDir, ip)
		t := template.New("etcd")

		pkg.Render(t, tp, etcd, etcdComposeFileNmae)
	}
}
