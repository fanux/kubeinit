## before starting
TODO

## generate config files
```
$ kubeinit gen --master 1.1.1.2 --master 1.1.1.3 --master 1.1.1.4 \
               --etcd  1.1.1.2 --etcd 1.1.1.3 --etcd 1.1.1.4 \
               --loadbalance 1.1.1.2
```

## install etcd and master
```
$ kubeinit apply
```
