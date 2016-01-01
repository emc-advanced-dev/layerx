package fakes
import "github.com/layer-x/layerx-core_v2/lxtypes"

func FakeTaskProvider(tpid, source string) *lxtypes.TaskProvider {
	return &lxtypes.TaskProvider{
		Id: tpid,
		Source: source,
	}
}