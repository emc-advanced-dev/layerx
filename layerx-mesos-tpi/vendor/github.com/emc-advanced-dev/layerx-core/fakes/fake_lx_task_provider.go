package fakes
import "github.com/emc-advanced-dev/layerx/layerx-core/lxtypes"

func FakeTaskProvider(tpid, source string) *lxtypes.TaskProvider {
	return &lxtypes.TaskProvider{
		Id: tpid,
		Source: source,
	}
}
