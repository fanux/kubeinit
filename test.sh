./kubeinit gen --master 192.168.11.16 --master 192.168.111.44  --master 192.168.111.50  \
    --etcd 192.168.11.16 --etcd 192.168.111.44 --etcd 192.168.111.50 --etcd-image gcr.io/google_containers/etcd-amd64:3.1.11\
    --version v1.9.2  \
    --loadbalance 10.1.245.93 
