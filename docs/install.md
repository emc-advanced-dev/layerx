# Installing **Layer-X**

Prerequisites:
* [Go](https://golang.org/) 1.5 or later
* [etcd](https://github.com/coreos/etcd)
* at least one [Mesos Cluster](http://mesos.apache.org/gettingstarted/) - or - [Kubernetes Cluster](http://christopher5106.github.io/continous/deployment/2016/05/02/deploy-instantly-from-your-host-to-AWS-EC2-and-Google-Cloud-with-kubernetes.html)
* any Mesos framework (in this example, we'll be using [Marathon](https://mesosphere.github.io/marathon/docs/))

There are 3 components (binaries) that must be created and run in order to have a functioning **Layer-X** installation:

1. **Layer-X Core**: The management layer / stateful component. This is the glue that binds *task providers* to *resource providers*, provides the **Brain API** and the **UI**.

2. **Task-Provider Interface**: a passthrough layer that collects client requests for docker containers and translates them into the abstract **Layer-X** notion of a "task". E.g.: the *Mesos TPI* collects Mesos TaskInfos from Mesos Frameworks, converts them into *LXTasks*, and propagates them to the Core.

3. **Resource-Provider Interface**: a passthrough layer that collects avaiable resource nodes from cluster managers, translates them to the abstract **Layer-X** notion of a "node", and propagates them to the Core.

---

### Optional: Build with Make

You can build all 4 binaries with `make` in the root dir. You will need >= `go1.5`, and go-bindata && go go-bindata-assetfs installed. Make sure that `$GOPATH/bin` is part of your `$PATH`. Run

```
go get github.com/jteeuwen/go-bindata/... && go get github.com/elazarl/go-bindata-assetfs/...
cd <layerx-root-dir>
make
```

to build binaries. You'll typically want the Mesos RPI to run on the same host as the Mesos Master (I have found it works best this way). The core can run anywhere. The Mesos TPI should be somewhere it can communicate with both the Core and and Mesos Frameworks you are running. The K8s RPI can be run anywhere (alongside the core works fine).

---

### 1. Deploying The **Layer-X Core**

The first component we'll deploy is the Core. The core uses [etcd](https://github.com/coreos/etcd) as a persistent datastore. Follow instructions for setting up an etcd node or cluster [here](https://github.com/coreos/etcd#getting-etcd).

Once etcd is installed and running, you can launch the Core.

First, build the core:

```bash
##if running with go version < 1.6:
export GO15VENDOREXPERIMENT=1
##

cd layerx-core/
go build
```

Then, launch the core:
```bash
./layerx-core --etcd <etcd_endpoint>
```
this will run the core on the default port of 5000. The Core is now waiting for at least one TPI and one RPI to register to it.

### 2.a Deploying The **Mesos RPI**

Next, we'll want to bind a Mesos master to the Core and start collecting its nodes for **Layer-X** to manage. Ensure you have a [running Mesos cluster](http://mesos.apache.org/gettingstarted/) before proceeding with this step. For test/development purposes, we recommend using [vagrant-mesos](https://github.com/everpeace/vagrant-mesos) for a quick local environment setup.

Once your cluster is ready, build & run the RPI. It's recommended you run the RPI directly on the Mesos Master (or anywhere you've had success running a Mesos framework, as the RPI itself is a framework).

build:
```bash
##if running with go version < 1.6:
export GO15VENDOREXPERIMENT=1
##

cd layerx-mesos-rpi/
go build
```

run:
```bash
./layerx-mesos-rpi \
  --master <mesos_master_endpoint> \
  --layerx <ip-of-core-host>:5000 \
  --localip <ip-of-current-host>
```

this will run the RPI on port 4000. required flags:
* `--layerx` endpoint for layerx core. in the format ip:port
* `--master` endpoint for mesos master. in the format ip:port
* `--localip` broadcast address for the rpi. necessary for the Core to be able to communicate back to the RPI.

The RPI will register to Mesos, start collecting Resource Offers, and propagate them to the Core as Nodes.

### 2.b Deploying the **Kubernetes RPI**

Layer-X can run Mesos Frameworks on top of Kubernetes if you use the Kubernetes RPI. Note that Kubernetes only supports running docker containers, so any Mesos framework you run will need to be using the Docker Executor (see http://mesos.apache.org/documentation/latest/docker-containerizer/). Marathon is a good framework for this, and we'll demonstrate that in this tutorial.

build:
```bash
##if running with go version < 1.6:
export GO15VENDOREXPERIMENT=1
##

cd layerx-k8s-rpi/
go build
```

You will need to point the k8s rpi at a valid `kubeconfig` file that is configured for a running kubernetes cluser. See http://kubernetes.io/docs/user-guide/kubeconfig-file/ for info on kubeconfig file

run:
```bash
./layerx-k8s-rpi \
  -layerx <ip-of-core-host>:5000 \
  -localip <ip-of-current-host> \
  -kubeconfig <path/to/kube/config> \
  -port <1234> #optional
```

`-port` is optional, in case you're running this rpi on the same host as another rpi. by default they all want port 4000.

### 3. Deploying The **Mesos TPI**
Finally, we'll set the whole service in motion by binding a Mesos Framework (Marathon) to the Core through the Mesos TPI. Ensure you have Marathon & the JDK 1.8: https://mesosphere.github.io/marathon/docs/ (note: don't run it yet!)

build & run the TPI:

build:
```bash
##if running with go version < 1.6:
export GO15VENDOREXPERIMENT=1
##

cd layerx-mesos-tpi/
go build
```

run:
```bash
./layerx-mesos-tpi \
  --layerx <ip-of-core-host>:5000 \
  --localip <ip-of-current-host>

```

this will run the RPI on port 4000. required flags:
* `--layerx` endpoint for layerx core. in the format ip:port
* `--localip` broadcast address for the tpi. necessary for the Core to be able to communicate back to the TPI.

The **Layer-X** service is now bootstrapped and ready to start accepting tasks. Go to `http://Layer-X-IP:Port` in your browser to see the **UI**.

Finally, let's register Marathon as a task provider and start scheduling some tasks with **Layer-X**.

Run Marathon with the following command to attach it to the **Layer-X TPI** *instead of the Mesos Master*:

```bash
marathon-x.y.z/bin/start \
  --master <layerx-tpi-endpoint> \
  --task_launch_confirm_timeout 1200000 \
  --task_launch_timeout 1200000
```

note that the `--master` flag specifies the endpoint for the TPI, **not** the Mesos master!

the timeout flags are necessary in order to prevent Marathon from panicking when its tasks aren't immediately scheduled (since Layer-X will sit on them until they are scheduled by a Brain or through the UI).

Navigate to the Marathon UI and launch some apps, and navigate to the Layer-X UI in order to launch them!

Note: you can schedule tasks in the UI by dragging-and-dropping. You can also migrate tasks between nodes with drag-and-drop.

# Schedulers

In Layer-X, we call the process in charge of making scheduling / migration decisions the **Brain**. A Brain is just a client that consumes the REST API of the Layer-X core, that exposes information and control over the cluster (all of tasks and resources being managed by Layer-X). Technically, the UI is one such Brain. An example brain is provided in [../example-brain/](../example-brain). More complex brains can be written that take full advantage of all of Layer-X's APIs, as well as the information that Layer-X stores about tasks and resources.
