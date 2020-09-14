/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1pre1 "github.com/vmware-tanzu/service-apis/apis/v1alpha1pre1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeGatewayClasses implements GatewayClassInterface
type FakeGatewayClasses struct {
	Fake *FakeNetworkingV1alpha1pre1
}

var gatewayclassesResource = schema.GroupVersionResource{Group: "networking.x-k8s.io", Version: "v1alpha1pre1", Resource: "gatewayclasses"}

var gatewayclassesKind = schema.GroupVersionKind{Group: "networking.x-k8s.io", Version: "v1alpha1pre1", Kind: "GatewayClass"}

// Get takes name of the gatewayClass, and returns the corresponding gatewayClass object, and an error if there is any.
func (c *FakeGatewayClasses) Get(name string, options v1.GetOptions) (result *v1alpha1pre1.GatewayClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(gatewayclassesResource, name), &v1alpha1pre1.GatewayClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1pre1.GatewayClass), err
}

// List takes label and field selectors, and returns the list of GatewayClasses that match those selectors.
func (c *FakeGatewayClasses) List(opts v1.ListOptions) (result *v1alpha1pre1.GatewayClassList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(gatewayclassesResource, gatewayclassesKind, opts), &v1alpha1pre1.GatewayClassList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1pre1.GatewayClassList{ListMeta: obj.(*v1alpha1pre1.GatewayClassList).ListMeta}
	for _, item := range obj.(*v1alpha1pre1.GatewayClassList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gatewayClasses.
func (c *FakeGatewayClasses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(gatewayclassesResource, opts))
}

// Create takes the representation of a gatewayClass and creates it.  Returns the server's representation of the gatewayClass, and an error, if there is any.
func (c *FakeGatewayClasses) Create(gatewayClass *v1alpha1pre1.GatewayClass) (result *v1alpha1pre1.GatewayClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(gatewayclassesResource, gatewayClass), &v1alpha1pre1.GatewayClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1pre1.GatewayClass), err
}

// Update takes the representation of a gatewayClass and updates it. Returns the server's representation of the gatewayClass, and an error, if there is any.
func (c *FakeGatewayClasses) Update(gatewayClass *v1alpha1pre1.GatewayClass) (result *v1alpha1pre1.GatewayClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(gatewayclassesResource, gatewayClass), &v1alpha1pre1.GatewayClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1pre1.GatewayClass), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeGatewayClasses) UpdateStatus(gatewayClass *v1alpha1pre1.GatewayClass) (*v1alpha1pre1.GatewayClass, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(gatewayclassesResource, "status", gatewayClass), &v1alpha1pre1.GatewayClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1pre1.GatewayClass), err
}

// Delete takes name of the gatewayClass and deletes it. Returns an error if one occurs.
func (c *FakeGatewayClasses) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(gatewayclassesResource, name), &v1alpha1pre1.GatewayClass{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGatewayClasses) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(gatewayclassesResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1pre1.GatewayClassList{})
	return err
}

// Patch applies the patch and returns the patched gatewayClass.
func (c *FakeGatewayClasses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1pre1.GatewayClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(gatewayclassesResource, name, pt, data, subresources...), &v1alpha1pre1.GatewayClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1pre1.GatewayClass), err
}
