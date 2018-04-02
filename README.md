## 2.0版本说明
此版本将整个项目重构，让其具备更好的灵活性与扩展性
![](https://github.com/fanux/kubeinit/blob/v2.0/docs/cluster.png)

## before starting
## ~~~generate config files~~~
```
$ kubeinit gen --master 1.1.1.2 --master 1.1.1.3 --master 1.1.1.4 \
               --etcd  1.1.1.2 --etcd 1.1.1.3 --etcd 1.1.1.4 \
               --loadbalance 1.1.1.2 --apply \
               --node 1.1.1.5 --node 1.1.1.6 --node 1.1.1.7 
```
--apply 生成配置立即执行 kubeinit apply, 不加这个参数只生成一些配置文件，方便定制需求去修改配置，修改完再apply

## install etcd and masters
```
$ kubeinit apply -bie
```
开关：
* -b 只初始化环境，拷贝bin文件，配置文件，导入镜像等
* -i 执行kubeadm init命令
* -e 启动etcd

## featrures
- [x] 自动生成多etcd节点compose文件
- [x] 自动生成kubeadm配置
- [x] 自动生成haproxy配置
- [x] 自动检测cgroup driver, 自动生成 kubelet配置文件

- [x] 自动初始化节点配置，拷贝bin文件，导入镜像
- [x] 自动初始化其它master节点配置
- [x] 自动启动etcd集群，master集群
- [x] 自动启动loadbalance
- [x] 自动join node节点, 修改node join参数为lb

- [x] 安装calico heapster dashboard

## create kubernetes HA cluster

> on master0:

```
$ kubeinit run --master 1.1.1.2 --master 1.1.1.3 --master 1.1.1.4 \
               --etcd  1.1.1.2 --etcd 1.1.1.3 --etcd 1.1.1.4 \
               --loadbalance 1.1.1.2 --loadbalance 1.1.1.2 --virturl-ip 1.1.2.2
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
