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

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/warroyo/nsxt-tag-controller/controllers"
	"github.com/warroyo/nsxt-tag-controller/nsxt"
	runtanzuvmwarecomv1alpha1 "gitlab.eng.vmware.com/core-build/guest-cluster-controller/apis/run.tanzu/v1alpha1"
	vmwarecomv1alpha1 "gitlab.eng.vmware.com/core-build/nsx-ujo/k8s-virtual-networking-client/pkg/apis/k8svirtualnetworking/v1alpha1"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = runtanzuvmwarecomv1alpha1.AddToScheme(scheme)

	_ = vmwarecomv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var nsxthost string
	var nsxtpass string
	var nsxtuser string

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&nsxthost, "nsxt-host", os.Getenv("NSXT_HOST"), "the nsxt mgr hostname or ip , NSXT_HOST env var")
	flag.StringVar(&nsxtpass, "nsxt-password", os.Getenv("NSXT_PASSWORD"), "the password for the nsxt manager,  NSXT_PASSOWRD env var")
	flag.StringVar(&nsxtuser, "nsxt-username", os.Getenv("NSXT_USERNAME"), "the username for the nsxt manager,  NSXT_USERNAME env var")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "8e47160d.field.vmware.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	nsxclient := nsxt.NsxtClient{
		Host: nsxthost,
		User: nsxtuser,
		Pass: nsxtpass,
	}
	ctrl.Log.Info("connecting to nsxt")
	err = nsxt.ConfigurePolicyConnectorData(&nsxclient)
	if err != nil {
		setupLog.Error(err, "unable to setup NSXT")
		os.Exit(1)
	}
	if err = (&controllers.TanzuKubernetesClusterReconciler{
		Client:     mgr.GetClient(),
		Log:        ctrl.Log.WithName("controllers").WithName("TanzuKubernetesCluster"),
		Scheme:     mgr.GetScheme(),
		NsxtClient: nsxclient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TanzuKubernetesCluster")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
