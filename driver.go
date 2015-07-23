package main

import "net/http"
import "log"
import "flag"

type Config struct {
	Sock           string
	StateFile      string
	ClusterAddress string
	MountBase      string
}

const dockerVersionMimeType = "application/vnd.docker.plugins.v1+json"

func pluginActivate(w http.ResponseWriter, r *http.Request)      {}
func pluginVolumeCreate(w http.ResponseWriter, r *http.Request)  {}
func pluginVolumeRemove(w http.ResponseWriter, r *http.Request)  {}
func pluginVolumeMount(w http.ResponseWriter, r *http.Request)   {}
func pluginVolumeUnmount(w http.ResponseWriter, r *http.Request) {}
func pluginVolumePath(w http.ResponseWriter, r *http.Request)    {}

func main() {
	v := new(Config)

	// parse command line.
	flag.StringVar(&v.Sock, "sock", "/run/docker/plugins/springpath", "socket docker talks to")
	flag.StringVar(&v.StateFile, "statefile", "/SYSTEM/volume-driver.json", "springpath volume driver metadata")
	flag.StringVar(&v.ClusterAddress, "clusteraddress", "localhost", "address of the springpath i/o dispatcher")
	flag.StringVar(&v.MountBase, "mountbase", "/run/springpath-docker-volumes", "base path for springpath volume mount points")
	flag.Parse()

	log.Println("starting docker volume plugin")

	// http server.
	http.HandleFunc("/Plugin.Activate", pluginActivate)
	http.HandleFunc("/VolumeDriver.Create", pluginVolumeCreate)
	http.HandleFunc("/VolumeDriver.Remove", pluginVolumeRemove)
	http.HandleFunc("/VolumeDriver.Mount", pluginVolumeMount)
	http.HandleFunc("/VolumeDriver.Unmount", pluginVolumeUnmount)
	http.HandleFunc("/VolumeDriver.Path", pluginVolumePath)

	err := http.ListenAndServe(v.Sock, nil)
	if err != nil {
		log.Fatal("failed to start http server", err)
	}
}
