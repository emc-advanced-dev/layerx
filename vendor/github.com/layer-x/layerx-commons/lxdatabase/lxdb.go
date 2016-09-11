package lxdatabase

import (
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/layer-x/layerx-commons/lxerrors"
	"golang.org/x/net/context"
	"strings"
	"sync"
	"time"
)

var c client.Client
var m *sync.Mutex

func Init(etcdEndpoints []string) error {
	cfg := client.Config{
		Endpoints:               etcdEndpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	var err error
	c, err = client.New(cfg)
	if err != nil {
		return lxerrors.New("initialize etcd", err)
	}
	m = &sync.Mutex{}
	return nil
}

func Get(key string) (string, error) {
	m.Lock()
	defer m.Unlock()
	key = prefixKey(key)
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Get(context.Background(), key, nil)
	if err != nil {
		return "", lxerrors.New("getting key/val", err)
	}
	if resp.Node.Dir {
		return "", lxerrors.New("get used on a dir", err)
	}
	return resp.Node.Value, nil
}

func Set(key string, value string) error {
	m.Lock()
	defer m.Unlock()
	key = prefixKey(key)
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Set(context.Background(), key, value, nil)
	if err != nil {
		return lxerrors.New("setting key/val pair", err)
	}
	if resp.Node.Key != key || resp.Node.Value != value {
		fmt.Printf("key was %s, value was %s", resp.Node.Key, resp.Node.Value)
		return lxerrors.New("key/value pair not set as expected", nil)
	}
	return nil
}

func Rm(key string) error {
	m.Lock()
	defer m.Unlock()
	key = prefixKey(key)
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Delete(context.Background(), key, nil)
	if err != nil {
		return lxerrors.New("deleting key/val pair", err)
	}
	if resp.Node.Key != key {
		return lxerrors.New("removed pair does not have expected key", nil)
	}
	return nil
}

func Mkdir(dir string) error {
	m.Lock()
	defer m.Unlock()
	dir = prefixKey(dir)
	kapi := client.NewKeysAPI(c)
	opts := &client.SetOptions{
		Dir: true,
	}
	resp, err := kapi.Set(context.Background(), dir, "", opts)
	if err != nil {
		return lxerrors.New("making directory", err)
	}
	if resp.Node.Key != dir || !resp.Node.Dir {
		return lxerrors.New("directory not created as expected", nil)
	}
	return nil
}

func Rmdir(dir string, recursive bool) error {
	m.Lock()
	defer m.Unlock()
	dir = prefixKey(dir)
	kapi := client.NewKeysAPI(c)
	opts := &client.DeleteOptions{
		Dir:       true,
		Recursive: recursive,
	}
	resp, err := kapi.Delete(context.Background(), dir, opts)
	if err != nil {
		return lxerrors.New("removing directory", err)
	}
	if resp.Node.Key != dir || !resp.Node.Dir {
		return lxerrors.New("directory not created as expected", nil)
	}
	return nil
}

func GetKeys(dir string) (map[string]string, error) {
	m.Lock()
	defer m.Unlock()
	dir = prefixKey(dir)
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Get(context.Background(), dir, nil)
	if err != nil {
		return map[string]string{}, lxerrors.New("getting key/vals for dir", err)
	}
	if !resp.Node.Dir {
		return map[string]string{}, lxerrors.New("ls used on a non-dir key", err)
	}
	result := make(map[string]string)
	for _, node := range resp.Node.Nodes {
		if !node.Dir {
			result[node.Key] = node.Value
		} //ignore directories
	}
	return result, nil
}

func GetSubdirectories(dir string) ([]string, error) {
	m.Lock()
	defer m.Unlock()
	dir = prefixKey(dir)
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Get(context.Background(), dir, nil)
	if err != nil {
		return []string{}, lxerrors.New("getting key/vals for dir", err)
	}
	if !resp.Node.Dir {
		return []string{}, lxerrors.New("ls used on a non-dir key", err)
	}
	result := []string{}
	for _, node := range resp.Node.Nodes {
		if node.Dir {
			result = append(result, node.Key)
		} //ignore keys
	}
	return result, nil
}

func prefixKey(key string) string {
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}
	return key
}
