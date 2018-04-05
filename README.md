```
 _          _          _       _ _   
| |        | |        (_)     (_) |  
| | ___   _| |__   ___ _ _ __  _| |_ 
| |/ / | | | '_ \ / _ \ | '_ \| | __|
|   <| |_| | |_) |  __/ | | | | | |_ 
|_|\_\\__,_|_.__/ \___|_|_| |_|_|\__|

_|                  _|                  _|            _|    _|      
_|  _|    _|    _|  _|_|_|      _|_|        _|_|_|        _|_|_|_|  
_|_|      _|    _|  _|    _|  _|_|_|_|  _|  _|    _|  _|    _|      
_|  _|    _|    _|  _|    _|  _|        _|  _|    _|  _|    _|      
_|    _|    _|_|_|  _|_|_|      _|_|_|  _|  _|    _|  _|      _|_| 

    __         __         _       _ __ 
   / /____  __/ /_  ___  (_)___  (_) /_
  / //_/ / / / __ \/ _ \/ / __ \/ / __/
 / ,< / /_/ / /_/ /  __/ / / / / / /_  
/_/|_|\__,_/_.___/\___/_/_/ /_/_/\__/ 
```

## 2.0版本说明
此版本将整个项目重构，让其具备更好的灵活性与扩展性
![](https://github.com/fanux/kubeinit/blob/v2.0/docs/cluster.png)

## featrures
- [x] 自动生成多etcd节点compose文件
- [x] 自动生成kubeadm配置
- [x] ~~~自动生成haproxy配置~~~ LVS代替
- [x] 自动检测cgroup driver, 自动生成 kubelet配置文件

- [x] 自动初始化节点配置，拷贝bin文件，导入镜像
- [x] 自动初始化其它master节点配置
- [x] 自动启动etcd集群，master集群
- [x] 自动启动loadbalance
- [x] 自动join node节点, 修改node join参数为lb

- [x] 安装calico heapster dashboard
- [ ] 自定义安装addons

## create kubernetes HA cluster
> on master0:

```
$ kubeinit run --master 1.1.1.2 --master 1.1.1.3 --master 1.1.1.4 \
               --etcd  1.1.1.2 --etcd 1.1.1.3 --etcd 1.1.1.4 \
               --loadbalance 1.1.1.2 --loadbalance 1.1.1.2 --virturl-ip 1.1.2.2 \
               --with-dashboard --with-heapster --with-promethus --with-EFK
```
output:
```
kubeinit serve on: 1.1.1.2:9527
join master command: kubeinit join master 1.1.1.2:9527 sha:xxxxxx
join node command: kubeinit join node 1.1.1.2:9527 sha:xxxxxx
```

> on other masters:
```
 kubeinit join master 1.1.1.2:9527 sha:xxxxxx
```

> on other nodes:
```
 kubeinit join node 1.1.1.2:9527 sha:xxxxxx
```

all this down, stop master0 kubeinit serve
