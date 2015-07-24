package main

import "net/http"
import "net"
import "log"
import "flag"
import "io/ioutil"

type Config struct {
	Sock           string
	StateFile      string
	ClusterAddress string
	MountBase      string
}

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
	Unmount(name string) (mountpoint string, err error)
}

type VolumeList map[string]Volume

type DockerVolumes struct {
	VolumeDriver
	volumes VolumeList
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

func pluginHandler(w http.ResponseWriter, r *http.Request) {
	var requestPath = r.URL.Path

	defer r.Body.Close()
	var requestBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("req %p: url %s, body %s\n", r, requestPath, requestBody)

	// Until we write all handlers.
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	v := new(Config)

	// parse command line.
	flag.StringVar(&v.Sock, "sockpath", "/run/docker/plugins/springpath.sock", "unix domain socket docker talks to")
	flag.StringVar(&v.StateFile, "statefile", "/SYSTEM/volume-driver.json", "springpath volume driver metadata")
	flag.StringVar(&v.ClusterAddress, "clusteraddress", "localhost", "address of the springpath i/o dispatcher")
	flag.StringVar(&v.MountBase, "mountbase", "/run/springpath-docker-volumes", "base path for springpath volume mount points")
	flag.Parse()

	log.Println("starting docker volume plugin")

	// http server.
	http.HandleFunc("/Plugin.Activate", pluginActivate)
	http.HandleFunc("/", pluginHandler)

	listener, err := net.Listen("unix", v.Sock)

	if err != nil {
		log.Fatal("failed to start http server", err)
	}

	log.Println("springpath volume driver listening on", v.Sock)
	err = http.Serve(listener, http.DefaultServeMux)
	if err != nil {
		log.Fatal("http serve failed", err)
	}
}
