CORE=layerx-core
K8S_RPI=layerx-k8s-rpi
MESOS_RPI=layerx-mesos-rpi
SWARM_RPI=layerx-swarm-rpi
MESOS_TPI=layerx-mesos-tpi

.PHONY: all
all: bin/$(CORE) bin/$(K8S_RPI) bin/$(MESOS_RPI) bin/$(MESOS_TPI)

bin/$(CORE): $(shell find $(CORE) -name '*.go') $(CORE)/bindata/bindata_assetfs.go
	cd $(CORE) && go build -v -o ../bin/$(CORE)

$(CORE)/bindata/bindata_assetfs.go: $(shell find $(CORE)/web)
	cd $(CORE) && mkdir -p bindata && go-bindata-assetfs -pkg bindata web/... && mv bindata_assetfs.go bindata

bin/$(K8S_RPI): $(shell find $(K8S_RPI) -name '*.go')
	cd $(K8S_RPI) && go build -v -o ../bin/$(K8S_RPI)

bin/$(MESOS_RPI): $(shell find $(MESOS_RPI) -name '*.go')
	cd $(MESOS_RPI) && go build -v -o ../bin/$(MESOS_RPI)

bin/$(SWARM_RPI): $(shell find $(SWARM_RPI) -name '*.go')
	cd $(SWARM_RPI) && go build -v -o ../bin/$(SWARM_RPI)

bin/$(MESOS_TPI): $(shell find $(MESOS_TPI) -name '*.go')
	cd $(MESOS_TPI) && go build -v -o ../bin/$(MESOS_TPI)

.PHONY: clean

clean:
	rm -rf bin/

