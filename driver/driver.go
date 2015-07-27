package driver

import "net/http"
import "io/ioutil"
import "encoding/json"
import "github.com/springpath/springpath-docker-plugin/volume"
import "log"

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

type nameFn func(string) error
type pathFn func(string) (string, error)

type pluginHandler struct {
	http.Handler
	VMap *volume.VolumeMap
	Fn   interface{} // one of nameFn or pathFn. sigh.
}

type Message struct {
	Name       string `json:"omitempty"`
	Err        error  `json:"omitempty"`
	Mountpoint string `json:"omitempty"`
}

// Common portions of all the plugin endpoints.
func (h pluginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var requestPath = r.URL.Path
	var request, response Message

	defer r.Body.Close()
	var requestBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("req %p: url %s, body %s\n", r, requestPath, requestBody)

	json.Unmarshal(requestBody, &request)

	switch fn := h.Fn.(type) {
	default:
		log.Fatalf("Unknown type %T", fn)
	case nameFn:
		response.Err = fn(request.Name)
	case pathFn:
		response.Mountpoint, response.Err = fn(request.Name)
	}

	resp, err := json.Marshal(response)

	w.Write(resp)
	return
}

func Register(mux *http.ServeMux, volMap *volume.VolumeMap) {
	// http server.
	mux.HandleFunc("/Plugin.Activate", pluginActivate)
	mux.Handle("/VolumeDriver.Create", pluginHandler{Fn: volMap.Create, VMap: volMap})
	mux.Handle("/VolumeDriver.Remove", pluginHandler{Fn: volMap.Remove, VMap: volMap})
	mux.Handle("/VolumeDriver.Mount", pluginHandler{Fn: volMap.Mount, VMap: volMap})
	mux.Handle("/VolumeDriver.Unmount", pluginHandler{Fn: volMap.Unmount, VMap: volMap})
	mux.Handle("/VolumeDriver.Path", pluginHandler{Fn: volMap.Path, VMap: volMap})
}
