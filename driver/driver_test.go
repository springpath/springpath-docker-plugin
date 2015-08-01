package driver

import "net/http/httptest"
import "net/http"
import "testing"
import "github.com/springpath/springpath-docker-plugin/volume"
import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"

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

func do(baseurl string, op string, name string) (mp string, err error) {
	var request, response Message
	request.Name = name
	var reqbody, respbody []byte
	var resp = new(http.Response)

	if reqbody, err = json.Marshal(request); err != nil {
		return
	}

	url := baseurl + "VolumeDriver." + op
	if resp, err = http.Post(url, dockerVersionMimeType, bytes.NewReader(reqbody)); err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return mp, errors.New("HttpError" + resp.Status)
	}

	if respbody, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if err = json.Unmarshal(respbody, &response); err != nil {
		return
	}

	return response.Mountpoint, response.Err
}

func activate(baseurl string) error {
	url := baseurl + "Plugin.Activate"
	resp, err := http.Post(url, dockerVersionMimeType, nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("HttpError" + resp.Status)
	}

	return nil
}

func TestAPI(t *testing.T) {
	mux := http.NewServeMux()
	vmap := DummyVolMap{}
	Register(mux, &vmap)
	server := httptest.NewServer(mux)

	ops := []string{"Create", "Remove", "Mount", "Unmount", "Path"}

	if err := activate(server.URL + "/"); err != nil {
		t.Fatalf("activation failed", err)
	}

	// write tests.
	for _, op := range ops {
		if _, err := do(server.URL+"/", op, "test"); err != nil {
			t.Fatalf(op, err)
		}
	}

	server.Close()
	return
}
