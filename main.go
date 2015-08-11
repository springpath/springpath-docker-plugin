package main

import "net/http"
import "flag"
import "net"
import "log"
import "github.com/springpath/springpath-docker-plugin/driver"
import "github.com/springpath/springpath-docker-plugin/volume"
import "os"
import "path"

// Global State.
type Config struct {
	Sock           string
	ClusterAddress string
	NFSServer      string
	MountBase      string
	StateFile      string
}

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
	flag.StringVar(&config.ClusterAddress, "clusteraddress", "", "address of the springpath cluster master")
	flag.StringVar(&config.NFSServer, "nfsd", "localhost", "address of the springpath nfs server")
	flag.StringVar(&config.MountBase, "mountbase", "/run/springpath-docker-volumes", "base path for springpath volume mount points")
	flag.StringVar(&config.StateFile, "statefile", "/run/springpath-docker-volumes.state", "base path for springpath volume mount points")
	flag.Parse()

	if config.ClusterAddress == "" {
		flag.Usage()
		log.Fatal("clusteraddress must be set")
	}

	log.Println("starting docker volume plugin")

	var volmap, err = volume.New(config.ClusterAddress,
		config.NFSServer,
		config.MountBase,
		config.StateFile)
	if err != nil {
		log.Fatalf("Failed to connect to cluster backend")
	}

	driver.Register(http.DefaultServeMux, volmap)

	if err := os.MkdirAll(path.Dir(config.Sock), 0700); err != nil {
		log.Fatal("failed to create socket dir: ", err)
	}

	socketCleanup(config.Sock)

	listener, err := net.Listen("unix", config.Sock)
	if err != nil {
		log.Fatalf("failed to bind %s\n", err)
	}
	defer listener.Close()

	log.Println("springpath volume driver listening on", config.Sock)

	// XXX Clean shutdown using a goroutine for accepts and a
	// channel to handle close.
	// https://github.com/braintree/manners
	http.Serve(listener, http.DefaultServeMux)
	if err != nil {
		log.Fatalln("http serve failed", err)
	}
}
