// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/fanux/kubeinit/define"
	"github.com/spf13/cobra"
)

//WriteFile is
func WriteFile(fileName string, content string) {
	b := []byte(content)
	err := ioutil.WriteFile(fileName, b, 0644)
	if err != nil {
		fmt.Println("write file error", err)
	}
}

//FileExists is
func FileExists(file string) bool {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		return true
	}
	return false

}

//GenEtcdYaml iko
func GenEtcdYaml(endpoint, endpoints, template string) (out string) {
	//TODO
	return
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

//Render is
func Render(t *template.Template, tp string, args interface{}, outFile string) {
	template.Must(t.Parse(tp))

	file, err := os.Create(outFile)
	defer file.Close()
	if err != nil {
		fmt.Println("create out file error: %s", err)
		return
	}

	err = t.Execute(file, args)
	if err != nil {
		fmt.Println("exec template file error: %s", err)
	}
}

func genEtcdyamls(etcdIPs []string, tp string) {
	var etcdComposeFileNmae string

	etcd := define.EtcdComposeTempST{}
	etcd.EndPoints = getEtcdEndpoints(etcdIPs)
	etcd.Image = define.EtcdImage

	for i, ip := range etcdIPs {
		etcd.EndPoint = ip
		etcd.Index = strconv.Itoa(i)

		etcdComposeFileNmae = fmt.Sprintf("etcd-docker-compose-%d.yml", i)
		t := template.New("etcd")

		Render(t, tp, etcd, etcdComposeFileNmae)
	}
}

func genKubeAdmConfigFile(etcdIPs []string, masterIPs []string, loadbalanceIP string, loadbalancePort string, template string) {
}

func genLoadbalanceConfigFile(loadbalanceIP string, loadbalancePort string, masterIPs []string, haproxyTemp string) {
}

func genKubeletSystemdConfig(kubeletSystemdTemp string) {
}

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate config files, include etcd docker compose file and kubeadm config file",
	Long:  `you can generate it then apply it, if using apply will generate configs if not exist`,
	Run: func(cmd *cobra.Command, args []string) {
		genEtcdyamls(define.EtcdIPs, define.EtcdComposeTemp)
		//genKubeAdmConfigFile(define.EtcdIPs, define.MasterIPs, define.LoadbalanceIP, define.LoadbalancePort, define.KubeadmTemp)
		//genLoadbalanceConfigFile(define.LoadbalanceIP, define.LoadbalancePort, define.MasterIPs, define.HaproxyTemp)
		//genKubeletSystemdConfig(define.KubeletSystemdTemp)
	},
}

func init() {
	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	genCmd.Flags().StringSliceVar(&define.EtcdIPs, "etcd", []string{"127.0.0.1"}, "etcd ips")
	genCmd.Flags().StringSliceVar(&define.MasterIPs, "master", []string{"127.0.0.1"}, "master ips")
	genCmd.Flags().StringVar(&define.LoadbalanceIP, "loadbalance", "127.0.0.1", "loadbalance ip")
	genCmd.Flags().StringVar(&define.LoadbalancePort, "loadbalance-port", ":6444", "loadbalance port")
	genCmd.Flags().StringVar(&define.EtcdImage, "etcd-image", "gcr.io/google_containers/etcd-amd64:3.0.17", "etcd docker image")
	genCmd.Flags().BoolVarP(&define.Apply, "apply", "a", false, "apply directly")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
