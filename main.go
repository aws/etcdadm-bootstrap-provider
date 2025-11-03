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

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	capiflags "sigs.k8s.io/cluster-api/util/flags"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	bootstrapv1alpha3 "github.com/aws/etcdadm-bootstrap-provider/api/v1alpha3"
	bootstrapv1beta1 "github.com/aws/etcdadm-bootstrap-provider/api/v1beta1"
	"github.com/aws/etcdadm-bootstrap-provider/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme               = runtime.NewScheme()
	setupLog             = ctrl.Log.WithName("setup")
	watchNamespace       string
	managerOptions       capiflags.ManagerOptions
	enableLeaderElection bool
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = clusterv1.AddToScheme(scheme)
	_ = bootstrapv1alpha3.AddToScheme(scheme)
	_ = bootstrapv1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create
func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	capiflags.AddManagerOptions(pflag.CommandLine, &managerOptions)

	pflag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	pflag.StringVar(&watchNamespace, "namespace", "",
		"Namespace that the controller watches to reconcile etcdadmConfig objects. If unspecified, the controller watches forobjects across all namespaces.")

	pflag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	_, metricsServerOpts, err := capiflags.GetManagerOptions(managerOptions)
	if err != nil {
		setupLog.Error(err, "Unable to start manager: invalid metrics server flags")
		os.Exit(1)
	}

	opts := ctrl.Options{
		Scheme:           scheme,
		Metrics:          *metricsServerOpts,
		LeaderElection:   enableLeaderElection,
		LeaderElectionID: "e949c6e2.cluster.x-k8s.io",
	}

	if watchNamespace != "" {
		opts.Cache = cache.Options{
			DefaultNamespaces: map[string]cache.Config{watchNamespace: {}},
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), opts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()
	if err = (&controllers.EtcdadmConfigReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("EtcdadmConfig"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "EtcdadmConfig")
		os.Exit(1)
	}

	if err = (&bootstrapv1beta1.EtcdadmConfig{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "EtcdadmConfig")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
