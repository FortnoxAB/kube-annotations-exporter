package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/koding/multiconfig"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var annotationsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "kube",
	Name:      "exporter",
	Subsystem: "annotations",
}, []string{"kind", "name", "annotation"})

type config struct {
	Port        int
	Annotations []string
}

var objects = map[string]func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error){
	"Deployment": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
	"Pod": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
	"Service": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
	"Ingress": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
	"ConfigMap": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
	"Secret": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
	"StatefulSet": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
	"NetworkPolicy": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.NetworkingV1().NetworkPolicies("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
	"CronJob": func(ctx context.Context, clientset *kubernetes.Clientset) ([]metav1.Object, error) {
		list, err := clientset.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var tmp []metav1.Object
		for _, item := range list.Items {
			item := item //magic
			tmp = append(tmp, item.GetObjectMeta())
		}

		return tmp, nil
	},
}

type exporter struct {
	clientset   *kubernetes.Clientset
	annotations []string
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	annotationsGauge.Describe(ch)
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	annotationsGauge.Reset()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	err := checkForAnnotations(ctx, e.clientset, e.annotations)

	if err != nil {
		slog.Info("failed to check ignored items", "error", err)
	}
	annotationsGauge.Collect(ch)
}

func main() {
	config := &config{
		Port: 8080,
	}
	multiconfig.MustLoad(config)

	clientset, err := GetClientset()
	if err != nil {
		log.Fatalf("failed to get clientset: %s", err)
	}

	prometheus.MustRegister(&exporter{
		clientset:   clientset,
		annotations: config.Annotations,
	})
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting webserver", "port", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), promhttp.Handler())
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %s", err)
	}
}

func checkForAnnotations(ctx context.Context, clientset *kubernetes.Clientset, annotations []string) error {
	for kind, getList := range objects {
		list, err := getList(ctx, clientset)
		if err != nil {
			return err
		}

		for _, item := range list {
			for _, annotation := range annotations {
				if _, ok := item.GetAnnotations()[annotation]; ok {
					annotationsGauge.WithLabelValues(kind, item.GetName(), annotation).Set(float64(1))
				}
			}
		}
	}
	return nil
}

func GetClientset() (*kubernetes.Clientset, error) {
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		slog.Info("No kubeconfig found. Using incluster...")
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		return kubernetes.NewForConfig(config)
	}

	return kubernetes.NewForConfig(config)
}
