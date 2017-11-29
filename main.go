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

package main

import (
	"fmt"

	"github.com/fanux/kubeinit/cmd"
)

var showCluster = ` cluster overview >
          kubectl dashboard
                 |
                 V 
     +------------------------+ join
     | LB  10.1.245.94 haproxy| <--- Nodes
     +------------------------+
     |                                                   
     |--master1 manager1 schedule1   10.1.245.93                                                
     |--master2 manager2 schedule2   10.1.245.95    =============>  etcd cluster  http://10.1.245.93:2379,http://10.1.245.94:2379,http://10.1.245.95:2379
     +--master3 manager3 schedule3   10.1.245.94   


   +---------------------------------+
   | qq群：98488045                  | 
   | phone NO. 15357921248           | 
   | Email: lamelegdog@gmail.com     |
   +---------------------------------+
`

func main() {
	fmt.Println(showCluster)
	cmd.Execute()
}
