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
	"encoding/json"
	"fmt"
	"html/template"
	"os"

	"github.com/fanux/kubeinit/define"
	"github.com/spf13/cobra"
)

var outDir = "out"

//FileExists is
func FileExists(file string) bool {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		return true
	}
	return false
}

func stringsIn(s []string, key string) bool {
	for _, i := range s {
		if i == key {
			return true
		}
	}
	return false
}

func genLoadbalanceConfigFile(loadbalancePort string, masterIPs []string, tp string) {
	haproxy := define.HaproxyTempST{
		loadbalancePort,
		masterIPs}

	t := template.New("haproxy")
	outfile := fmt.Sprintf("%s/haproxy.cfg", outDir)
	Render(t, tp, haproxy, outfile)
}

func dumpKubeInitConfig(flags define.Flags) {
	fileName := fmt.Sprintf("%s/kubeinit.json", outDir)
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return // exit the function on error
	}
	defer output.Close()
	json.NewEncoder(output).Encode(&flags)
	//json.NewDecoder(os.Stdin)Decode(&flags)
}

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate config files, include etcd docker compose file and kubeadm config file",
	Long:  `you can generate it then apply it, if using apply will generate configs if not exist`,
	Run: func(cmd *cobra.Command, args []string) {
		/*
			genEtcdyamls(define.KubeFlags.EtcdIPs, define.EtcdComposeTemp)
			genKubeAdmConfigFile(define.KubeFlags.EtcdIPs, define.KubeFlags.MasterIPs, define.KubeFlags.LoadbalanceIP,
				define.KubeFlags.LoadbalancePort, define.KubeFlags.Subnet, define.KubeFlags.Version, define.KubeadmTemp)
			genLoadbalanceConfigFile(define.KubeFlags.LoadbalancePort, define.KubeFlags.MasterIPs, define.HaproxyTemp)
			genKubeletSystemdConfig(define.KubeletSystemdTemp)

			//kubeinit appy needs this arguments
			dumpKubeInitConfig(define.KubeFlags)
		*/
	},
}

func init() {
	err := os.Mkdir(outDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	genCmd.Flags().StringSliceVar(&define.KubeFlags.EtcdIPs, "etcd", []string{"127.0.0.1"}, "etcd ips")
	genCmd.Flags().StringSliceVar(&define.KubeFlags.MasterIPs, "master", []string{"127.0.0.1"}, "master ips")
	genCmd.Flags().StringSliceVar(&define.KubeFlags.OtherAPIServerCertSANs, "cert-sans", []string{}, "other api server cert sans, like floating ips")
	genCmd.Flags().StringVar(&define.KubeFlags.LoadbalanceIP, "loadbalance", "127.0.0.1", "loadbalance ip")
	genCmd.Flags().StringVar(&define.KubeFlags.LoadbalancePort, "loadbalance-port", "6444", "loadbalance port")
	genCmd.Flags().StringVar(&define.KubeFlags.EtcdImage, "etcd-image", "gcr.io/google_containers/etcd-amd64:3.0.17", "etcd docker image")
	genCmd.Flags().BoolVarP(&define.KubeFlags.Apply, "apply", "a", false, "apply directly")
	genCmd.Flags().StringVar(&define.KubeFlags.Subnet, "pod-subnet", "10.122.0.0/16", "pod subnet")
	genCmd.Flags().StringVar(&define.KubeFlags.Version, "version", "v1.9.1", "kubernetes version")
	genCmd.Flags().StringSliceVar(&define.KubeFlags.NodeIPs, "node", []string{}, "node ips")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
