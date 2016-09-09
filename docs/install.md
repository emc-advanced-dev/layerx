# Installing **Layer-X**

Prerequisites:
* [Go](https://golang.org/) 1.5 or later
* [etcd](https://github.com/coreos/etcd)
* at least one [Mesos Cluster](http://mesos.apache.org/gettingstarted/)
* any Mesos framework (in this example, we'll be using [Marathon](https://mesosphere.github.io/marathon/docs/))

There are 3 components (binaries) that must be created and run in order to have a functioning **Layer-X** installation:

1. **Layer-X Core**: The management layer / stateful component. This is the glue that binds *task providers* to *resource providers*, provides the **Brain API** and the **UI**.

2. **Task-Provider Interface**: a passthrough layer that collects client requests for docker containers and translates them into the abstract **Layer-X** notion of a "task". E.g.: the *Mesos TPI* collects Mesos TaskInfos from Mesos Frameworks, converts them into *LXTasks*, and propagates them to the Core.

3. **Resource-Provider Interface**: a passthrough layer that collects avaiable resource nodes from cluster managers, translates them to the abstract **Layer-X** notion of a "node", and propagates them to the Core.

Note: as of the release of Layer-X, the only finished TPI and RPI are a TPI for Mesos frameworks, and an RPI for Mesos slaves. We are currently developing our Kubernetes TPI and RPI.

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
this will run the core on the default port of 5555. The Core is now waiting for at least one TPI and RPI to register to it.

### 2. Deploying The **Mesos RPI**

Next, we'll want to bind a Mesos master to the Core and start collecting its nodes for **Layer-X** to manage. Ensure you have a [running Mesos cluster](http://mesos.apache.org/gettingstarted/) before proceeding with this step. For test/development purposes, we recommend using [vagrant-mesos](https://github.com/everpeace/vagrant-mesos) for a quick local environment setup.

Once your cluster is ready, build & run the RPI.

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
  --layerx 127.0.0.1:5050 \
  --localip 192.168.0.20

```

this will run the RPI on port 4000. required flags:
* `--layerx` endpoint for layerx core. in the format ip:port
* `--master` endpoint for mesos master. in the format ip:port
* `--localip` broadcast address for the rpi. necessary for the Core to be able to communicate back to the RPI.

The RPI will register to Mesos, start collecting Resource Offers, and propagate them to the Core as Nodes.

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
  --layerx 127.0.0.1:5050 \
  --localip 192.168.0.20

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
