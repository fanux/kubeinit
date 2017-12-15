./kubeinit gen --master 172.16.22.152   \
        --etcd 172.16.22.152  \
            --version 1.8.5 \
                --loadbalance  172.16.22.152

./kubeinit apply -bie --passwd Fanux#123
