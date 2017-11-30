## before starting
TODO

## generate config files
```
$ kubeinit gen --master 1.1.1.2 --master 1.1.1.3 --master 1.1.1.4 \
               --etcd  1.1.1.2 --etcd 1.1.1.3 --etcd 1.1.1.4 \
               --loadbalance 1.1.1.2 --apply
```

## install etcd and master
```
$ kubeinit apply
```

## featrures
- [x] 自动生成多etcd节点compose文件
- [x] 自动生成kubeadm配置
- [x] 自动生成haproxy配置
- [x] 自动检测cgroup driver, 自动生成 kubelet配置文件

- [ ] 自动初始化节点配置，拷贝bin文件，导入镜像
- [ ] 自动初始化其它master节点配置
- [ ] 自动启动etcd集群，master集群
- [ ] 自动启动loadbalance
