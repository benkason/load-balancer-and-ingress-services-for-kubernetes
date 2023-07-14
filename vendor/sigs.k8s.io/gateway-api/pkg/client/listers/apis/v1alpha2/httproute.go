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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha2

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// HTTPRouteLister helps list HTTPRoutes.
// All objects returned here must be treated as read-only.
type HTTPRouteLister interface {
	// List lists all HTTPRoutes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha2.HTTPRoute, err error)
	// HTTPRoutes returns an object that can list and get HTTPRoutes.
	HTTPRoutes(namespace string) HTTPRouteNamespaceLister
	HTTPRouteListerExpansion
}

// hTTPRouteLister implements the HTTPRouteLister interface.
type hTTPRouteLister struct {
	indexer cache.Indexer
}

// NewHTTPRouteLister returns a new HTTPRouteLister.
func NewHTTPRouteLister(indexer cache.Indexer) HTTPRouteLister {
	return &hTTPRouteLister{indexer: indexer}
}

// List lists all HTTPRoutes in the indexer.
func (s *hTTPRouteLister) List(selector labels.Selector) (ret []*v1alpha2.HTTPRoute, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.HTTPRoute))
	})
	return ret, err
}

// HTTPRoutes returns an object that can list and get HTTPRoutes.
func (s *hTTPRouteLister) HTTPRoutes(namespace string) HTTPRouteNamespaceLister {
	return hTTPRouteNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// HTTPRouteNamespaceLister helps list and get HTTPRoutes.
// All objects returned here must be treated as read-only.
type HTTPRouteNamespaceLister interface {
	// List lists all HTTPRoutes in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha2.HTTPRoute, err error)
	// Get retrieves the HTTPRoute from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha2.HTTPRoute, error)
	HTTPRouteNamespaceListerExpansion
}

// hTTPRouteNamespaceLister implements the HTTPRouteNamespaceLister
// interface.
type hTTPRouteNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all HTTPRoutes in the indexer for a given namespace.
func (s hTTPRouteNamespaceLister) List(selector labels.Selector) (ret []*v1alpha2.HTTPRoute, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.HTTPRoute))
	})
	return ret, err
}

// Get retrieves the HTTPRoute from the indexer for a given namespace and name.
func (s hTTPRouteNamespaceLister) Get(name string) (*v1alpha2.HTTPRoute, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha2.Resource("httproute"), name)
	}
	return obj.(*v1alpha2.HTTPRoute), nil
}
