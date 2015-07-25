package main

import "net/http"
import "net"
import "log"
import "flag"
import "io/ioutil"
import "encoding/json"

// Global State.
type Config struct {
	Sock           string
	StateFile      string
	ClusterAddress string
	MountBase      string
}

// Individual mounted volumes.
type Volume struct {
	Name          string // docker volume name
	DatastorePath string // springpath datastore name + path
	Size          uint64 // size of the datastore in bytes.
	Mounted       bool   // locally mounted.
	Alive         bool   // scvmclient is reachable.
}

type VolumeDriver interface {
	Create(name string) error
	Remove(name string) error
	Path(name string) (mountpoint string, err error)
	Mount(name string) (mountpoint string, err error)
	Unmount(name string) error
}

// Set of known Volumes.
type VolumeMap struct {
	VolumeDriver
	m map[string]Volume
}

func (v *VolumeMap) Create(name string) error {
	return nil
}

func (v *VolumeMap) Remove(name string) error {
	return nil
}

func (v *VolumeMap) Path(name string) (mountpoint string, err error) {
	return
}

func (v *VolumeMap) Mount(name string) (mountpoint string, err error) {
	return
}

func (v *VolumeMap) Unmount(name string) error {
	return nil
}

const dockerPluginManifestJSON = `
{
	"Implements" : [ "VolumeDriver" ]
}
`

const dockerVersionMimeType = "application/vnd.docker.plugins.v1+json"

func pluginActivate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", dockerVersionMimeType)
	w.Write([]byte(dockerPluginManifestJSON))
}

type PluginHandler struct {
	http.Handler
	VMap *VolumeMap
	Fn   interface{} // sigh.
}

type Message struct {
	Name       string `json:"omitempty"`
	Err        string `json:"omitempty"`
	Mountpoint string `json:"omitempty"`
}

// Common portions of all the plugin endpoints.
func (h PluginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var requestPath = r.URL.Path

	defer r.Body.Close()
	var requestBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("req %p: url %s, body %s\n", r, requestPath, requestBody)

	var request Message
	json.Unmarshal(requestBody, &request)

	return
}

func main() {
	config := new(Config)

	// parse command line.
	flag.StringVar(&config.Sock, "sockpath", "/run/docker/plugins/springpath.sock", "unix domain socket docker talks to")
	flag.StringVar(&config.StateFile, "statefile", "/SYSTEM/volume-driver.json", "springpath volume driver metadata")
	flag.StringVar(&config.ClusterAddress, "clusteraddress", "localhost", "address of the springpath i/o dispatcher")
	flag.StringVar(&config.MountBase, "mountbase", "/run/springpath-docker-volumes", "base path for springpath volume mount points")
	flag.Parse()

	log.Println("starting docker volume plugin")

	volMap := new(VolumeMap)

	// http server.
	http.HandleFunc("/Plugin.Activate", pluginActivate)
	http.Handle("/VolumeDriver.Create", PluginHandler{Fn: volMap.Create, VMap: volMap})
	http.Handle("/VolumeDriver.Remove", PluginHandler{Fn: volMap.Create, VMap: volMap})
	http.Handle("/VolumeDriver.Mount", PluginHandler{Fn: volMap.Create, VMap: volMap})
	http.Handle("/VolumeDriver.Unmount", PluginHandler{Fn: volMap.Create, VMap: volMap})
	http.Handle("/VolumeDriver.Path", PluginHandler{Fn: volMap.Create, VMap: volMap})

	// XXX: set SO_REUSEADDR
	listener, err := net.Listen("unix", config.Sock)

	if err != nil {
		log.Fatal("failed to start http server", err)
	}

	log.Println("springpath volume driver listening on", config.Sock)
	err = http.Serve(listener, http.DefaultServeMux)
	if err != nil {
		log.Fatal("http serve failed", err)
	}
}
