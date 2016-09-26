CORE=layerx-core
K8S_RPI=layerx-k8s-rpi
MESOS_RPI=layerx-mesos-rpi
MESOS_TPI=layerx-mesos-tpi

.PHONY: all
all: bin/$(CORE) bin/$(K8S_RPI) bin/$(MESOS_RPI) bin/$(MESOS_TPI)

bin/$(CORE): $(shell find $(CORE) -name '*.go')
	pushd $(CORE) && go build -v -o ../bin/$(CORE) && popd

bin/$(K8S_RPI): $(shell find $(K8S_RPI) -name '*.go')
	pushd $(K8S_RPI) && go build -v -o ../bin/$(K8S_RPI) && popd

bin/$(MESOS_RPI): $(shell find $(MESOS_RPI) -name '*.go')
	pushd $(MESOS_RPI) && go build -v -o ../bin/$(MESOS_RPI) && popd

bin/$(MESOS_TPI): $(shell find $(MESOS_TPI) -name '*.go')
	pushd $(MESOS_TPI) && go build -v -o ../bin/$(MESOS_TPI) && popd

.PHONY: clean

clean:
	rm -rf bin/

