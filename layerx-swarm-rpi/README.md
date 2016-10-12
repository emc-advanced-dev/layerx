swarm cluster

docker-machine create -d amazonec2 \
swarm-master

docker-machine create -d amazonec2 \
swarm-worker-01

docker-machine create -d amazonec2 \
swarm-worker-02

eval $(docker-machine env swarm-master)

# Layer-X RPI for Docker Swarm (the 1.12 version, not standalone)
Prerequisites:
- Running Swarm cluster
- Env vars set up on local terminal where you're running the RPI (easiest way to do this is with `eval $(docker-machine env swarm-manager)` assuming your swarm master node is managed by docker-machine).
- Docker 1.12 (or greater) installed on local machine.

build: use make or 
```bash
go build
```

run:
```bash
./layerx-mesos-rpi \
  --layerx <ip-of-core-host>:5000 \
  --localip <ip-of-current-host>
```