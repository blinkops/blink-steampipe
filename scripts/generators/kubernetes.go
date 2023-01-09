package generators

import (
	"fmt"
	"github.com/blinkops/blink-steampipe/scripts/consts"
	"os"
	"strconv"

	"github.com/ghodss/yaml"
	uuid "github.com/satori/go.uuid"
)

const (
	kubernetesConnectionIdentifier = "KUBERNETES_CONNECTION"
	kubernetesApiUrl               = "KUBERNETES_API_URL"
	kubernetesBearerToken          = "KUBERNETES_BEARER_TOKEN"
	kubernetesVerifyCertificate    = "KUBERNETES_VERIFY_CERTIFICATE"
	kubeConfigDirectoryPath        = consts.SteampipeBasePath + ".kube/"
	kubeConfigFilePath             = kubeConfigDirectoryPath + "config"
)

type KubernetesCredentialGenerator struct{}

func (gen KubernetesCredentialGenerator) Generate() error {
	if _, ok := os.LookupEnv(kubernetesConnectionIdentifier); !ok {
		return nil
	}

	apiUrl, ok := os.LookupEnv(kubernetesApiUrl)
	if !ok {
		return fmt.Errorf("invalid kubernetes connection was provided")
	}

	token, ok := os.LookupEnv(kubernetesBearerToken)
	if !ok {
		return fmt.Errorf("invalid kubernetes connection was provided")
	}

	var verifyCert bool
	if verify, ok := os.LookupEnv(kubernetesVerifyCertificate); ok {
		if verifyAsBool, err := strconv.ParseBool(verify); err == nil {
			verifyCert = verifyAsBool
		}
	}

	if err := os.MkdirAll(kubeConfigDirectoryPath, 0o770); err != nil {
		return fmt.Errorf("error creating configuration path: %v", err)
	}

	ctxName, clusterName, userName := uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String()
	configuration := KubectlConfig{
		Kind:           "Config",
		ApiVersion:     "v1",
		CurrentContext: ctxName,
		Clusters: []*KubectlClusterWithName{
			{
				Name:    clusterName,
				Cluster: KubectlCluster{Server: apiUrl, InsecureSkipTLSVerify: !verifyCert},
			},
		},
		Users: []*KubectlUserWithName{
			{
				Name: userName,
				User: KubectlUser{Token: token},
			},
		},
		Contexts: []*KubectlContextWithName{
			{
				Name:    ctxName,
				Context: KubectlContext{User: userName, Cluster: clusterName},
			},
		},
	}

	manifest, err := yaml.Marshal(configuration)
	if err != nil {
		return fmt.Errorf("error marshaling authentication config to yaml: %v", err)
	}

	if err = os.WriteFile(kubeConfigFilePath, manifest, 0o600); err != nil {
		return fmt.Errorf("unable to prepare kuberenetes credentials: %w", err)
	}

	return nil
}
