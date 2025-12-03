package k8s

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	clientset     *kubernetes.Clientset
	metricsClient *metricsv.Clientset
	config        *rest.Config
	context       string
	namespace     string
}

func NewClient() (*Client, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes config: %w", err)
		}
	}

	config.Timeout = 30 * time.Second

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	metricsClient, _ := metricsv.NewForConfig(config)

	rawConfig, _ := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	currentContext := ""
	if rawConfig != nil {
		currentContext = rawConfig.CurrentContext
	}

	return &Client{
		clientset:     clientset,
		metricsClient: metricsClient,
		config:        config,
		context:       currentContext,
		namespace:     "default",
	}, nil
}

func (c *Client) Clientset() *kubernetes.Clientset {
	return c.clientset
}

func (c *Client) MetricsClient() *metricsv.Clientset {
	return c.metricsClient
}

func (c *Client) Context() string {
	return c.context
}

func (c *Client) Namespace() string {
	return c.namespace
}

func (c *Client) SetNamespace(ns string) {
	c.namespace = ns
}

func (c *Client) ListNamespaces(ctx context.Context) ([]string, error) {
	return ListNamespaces(ctx, c.clientset)
}

func (c *Client) ListContexts() ([]string, string, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := rules.Load()
	if err != nil {
		return nil, "", err
	}

	var contexts []string
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}
	return contexts, config.CurrentContext, nil
}
