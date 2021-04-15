/*


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

package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/warroyo/nsxt-tag-controller/nsxt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runtanzuvmwarecomv1alpha1 "gitlab.eng.vmware.com/core-build/guest-cluster-controller/apis/run.tanzu/v1alpha1"
	vmwarecomv1alpha1 "gitlab.eng.vmware.com/core-build/nsx-ujo/k8s-virtual-networking-client/pkg/apis/k8svirtualnetworking/v1alpha1"
)

// TanzuKubernetesClusterReconciler reconciles a TanzuKubernetesCluster object
type TanzuKubernetesClusterReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	NsxtClient nsxt.NsxtClient
}

// +kubebuilder:rbac:groups=run.tanzu.vmware.com,resources=tanzukubernetesclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=run.tanzu.vmware.com,resources=tanzukubernetesclusters/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=vmware.com,resources=virtualnetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmware.com,resources=virtualnetworks/status,verbs=get;update;patch

func (r *TanzuKubernetesClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("tanzukubernetescluster", req.NamespacedName)

	var tkCluster runtanzuvmwarecomv1alpha1.TanzuKubernetesCluster
	if err := r.Get(ctx, req.NamespacedName, &tkCluster); err != nil {
		r.Log.Error(err, "unable to fetch Cluster")

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cluster := tkCluster.Name
	namespacedVnet := types.NamespacedName{Namespace: req.Namespace, Name: fmt.Sprintf("%s-vnet", cluster)}
	var virtualNetwork vmwarecomv1alpha1.VirtualNetwork
	if err := r.Get(ctx, namespacedVnet, &virtualNetwork); err != nil {
		r.Log.Error(err, "unable to fetch virtualNetwork")

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	nsxsegmentid := fmt.Sprintf("vnet_%s_0", virtualNetwork.UID)
	connector := nsxt.GetPolicyConnector(r.NsxtClient)
	client := infra.NewDefaultSegmentsClient(connector)
	segment, err := client.Get(nsxsegmentid)
	if err != nil {
		r.Log.Info(err.Error())
	}
	tags := segment.Tags
	for key, value := range tkCluster.Labels {
		if strings.Contains(key, "policytag/") {
			scope := strings.ReplaceAll(key, "policytag/", "")
			tag := value
			item := model.Tag{
				Scope: &scope,
				Tag:   &tag,
			}
			add := true
			for _, i := range segment.Tags {
				if scope == *i.Scope && tag == *i.Tag {
					add = false
				}
			}
			if add {
				tags = append(tags, item)
			}

		}
	}

	if len(tags) > len(segment.Tags) {
		r.Log.Info(fmt.Sprintf("updating tags on %s", cluster))
		obj := model.Segment{
			Tags: tags,
		}

		err = client.Patch(nsxsegmentid, obj)
		if err != nil {
			r.Log.Info(err.Error())
		}
	}

	return ctrl.Result{}, nil
}

func (r *TanzuKubernetesClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&runtanzuvmwarecomv1alpha1.TanzuKubernetesCluster{}).
		Complete(r)
}
