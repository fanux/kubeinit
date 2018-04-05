package etcd

//EtcdComposeTemp is
var EtcdComposeTemp = `version: '2.1'
services:
    etcd{{.Index}}:
        container_name: etcd_infra{{.Index}}
        image: {{.Image}}
        command: |
                etcd --name infra{{.Index}}
                --initial-advertise-peer-urls http://{{.EndPoint}}:2380
                --listen-peer-urls http://{{.EndPoint}}:2380
                --listen-client-urls http://{{.EndPoint}}:2379,http://127.0.0.1:2379
                --advertise-client-urls http://{{.EndPoint}}:2379
                --data-dir /etcd-data.etcd
                --initial-cluster-token etcd-cluster-1
                --initial-cluster {{.EndPoints}}
                --initial-cluster-state new
        restart: always
        volumes:
           - /data/etcd-data.etcd:/etcd-data.etcd
        network_mode: "host"

`

//ComposeTempST is
type ComposeTempST struct {
	Index     string
	Image     string
	EndPoint  string
	EndPoints string
}
