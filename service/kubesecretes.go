package service

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/karim-w/ksec/models"

	v1 "k8s.io/api/core/v1"
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
	if conf.EnvPath != "" {
		handleEnvPath(conf.EnvPath, conf.Namespace, conf.Secret)
		return
	}
	if conf.Set {
		addKubeSecret(conf.Namespace, conf.Secret, conf.Key, conf.Value)
		return
	}
	if conf.Get {
		getkubeSecretValue(conf.Namespace, conf.Secret, conf.Key)
		return
	}
	if conf.Delete {
		deleteKubeSecret(conf.Namespace, conf.Secret, conf.Key)
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

func deleteKubeSecret(namespace string, secret string, key string) {
	// get secret
	s, err := GetKubeClient().CoreV1().Secrets(namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	delete(s.Data, key)
	// delete Key through update
	s, err = GetKubeClient().CoreV1().Secrets(namespace).Update(context.TODO(), s, metav1.UpdateOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	for k, v := range s.Data {
		fmt.Println(k, string(v))
	}
}

func extractValuesFromEnv(path string, namespace string, secret string) *map[string][]byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sList := make(map[string][]byte)

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			//get location of first "="
			index := strings.Index(line, "=")
			//get key and value
			key := line[:index]
			value := line[index+1:]
			sList[key] = []byte(value)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &sList
}

func handleEnvPath(path string, namespace string, secret string) {
	vals := extractValuesFromEnv(path, namespace, secret)
	//Create secret
	s, err := GetKubeClient().CoreV1().Secrets(namespace).Create(
		context.TODO(),
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secret,
				Namespace: namespace,
			},
			Data: *vals,
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		println(err.Error())
		return
	}
	generateDecalrationFile(secret, &s.Data)
	fmt.Println("Secret created")
}

func generateDecalrationFile(secret string, vals *map[string][]byte) {
	//create file
	file, err := os.Create(secret + ".yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	overall := "env:"
	//write to file
	for k, _ := range *vals {
		line := fmt.Sprintf("\n- name: %s\n\tvalueFrom:\n\t\tsecretKeyRef:\n\t\t\tname: %s\n\t\t\tkey: %s", k, secret, k)
		overall += line
	}
	_, err = file.WriteString(overall)
	if err != nil {
		log.Fatal(err)
	}
}
