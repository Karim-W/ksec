package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/karim-w/ksec/models"
	// v1 "k8s.io/client-go/applyconfigurations/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var instance *kubernetes.Clientset
var mtx sync.Mutex

func boot() *kubernetes.Clientset {
	mtx.Lock()
	defer mtx.Unlock()
	if instance == nil {
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
		config, err := kubeconfig.ClientConfig()
		if err != nil {
			panic(err)
		}
		instance = kubernetes.NewForConfigOrDie(config)
	}
	return instance
}

func GetKubeClient() *kubernetes.Clientset {
	return boot()
}

func KubectlSecretsSvc(conf *models.Secrets) {
	if conf.Set {
		addKubeSecret(conf.Namespace, conf.Secret, conf.Key, conf.Value)
		return
	}
	if conf.Get {
		getkubeSecretValue(conf.Namespace, conf.Secret, conf.Key)
		return
	}
	if conf.Delete {
		fmt.Println("Delete")
		return
	}
	if conf.List {
		listKubeSecrets(conf.Namespace, conf.Secret)
		return
	}
}

func getkubeSecretValue(namespace string, secret string, key string) {
	s, err := GetKubeClient().CoreV1().Secrets(namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(string(s.Data[key]))
}

func listKubeSecrets(namespace string, secret string) {
	s, err := GetKubeClient().CoreV1().Secrets(namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	for k, v := range s.Data {
		fmt.Println(k, string(v))
	}
}

func addKubeSecret(namespace string, secret string, key string, value string) {
	//get secrets
	s, err := GetKubeClient().CoreV1().Secrets(namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	s.Data[key] = []byte(value)
	//update secret
	s, err = GetKubeClient().CoreV1().Secrets(namespace).Update(context.TODO(), s, metav1.UpdateOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	for k, v := range s.Data {
		fmt.Println(k, string(v))
	}
}
