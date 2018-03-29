package define

//EtcdComposeTempST is
type EtcdComposeTempST struct {
	Index     string
	Image     string
	EndPoint  string
	EndPoints string
}

//HaproxyTemp  is
var HaproxyTemp = `global
  daemon
  log 127.0.0.1 local0
  log 127.0.0.1 local1 notice
  maxconn 4096

defaults
  log               global
  retries           3
  maxconn           2000
  timeout connect   5s
  timeout client    50s
  timeout server    50s

frontend k8s
  bind *:{{.LoadbalancePort}}
  mode tcp
  default_backend k8s-backend

backend k8s-backend
  balance roundrobin
  mode tcp
{{range $index, $element := .BackendEndPoint}}  server k8s-{{$index}} {{$element}}:6443 check {{"\n"}}{{end}}
`

//HaproxyTempST is
type HaproxyTempST struct {
	LoadbalancePort string
	BackendEndPoint []string
}

//Flags is
type Flags struct {
	MasterIPs              []string
	NodeIPs                []string
	OtherAPIServerCertSANs []string
	EtcdIPs                []string
	MasterEndPoints        []string
	APIServerCertSANs      []string
	PodSubnet              string
	KubernetesVersion      string
	LoadbalancePort        string
	LoadbalanceIP          string
	Apply                  bool
	EtcdImage              string
	Version                string
	Subnet                 string

	// render out compose file and kubeadm config
	ConfigOutDir string
	// /etc/kubernetes files, this need copy to other nodes, and change ips
	KubernetesDir string
}

//KubeFlags is
var (
	KubeFlags Flags

	InitBaseEnvironment bool
	InitKubeadm         bool
	InitOtherMasters    bool
	StartEtcdCluster    bool
	Distribute          bool
	Pssh                bool
	User                string
	Password            string
)
