package main

import (
	"context"
	"fmt"
	"github.com/akamensky/argparse"
	iapi "github.com/eranco74/inventory/pkg/apis"
	nv1alpha1 "github.com/eranco74/inventory/pkg/apis/eranco74/v1alpha1"
	"io/ioutil"
	runtime "k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	parser := argparse.NewParser("health", "Health check")
	// Create string flag
	ip := parser.String("i", "ip", &argparse.Options{Required: true, Help: "Ip to connect"})
	port := parser.String("p", "port", &argparse.Options{Required: true, Help: "port to connect"})
	namespace := parser.String("s", "namespace", &argparse.Options{Required: true, Help: "Machine namespace"})
	name := parser.String("n", "name", &argparse.Options{Required: true, Help: "Machine name"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	cl := GetClient()
	key := client.ObjectKey{*namespace, *name}
	fmt.Println("Helathcheck loop started")
	for {
		time.Sleep(15 * time.Second)
		go checkMachineHealth(fmt.Sprintf("http://%s:%s/hello", *ip, *port), cl, key)
	}
}

func checkMachineHealth(connectionString string, cl client.Client, key client.ObjectKey) string {
	resp, err := http.Get(connectionString)
	if err != nil {
		fmt.Println("Healthcheck failed")
		fmt.Println(err.Error())
		updateHealth(cl, key, "Unknown")
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(fmt.Sprintf("RESPONSE %s", body))
	updateHealth(cl, key, "Ready")

	return string(body)
}

func GetClient() client.Client {
	scheme := runtime.NewScheme()

	clientgoscheme.AddToScheme(scheme)

	iapi.AddToScheme(scheme)
	log.Println("Creating client")
	cl, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		log.Println("Failed to create client")
		os.Exit(1)
	}
	return cl
}

func updateHealth(cl client.Client, key client.ObjectKey, status string) error {

	machine := &nv1alpha1.MachineHealth{}
	err := cl.Get(context.Background(), key, machine)
	if machine.Spec.MachineHelath == status {
		log.Println("Nothing to update")
		return nil
	}
	if err != nil {
		log.Printf("Failed to get machinehealth with key: %s %v\n", key, err)
		return err
	}
	log.Printf("Updateing MachineHealth from: %s to: %s", machine.Spec.MachineHelath, status)

	updated := machine.DeepCopy()
	updated.Spec.MachineHelath = status
	err = cl.Update(context.Background(), updated)
	if err != nil {
		fmt.Printf("failed to update nodes in namespace default: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Updated")
	return nil
}
