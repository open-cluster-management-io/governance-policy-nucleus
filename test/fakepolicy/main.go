// Copyright Contributors to the Open Cluster Management project

package fakepolicy

import (
	"context"
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	policyv1beta1 "open-cluster-management.io/governance-policy-nucleus/test/fakepolicy/api/v1beta1"
	"open-cluster-management.io/governance-policy-nucleus/test/fakepolicy/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
	utilruntime.Must(policyv1beta1.AddToScheme(scheme))
}

//nolint:unused
func main() {
	if err := Run(context.Background(), nil); err != nil {
		os.Exit(1)
	}
}

func Run(parentCtx context.Context, cfg *rest.Config) error {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string

	flag.StringVar(&metricsAddr, "metrics-bind-address", "0",
		"The address the metric endpoint binds to. Disabled by default, but conventionally :8080")
	flag.StringVar(&probeAddr, "health-probe-bind-address", "0",
		"The address the probe endpoint binds to. Disabled by default, but conventionally :8081")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	opts := zap.Options{
		Development: true,
	}

	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if cfg == nil {
		var err error

		cfg, err = ctrl.GetConfig()
		if err != nil {
			setupLog.Error(err, "unable to get kubernetes config")

			return err
		}
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "8b5e65ab.open-cluster-management.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")

		return err
	}

	if err = (&controllers.FakePolicyReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FakePolicy")

		return err
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")

		return err
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")

		return err
	}

	setupLog.Info("starting manager")

	managerCtx, cancel := context.WithCancel(parentCtx)

	go func() {
		// It would be nicer if this could take a parent context,
		// but this will work to cancel the manager on those signals.
		<-ctrl.SetupSignalHandler().Done()
		cancel()
	}()

	if err := mgr.Start(managerCtx); err != nil {
		setupLog.Error(err, "problem running manager")

		return err
	}

	return nil
}
