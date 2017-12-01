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

	"github.com/fanux/kubeinit/define"
	"github.com/fanux/kubeinit/pkg"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "init env and appy kubeadmin init",
	Long: `初始化环境，拷贝bin程序，配置文件，加载镜像，执行kubeadm, 初始化其它节点
	init env, copy binarys config files, load docker images, exec kubeadm init and setup other master nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("appy called")
		pkg.Apply()
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	applyCmd.Flags().BoolVarP(&define.InitBaseEnvironment, "init-base-env", "b", false, "init base environment, close firewalld selinux swap, copy bin and configs, load docker images")
	applyCmd.Flags().BoolVarP(&define.InitKubeadm, "init-kubeadm", "i", false, "exec kubeadm init")
	applyCmd.Flags().BoolVarP(&define.StartEtcdCluster, "start-etcd", "e", false, "docker compose up etcd compose files")

}
