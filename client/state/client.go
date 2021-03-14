package state

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	apiv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/api/v1"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Client is the interface for retrieving and storing state
type Client interface {
	ReadState(ctx context.Context) (state *apiv1.State, err error)
	StoreState(ctx context.Context, state apiv1.State) (err error)
}

// NewClient returns new bigquery.Client
func NewClient(kubeClientset *kubernetes.Clientset, stateFilePath, stateFileConfigMapName string) (Client, error) {
	return &client{
		kubeClientset:          kubeClientset,
		stateFilePath:          stateFilePath,
		stateFileConfigMapName: stateFileConfigMapName,
	}, nil
}

type client struct {
	kubeClientset          *kubernetes.Clientset
	stateFilePath          string
	stateFileConfigMapName string
}

func (c *client) ReadState(ctx context.Context) (state *apiv1.State, err error) {

	// check if last measurement file exists in configmap
	if _, err := os.Stat(c.stateFilePath); !os.IsNotExist(err) {
		log.Info().Msgf("File %v exists, reading contents...", c.stateFilePath)

		// read state file
		data, err := ioutil.ReadFile(c.stateFilePath)
		if err != nil {
			return nil, fmt.Errorf("Failed reading file from path %v: %w", c.stateFilePath, err)
		}

		log.Info().Msgf("Unmarshalling file %v contents...", c.stateFilePath)

		// unmarshal state file
		if err := json.Unmarshal(data, &state); err != nil {
			return nil, fmt.Errorf("Failed unmarshalling last measurement file: %w", err)
		}
	}

	return
}

func (c *client) StoreState(ctx context.Context, state apiv1.State) (err error) {

	currentNamespace, err := c.getCurrentNamespace()
	if err != nil {
		return
	}

	// retrieve configmap
	configMap, err := c.kubeClientset.CoreV1().ConfigMaps(currentNamespace).Get(ctx, c.stateFileConfigMapName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Failed retrieving configmap %v: %w", c.stateFileConfigMapName, err)
	}

	// marshal state to json
	stateData, err := json.Marshal(state)
	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}

	configMap.Data[filepath.Base(c.stateFilePath)] = string(stateData)

	// update configmap to have measurement available when the application runs the next time and for other applications
	_, err = c.kubeClientset.CoreV1().ConfigMaps(currentNamespace).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Failed updating configmap %v: %w", c.stateFileConfigMapName, err)
	}

	log.Info().Msgf("Stored state in configmap %v...", c.stateFileConfigMapName)

	return nil
}

func (c *client) getCurrentNamespace() (namespace string, err error) {
	ns, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return namespace, fmt.Errorf("Failed reading namespace: %w", err)
	}

	return string(ns), nil
}
