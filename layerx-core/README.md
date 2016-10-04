# Core

Core just requires etcd to be running. Hook up an RPI and TPI to it to get it running. You can view the UI of the core at <core-ip:port> in browser.

If you get Layer X in a bad state, you can kill all the layerx processes and clear etcd with `curl -XDELETE <etcd_endpoint>/v2/keys/state?recursive=true`

Note this will orphan anything Layer-X has created and didn't clean up (like kubernetes pods).

### set up etcd nice and ezpz

install etcd:
```
curl -L  https://github.com/coreos/etcd/releases/download/v2.3.0-alpha.1/etcd-v2.3.0-alpha.1-darwin-amd64.zip -o etcd-v2.3.0-alpha.1-darwin-amd64.zip
unzip etcd-v2.3.0-alpha.1-darwin-amd64.zip
cd etcd-v2.3.0-alpha.1-darwin-amd64
mv etcd /usr/local/bin/
```
