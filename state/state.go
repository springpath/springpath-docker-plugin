// Package state persists information data volumes, containers, and their
// relationships.  This could in turn be backed by systems zookeeper, etcd, or
// a shared filesystem.
package state

import "net/url"
import "sync"

type MountInfo struct {
	ContainerId string
	Volumename  string
	Mountpath   string
	Hostname    string
}

type StateFile struct {
	generation uint64 // In memory version to track updates.
	dbVersion  string // database format version.
	dirty      bool
	fpath      url.URL
	hostname   string
	sync.Mutex
}

func New(stateRoot url.URL) (s *StateFile, err error) {
	s = &StateFile{
		dirty:     true,
		dbVersion: "", // unknown.
		fpath:     stateRoot,
	}
	return s, nil
}

func (s *StateFile) Sync() error {
	return nil
}

func (s *StateFile) AddMountInfo(m MountInfo) error {
	return nil
}

func (s *StateFile) RemoveMountInfo(mountpath string) error {
	return nil
}

func (s *StateFile) GetMountInfoByVolume(name string) (info []MountInfo, err error) {
	return
}

func (s *StateFile) GetMountInfoByHost(host string) (info []MountInfo, err error) {
	return
}
