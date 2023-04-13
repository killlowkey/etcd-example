# Etcd 学习
1. CRUD 操作
2. Watch 机制
3. 事务操作
4. Etcd 集群

## Docker
docker 运行 etcd
```shell
docker run -d --name Etcd-server \
    --network app-tier \
    --publish 2379:2379 \
    --publish 2380:2380 \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    --env ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 \
    bitnami/etcd:latest
```

docker-compose 运行 etcd
```yml
version: '3'

services:
  Etcd:
    image: 'bitnami/etcd:latest'
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379
    ports:
      - "2379:2379"
      - "2380:2380"
```

## 文档
1. https://etcd.io/
2. https://pkg.go.dev/go.etcd.io/etcd/client/v3