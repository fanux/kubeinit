./kubeinit gen --master 172.16.22.156  --master  172.16.22.155 --master  172.16.22.157  \
               --etcd  172.16.22.155  --etcd  172.16.22.156 --etcd  172.16.22.157 \
               --node  172.16.22.153  --node  172.16.22.154  \
               --version v1.8.5 \
               --loadbalance  172.16.22.155

./kubeinit apply -bie --passwd Fanux#123
