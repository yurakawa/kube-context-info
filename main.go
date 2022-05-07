package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type kubeContext struct {
	Namespace string `yaml:"namespace"`
}

type kubeContextWithName struct {
	Context kubeContext `yaml:"context"`
	Name    string      `yaml:"name"`
}
type kubeConfig struct {
	APIVersion     string                `yaml:"apiVersion"`
	CurrentContext string                `yaml:"current-context"`
	Contexts       []kubeContextWithName `yaml:"contexts"`
}

func main() {
	err := displayContextInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
func displayContextInfo() error {
	usr, _ := user.Current()
	home := usr.HomeDir
	path := filepath.Join(home, ".kube", "config")
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open the file: %s", err.Error())
	}
	buf, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read the file: %w", err)
	}
	var config kubeConfig
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the config: %w", err)
	}

	var namespace string
	for _, ctxWithName := range config.Contexts {
		if ctxWithName.Name == config.CurrentContext {
			namespace = ctxWithName.Context.Namespace
		}
	}
	fmt.Printf("%s/%s", config.CurrentContext, namespace)
	return nil
}
