package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	iapi "github.com/eranco74/inventory/pkg/apis"
	nv1alpha1 "github.com/eranco74/inventory/pkg/apis/eranco74/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"sigs.k8s.io/controller-runtime/pkg/client/config"
	"k8s.io/client-go/tools/clientcmd"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func main() {
	log.Println("Starting machine")

	log.Println("Adding to scheme")

	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	iapi.AddToScheme(scheme)

	var kubeconfig *string
	if home := 	os.Getenv("HOME"); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	log.Println("Getting config")

	// use the current context in kubeconfig
	conf, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Creating client")
	cl, err := client.New(conf, client.Options{Scheme: scheme})
	if err != nil {
		log.Println("Failed to create client")
		log.Println(err)
		//os.Exit(1)
	}

	labels := map[string]string{
		"app": "inventory",
	}
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	machine := &nv1alpha1.MachineHealth{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels:    labels,
		},
		Spec: nv1alpha1.MachineHealthSpec{
			getMyIp(),
			"8090",
			5,
			"Started",
		},
	}

	log.Println("Creating machine health in k8s")
	cl.Create(context.Background(), machine)

	key, err := client.ObjectKeyFromObject(machine)
	if err != nil {
		fmt.Printf("failed to get key from object: %s %v\n", key, err)
		os.Exit(1)
	}
	log.Println("My machineHealth key is: ", key)

	log.Println("Serving healthcheck...")
	http.HandleFunc("/hello", handleHello)
	http.ListenAndServe(":8090", nil)
}

func handleHello(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Got healthcheck from %s\n", req.RemoteAddr)
	fmt.Fprintf(w, "OK\n")
}

func getMyIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalln("Oops: " + err.Error() + "\n")
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	log.Fatalln("Oops: failed to find IP")
	return ""
}