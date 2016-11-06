# Layer-X All-In-One

This directory uses a set of scripts ([vagrant]() & bash) to deploy each of the following on virtualbox or AWS:
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

Prerequisites:

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
