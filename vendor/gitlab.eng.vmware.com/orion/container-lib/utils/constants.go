/*
 * [2013] - [2018] Avi Networks Incorporated
 * All Rights Reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *   http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

const (
	GraphLayer            = "GraphLayer"
	ObjectIngestionLayer  = "ObjectIngestionLayer"
	LeastConnection       = "LB_ALGORITHM_LEAST_CONNECTIONS"
	RandomConnection      = "RANDOM_CONN"
	PassthroughConnection = "PASSTHROUGH_CONN"
	RoundRobinConnection  = "LB_ALGORITHM_ROUND_ROBIN"
	ServiceInformer       = "ServiceInformer"
	PodInformer           = "PodInformer"
	SecretInformer        = "SecretInformer"
	EndpointInformer      = "EndpointInformer"
	IstioMutualKey        = "key.pem"
	IstioMutualCertChain  = "cert-chain.pem"
	IstioMutualRootCA     = "root-cert.pem"
	IngressInformer       = "IngressInformer"
	RouteInformer         = "RouteInformer"
	L4LBService           = "L4LBService"
	LoadBalancer          = "LoadBalancer"
	Endpoints             = "Endpoints"
	Ingress               = "Ingress"
	Service               = "Service"
)
