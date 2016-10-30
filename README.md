# <img src="http://i.imgur.com/idwFRSK.png" alt="Container Scheduling across Clusters" width="159" height="50">

**Layer X** a single centralized service that manages containers across multiple clusters, multiple clouds, and multiple cluster managers. it joins together pools cluster managers such as [*Kubernetes*](http://kubernetes.io/), [*Docker Swarm*](https://docs.docker.com/swarm/), and [*Mesos*](http://mesos.apache.org/) into a single shared group of resources, where any container can be managed, scaled, migrated across a unified resource pool. 

**Layer X** exposes a set of cluster-management APIs including scheduling & resource management to so-called “brains” which enables advanced features such as live migration & continuous placement of containers, easy flex up & flex down of cluster, and dynamic resizing of resource constraints.


![Layer-X Arch](http://i.imgur.com/O5TxJLF.png "Architecture")


In the **Layer-X** world:
* **Resource Providers** offer nodes with useable resources (e.g. a [Mesos Slave](https://open.mesosphere.com/reference/mesos-slave/), a [Kubelet](http://kubernetes.io/docs/admin/kubelet/), or a [Swarm Node](https://docs.docker.com/engine/swarm/swarm-tutorial/add-nodes/)) to be consumed by
* **Task Providers**, which create services and perform work by launching Docker containers (e.g. [Marathon](https://mesosphere.github.io/marathon/), [Deis](http://deis.io/))
* **Brains** interface with the **Layer-X** HTTP API to manage the global state of the cluster, including scheduling *tasks* on *nodes*.

**Layer-X** abstracts specific implementations of **task providers** and **resource-providers** with **Task Provider Interfaces** (TPIs) and **Resource Provider Interfaces** (RPIs). Currently, TPI and RPI have only been completed for Mesos, but Kubernetes and Docker Swarm will be available soon™.

**Layer-X** implements a simple interface for adding TPIs and RPIs to encourage easy community contributions for new resource providers (Nomad, Hadoop, etc.) and task providers (Dokku). Client APIs for TPIs are visible [here](./layerx-core/layerx_tpi_client/layerx_tpi.go) and for RPIs [here](./layerx-core/layerx_rpi_client/layerx_rpi.go).

**Brains** are also designed to be pluggable, 3rd party extensions for **Layer-X**. The Brain API is visible through the Brain Client package [here](./layerx-core/layerx_brain_client/layerx_brain_client.go)

Get started with [installation documentation](docs/install.md).

---

A rough diagram of Layer-X architecture:

![Layer-X Arch Diagram](http://i.imgur.com/GcYh5ug.png "Architecture")
