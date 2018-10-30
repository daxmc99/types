package v1

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ComponentStatusLifecycle interface {
	Create(obj *v1.ComponentStatus) (*v1.ComponentStatus, error)
	Remove(obj *v1.ComponentStatus) (*v1.ComponentStatus, error)
	Updated(obj *v1.ComponentStatus) (*v1.ComponentStatus, error)
}

type componentStatusLifecycleAdapter struct {
	lifecycle ComponentStatusLifecycle
}

func (w *componentStatusLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1.ComponentStatus))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *componentStatusLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1.ComponentStatus))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *componentStatusLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1.ComponentStatus))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewComponentStatusLifecycleAdapter(name string, clusterScoped bool, client ComponentStatusInterface, l ComponentStatusLifecycle) ComponentStatusHandlerFunc {
	adapter := &componentStatusLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.ComponentStatus) (*v1.ComponentStatus, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(*v1.ComponentStatus); ok {
			return o, err
		}
		return nil, err
	}
}
