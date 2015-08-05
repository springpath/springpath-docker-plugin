package volume

import "sync"

// Package Volume implements the volume management
// by calling child processes.
//
// XXX implement direct system calls/RPC calls, to mount etc.

// Individual mounted volumes.
type Volume struct {
	Name          string // docker volume name
	DatastorePath string // springpath datastore name + path
	MountedPath   string // path where datastore is currently mounted
	Size          uint64 // size of the datastore in bytes.
	Mounted       bool   // locally mounted.
	Alive         bool   // scvmclient is reachable.
	Created       bool   // volume exists.
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

	volumes    map[string]Volume
	nfsServer  string
	mountBase  string
	routerHost string
	sync.Mutex
}

func New(routerHost string, nfsServer string, mountBase string) (m *VolumeMap, err error) {
	m = &VolumeMap{
		volumes:    make(map[string]Volume),
		mountBase:  mountBase,
		routerHost: routerHost,
		nfsServer:  nfsServer,
	}
	return m, nil
}

func (m *VolumeMap) Create(name string) error {

	v := Volume{Name: name}

	// all datastores are at the root.
	v.DatastorePath = "/" + v.Name

	// Add the volume to our list.
	m.Lock()
	defer m.Unlock()
	m.volumes[v.Name] = v

	return nil
}

func (m *VolumeMap) Unmount(name string) error {
	return nil
}

func (m *VolumeMap) Remove(name string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.volumes, name)

	return nil
}

func (m *VolumeMap) Path(name string) (mountpoint string, err error) {
	m.Lock()
	defer m.Unlock()

	v := m.volumes[name]

	return v.MountedPath, nil
}

func (m *VolumeMap) Mount(name string) (mountpoint string, err error) {
	m.Lock()
	defer m.Unlock()

	v := m.volumes[name]

	return v.MountedPath, nil
}
