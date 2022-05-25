/*
 * Copyright 2019-2020 VMware, Inc.
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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	avicache "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/cache"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/k8s"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/lib"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/api"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/api/models"
	crd "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1alpha1/clientset/versioned"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/utils"
	advl4 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/third_party/service-apis/client/clientset/versioned"

	oshiftclient "github.com/openshift/client-go/route/clientset/versioned"
	istiocrd "istio.io/client-go/pkg/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	svcapi "sigs.k8s.io/service-apis/pkg/client/clientset/versioned"
)

var (
	masterURL  string
	kubeconfig string
	version    = "dev"
)

func main() {

	InitializeAKOApi()

	InitializeAKC()
}

func InitializeAKOApi() {
	akoApi := api.NewServer(lib.GetAkoApiServerPort(), []models.ApiModel{})
	akoApi.InitApi()
	lib.SetApiServerInstance(akoApi)
}

func InitializeAKC() {
	var err error
	kubeCluster := false
	utils.AviLog.Info("AKO is running with version: ", version)
	// Check if we are running inside kubernetes. Hence try authenticating with service token
	cfg, err := rest.InClusterConfig()
	if err != nil {
		utils.AviLog.Warnf("We are not running inside kubernetes cluster. %s", err.Error())
	} else {
		utils.AviLog.Info("We are running inside kubernetes cluster. Won't use kubeconfig files.")
		kubeCluster = true
	}

	if kubeCluster == false {
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
		utils.AviLog.Infof("master: %s", masterURL)
		if err != nil {
			utils.AviLog.Fatalf("Error building kubeconfig: %s", err.Error())
		}
	}

	// Initialize akoControlConfig
	akoControlConfig := lib.AKOControlConfig()
	//Used to set vrf context, static routes
	isPrimaryAKO, _ := strconv.ParseBool(os.Getenv("PRIMARY_AKO_FLAG"))
	akoControlConfig.SetAKOInstanceFlag(isPrimaryAKO)
	var crdClient *crd.Clientset
	var advl4Client *advl4.Clientset
	var svcAPIClient *svcapi.Clientset
	var istioClient *istiocrd.Clientset
	if lib.IsWCP() {
		advl4Client, err = advl4.NewForConfig(cfg)
		if err != nil {
			utils.AviLog.Fatalf("Error building service-api v1alpha1pre1 clientset: %s", err.Error())
		}
		akoControlConfig.SetAdvL4Clientset(advl4Client)
	} else {
		if lib.UseServicesAPI() {
			svcAPIClient, err = svcapi.NewForConfig(cfg)
			if err != nil {
				utils.AviLog.Fatalf("Error building service-api clientset: %s", err.Error())
			}
			akoControlConfig.SetServicesAPIClientset(svcAPIClient)
		}

		crdClient, err = crd.NewForConfig(cfg)
		if err != nil {
			utils.AviLog.Fatalf("Error building AKO CRD clientset: %s", err.Error())
		}
		akoControlConfig.SetCRDClientset(crdClient)

		// Handle Istio code.
		if lib.IsIstioEnabled() {
			istioClient, err = istiocrd.NewForConfig(cfg)
			if err != nil {
				utils.AviLog.Fatalf("Error building Istio CRD clientset: %s", err.Error())
			}
			akoControlConfig.SetIstioClientset(istioClient)
		}
	}

	dynamicClient, err := lib.NewDynamicClientSet(cfg)
	if err != nil {
		utils.AviLog.Warnf("Error while creating dynamic client %v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		utils.AviLog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	akoControlConfig.SetEventRecorder(lib.AKOEventComponent, kubeClient, false)
	pod, err := kubeClient.CoreV1().Pods(utils.GetAKONamespace()).Get(context.TODO(), os.Getenv("POD_NAME"), metav1.GetOptions{})
	if err != nil {
		utils.AviLog.Warnf("Error getting AKO pod details, %s.", err.Error())
	}
	akoControlConfig.SaveAKOPodObjectMeta(pod)

	// Check for kubernetes apiserver version compatibility with AKO version.
	if serverVersionInfo, err := kubeClient.Discovery().ServerVersion(); err != nil {
		utils.AviLog.Warnf("Error while fetching kubernetes apiserver version")
	} else {
		serverVersion := fmt.Sprintf("%s.%s", serverVersionInfo.Major, serverVersionInfo.Minor)
		utils.AviLog.Infof("Kubernetes cluster apiserver version %s", serverVersion)
		if lib.CompareVersions(serverVersion, ">", lib.GetK8sMaxSupportedVersion()) ||
			lib.CompareVersions(serverVersion, "<", lib.GetK8sMinSupportedVersion()) {
			akoControlConfig.PodEventf(corev1.EventTypeWarning, lib.AKOShutdown, "Unsupported kubernetes apiserver %s version detected", serverVersion)
			utils.AviLog.Fatalf("Unsupported kubernetes apiserver version detected. Please check the supportability guide.")
		}
	}

	oshiftClient, err := oshiftclient.NewForConfig(cfg)
	if err != nil {
		utils.AviLog.Warnf("Error in creating openshift clientset")
	}

	registeredInformers, err := lib.InformersToRegister(kubeClient, oshiftClient, false)
	if err != nil {
		utils.AviLog.Fatalf("Failed to initialize informers: %v, shutting down AKO, going to reboot", err)
	}

	informersArg := make(map[string]interface{})
	informersArg[utils.INFORMERS_OPENSHIFT_CLIENT] = oshiftClient
	informersArg[utils.INFORMERS_AKO_CLIENT] = crdClient

	if lib.GetNamespaceToSync() != "" {
		informersArg[utils.INFORMERS_NAMESPACE] = lib.GetNamespaceToSync()
	}

	// Namespace bound Secret informers should be initialized for AKO in VDS,
	// For AKO in VCF, we will need to watch on Secrets across all namespaces.
	if !utils.IsVCFCluster() && lib.GetAdvancedL4() {
		informersArg[utils.INFORMERS_ADVANCED_L4] = true
	}

	utils.NewInformers(utils.KubeClientIntf{ClientSet: kubeClient}, registeredInformers, informersArg)
	lib.NewDynamicInformers(dynamicClient, false)
	if lib.IsWCP() {
		k8s.NewAdvL4Informers(advl4Client)
	} else {
		k8s.NewCRDInformers(crdClient)
		if lib.UseServicesAPI() {
			k8s.NewSvcApiInformers(svcAPIClient)
		}
	}

	// Set Istio Informers.
	if lib.IsIstioEnabled() {
		k8s.NewIstioCRDInformers(istioClient)
	}

	informers := k8s.K8sinformers{Cs: kubeClient, DynamicClient: dynamicClient, OshiftClient: oshiftClient}
	c := k8s.SharedAviController()
	stopCh := utils.SetupSignalHandler()
	ctrlCh := make(chan struct{})
	quickSyncCh := make(chan struct{})

	if utils.IsVCFCluster() {
		c.InitVCFHandlers(kubeClient, ctrlCh, stopCh)
	}

	aviRestClientPool := avicache.SharedAVIClients()
	if aviRestClientPool == nil {
		utils.AviLog.Fatalf("Avi client not initialized")
	}

	if aviRestClientPool != nil && !avicache.IsAviClusterActive(aviRestClientPool.AviClient[0]) {
		akoControlConfig.PodEventf(corev1.EventTypeWarning, lib.AKOShutdown, "Avi Controller Cluster state is not Active")
		utils.AviLog.Fatalf("Avi Controller Cluster state is not Active, shutting down AKO")
	}

	akoControlConfig.SetLicenseType(aviRestClientPool.AviClient[0])

	err = c.HandleConfigMap(informers, ctrlCh, stopCh, quickSyncCh)
	if err != nil {
		utils.AviLog.Errorf("Handle configmap error during reboot, shutting down AKO. Error is: %v", err)
		return
	}

	if _, err := lib.GetVipNetworkListEnv(); !utils.IsVCFCluster() && err != nil {
		utils.AviLog.Fatalf("Error in getting VIP network %s, shutting down AKO", err)
	}

	c.InitializeNamespaceSync()
	k8s.PopulateNodeCache(kubeClient)
	waitGroupMap := make(map[string]*sync.WaitGroup)
	wgIngestion := &sync.WaitGroup{}
	waitGroupMap["ingestion"] = wgIngestion
	wgFastRetry := &sync.WaitGroup{}
	waitGroupMap["fastretry"] = wgFastRetry
	wgSlowRetry := &sync.WaitGroup{}
	waitGroupMap["slowretry"] = wgSlowRetry
	wgGraph := &sync.WaitGroup{}
	waitGroupMap["graph"] = wgGraph
	wgStatus := &sync.WaitGroup{}
	waitGroupMap["status"] = wgStatus
	go c.InitController(informers, registeredInformers, ctrlCh, stopCh, quickSyncCh, waitGroupMap)
	<-stopCh
	close(ctrlCh)
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		wgIngestion.Wait()
		wgGraph.Wait()
		wgFastRetry.Wait()
		wgStatus.Wait()
	}()
	// Timeout after 60 seconds.
	timeout := 60 * time.Second
	select {
	case <-doneChan:
		return
	case <-time.After(timeout):
		utils.AviLog.Warnf("Timed out while waiting for threads to return, going to stop AKO. Time waited 60 seconds")
		return
	}

}

func init() {
	def_kube_config := os.Getenv("HOME") + "/.kube/config"
	flag.StringVar(&kubeconfig, "kubeconfig", def_kube_config, "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
