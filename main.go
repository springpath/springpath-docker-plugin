package main

import "net/http"
import "flag"
import "net"
import "log"
import "github.com/springpath/springpath-docker-plugin/driver"
import "github.com/springpath/springpath-docker-plugin/volume"
import "os"

// Global State.
type Config struct {
	Sock           string
	StateFile      string
	ClusterAddress string
	MountBase      string
}

// Clean up socket when we exit.
func socketCleanup(sock string) {
	_, err := os.Stat(sock)

	if err != nil {
		os.Remove(sock)
		log.Println("cleaning up socket", sock)
	}
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

	var volmap = new(volume.VolumeMap)

	driver.Register(http.DefaultServeMux, volmap)

	listener, err := net.Listen("unix", config.Sock)
	if err != nil {
		log.Fatalf("failed to bind %s\n", err)
	}
	defer listener.Close()
	defer socketCleanup(config.Sock)

	log.Println("springpath volume driver listening on", config.Sock)

	// XXX Clean shutdown using a goroutine for accepts and a
	// channel to handle close.
	// https://github.com/braintree/manners
	http.Serve(listener, http.DefaultServeMux)
	if err != nil {
		log.Fatalln("http serve failed", err)
	}
}
