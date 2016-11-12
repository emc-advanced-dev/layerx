# Layer-X All-In-One

This directory uses a set of scripts [vagrant](https://www.vagrantup.com/), [minkube](https://github.com/kubernetes/minikube), and [docker-machine](https://docs.docker.com/machine/) to deploy each of the following on virtualbox:
- One Layer X VM with:
  - ETCD
  - Layer-X RPIs and TPIs
  - mesos master
  - Marathon
- One Zookeeper vm (required for mesos)
- One Mesos Slave 
- (Optionally):
  - One K8s vm
  - One docker-swarm node

# Prerequisites:

- [Virtualbox](https://www.virtualbox.org/wiki/Downloads) or an AWS account
- [Vagrant](https://www.vagrantup.com/docs/installation/)
- The following Vagrant plugins:
   - `vagrant-berkshelf`
   - `vagrant-omnibus`
   - `vagrant-hosts`
   - These can be installed with the `vagrant plugin install <plugin-name>` command

  If you want a Kubernetes node in the cluster:
- [minikube](http://kubernetes.io/docs/getting-started-guides/minikube/) and [kubectl](http://kubernetes.io/docs/user-guide/prereqs/)

  If you want a Docker Swarm node in the cluster:
- [docker-machine](https://docs.docker.com/machine/install-machine/) and [docker cli >= 1.12](https://github.com/docker/docker/releases)


## Kubernetes Prerequisities

Install Kubectl:
```bash
# linux/amd64
curl -Lo kubectl http://storage.googleapis.com/kubernetes-release/release/v1.3.0/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
# OS X/amd64 
curl -Lo kubectl http://storage.googleapis.com/kubernetes-release/release/v1.3.0/bin/darwin/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
```
Install Minikube:
```bash
# linux/amd64
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.12.2/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
# OS X/amd64 
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.12.2/minikube-darwin-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
```

## Docker Swarm Prerequisites

Install Docker & Docker-Machine
```bash
#docker
curl -sSL https://get.docker.com/ | sh
#docker-machine
curl -L https://github.com/docker/machine/releases/download/v0.8.2/docker-machine-`uname -s`-`uname -m` >/usr/local/bin/docker-machine && \
chmod +x /usr/local/bin/docker-machine
```

# Deploy the Cluster

## Configure

If you want to change any configuration options, see `cluster.yml`
To turn off deploying Kubernetes or Docker Swarm, set `deploy_kubernetes` or `deploy_swarm` to `false`, respectively. These are on by default.

## Deploy!

`cd` to this directory and simply run `vagrant up`. 

Deploying may take 15-30 minutes to deploy, and consume several GB of memory for the VMs, as well as a few tens of GB of disk space. When deploy finishes, direct your browser to http://172.31.1.11:5000/ to see the Layer-x UI. 

To kill the cluster, run `vargant destroy`