package driver

import "net/http/httptest"
import "testing"

func TestAPI(t *testing.T) {
	s := new(httptest.Server)
	s.Close()
	return
}
