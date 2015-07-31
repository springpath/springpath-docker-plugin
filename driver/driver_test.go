package driver

import "net/http/httptest"
import "net/http"
import "testing"
import "github.com/springpath/springpath-docker-plugin/volume"

type DummyVolMap struct {
	volume.VolumeDriver

	retErr error // return this error for all operations.
}

func (m *DummyVolMap) Create(name string) error {
	return m.retErr
}

func (m *DummyVolMap) Remove(name string) error {
	return m.retErr
}

func (m *DummyVolMap) Mount(name string) (mp string, err error) {
	return name, m.retErr
}

func (m *DummyVolMap) Unmount(name string) error {
	return m.retErr
}

func (m *DummyVolMap) Path(name string) (mp string, err error) {
	return name, m.retErr
}

func callNameOp(url string, op string, name string) (err error) {
}

func callPathOp(url string, op string, name string) (mp string, err error) {
	return
}

func callActivate() error {

}

func TestAPI(t *testing.T) {
	mux := http.NewServeMux()
	vmap := DummyVolMap{}
	Register(mux, &vmap)
	server := httptest.NewServer(mux)

	// write tests.
	callNameOp(server.URL, "Create", "test")
	callNameOp(server.URL, "Remove", "test")
	callPathOp(server.URL, "Mount", "test")
	callNameOp(server.URL, "Unmount", "test")
	callPathOp(server.URL, "Path", "test")

	server.Close()
	return
}
