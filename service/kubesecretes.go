package service

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/karim-w/ksec/models"
	"gopkg.in/yaml.v3"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	instance *kubernetes.Clientset
	mtx      sync.Mutex
)

func boot() *kubernetes.Clientset {
	mtx.Lock()
	defer mtx.Unlock()
	if instance == nil {
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			rules,
			&clientcmd.ConfigOverrides{},
		)
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
	if conf.FillPath != "" {
		handleFillPath(conf.FillPath, conf.Namespace, conf.Secret)
		return
	}
	if conf.Modify {
		modifyKubeSecret(conf.Namespace, conf.Secret, conf.FileFormat)
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
	s, err := GetKubeClient().CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(string(s.Data[key]))
}

func listKubeSecrets(namespace string, secret string) {
	s, err := GetKubeClient().CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	for k, v := range s.Data {
		fmt.Println(k, string(v))
	}
}

func addKubeSecret(namespace string, secret string, key string, value string) {
	// get secrets
	s, err := GetKubeClient().CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	s.Data[key] = []byte(value)
	// update secret
	s, err = GetKubeClient().CoreV1().
		Secrets(namespace).
		Update(context.TODO(), s, metav1.UpdateOptions{})
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
	s, err := GetKubeClient().CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	delete(s.Data, key)
	// delete Key through update
	s, err = GetKubeClient().CoreV1().
		Secrets(namespace).
		Update(context.TODO(), s, metav1.UpdateOptions{})
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
			// get location of first "="
			index := strings.Index(line, "=")
			// get key and value
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
	// Create secret
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
	// create file
	file, err := os.Create(secret + ".yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	overall := "env:"
	// write to file
	for k := range *vals {
		line := fmt.Sprintf(
			"\n- name: %s\n\tvalueFrom:\n\t\tsecretKeyRef:\n\t\t\tname: %s\n\t\t\tkey: %s",
			k,
			secret,
			k,
		)
		overall += line
	}
	_, err = file.WriteString(overall)
	if err != nil {
		log.Fatal(err)
	}
}

func handleFillPath(path string, namespace string, secret string) {
	// fetch secrets
	s, err := GetKubeClient().CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		println(err.Error())
		return
	}
	secretsMap := make(map[string][]byte)
	for k, v := range s.Data {
		secretsMap[k] = v
	}
	// create the file
	handleGenerateFileWithSecrets(path, &secretsMap)
}

func handleGenerateFileWithSecrets(path string, secrets *map[string][]byte) {
	// create file
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// write to file
	for k, v := range *secrets {
		file.WriteString(k + "=" + string(v) + "\n")
	}
}

func modifyKubeSecret(namespace, secret, format string) {
	// check if file type is supported
	_, ok := models.SupportedFormats[format]
	if !ok {
		fmt.Println(format + " format is not supported")
		return
	}
	// check if secret exists, if not create it
	s, err := GetKubeClient().CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		if err.Error() != `secrets "`+secret+`" not found` {
			println(err.Error())
			return
		}

		// Create generic secret
		s, err = GetKubeClient().CoreV1().Secrets(namespace).Create(
			context.TODO(),
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secret,
					Namespace: namespace,
				},
				Data: map[string][]byte{
					"key": []byte("value"),
				},
			},
			metav1.CreateOptions{},
		)
	}

	m := map[string]interface{}{}

	for k, v := range s.Data {
		m[k] = string(v)
	}

	var out []byte

	switch format {
	case "yaml":
		out, err = yaml.Marshal(m)
	case "json":
		out, err = json.MarshalIndent(m, "", "  ")
	}

	if err != nil {
		println(err.Error())
		return
	}

	file, err := os.Create("." + secret + "." + format)
	if err != nil {
		println(err.Error())
		return
	}

	defer func() {
		defer file.Close()
		// delete file
		err := os.Remove("." + secret + "." + format)
		if err != nil {
			println(err.Error())
		}
	}()

	before_hash := md5.Sum(out)
	before_check_sum := hex.EncodeToString(before_hash[:])

	_, err = file.Write(out)
	if err != nil {
		println(err.Error())
		return
	}

	// open file in editor
	err = openFileInEditor("." + secret + "." + format)
	if err != nil {
		println(err.Error())
		return
	}

	// read file
	file, err = os.Open("." + secret + "." + format)
	if err != nil {
		println(err.Error())
		return
	}

	defer file.Close()

	size, err := file.Stat()
	if err != nil {
		println(err.Error())
		return
	}

	// read file whole
	byts := make([]byte, size.Size())
	n, err := file.Read(byts)
	if err != nil {
		println(err.Error())
		return
	}

	after_hash := md5.Sum(byts[:n])
	after_check_sum := hex.EncodeToString(after_hash[:])

	if before_check_sum == after_check_sum {
		println("No changes detected")
		return
	}

	dat := make(map[string]string)

	// unmarshal
	switch format {
	case "yaml":
		err = yaml.Unmarshal(byts[:n], &dat)
	case "json":
		err = json.Unmarshal(byts[:n], &dat)
	}

	if err != nil {
		fmt.Println("bytssss", string(byts))
		println(err.Error())
		return
	}

	// check if there is a change

	// update secret
	s.Data = make(map[string][]byte)
	for k, v := range dat {
		s.Data[k] = []byte(v)
	}

	_, err = GetKubeClient().CoreV1().
		Secrets(namespace).
		Update(context.TODO(), s, metav1.UpdateOptions{})
	if err != nil {
		println(err.Error())
		return
	}
}

func openFileInEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
