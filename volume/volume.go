package volume

import "sync"
import "errors"
import "log"

var ErrVolumeCreate = errors.New("Failed to Create Volume")
var ErrVolumeRemove = errors.New("Failed to Remove Volume")
var ErrVolumeMount = errors.New("Failed to Mount Volume")
var ErrVolumeUnmount = errors.New("Failed to Unmount Volume")
var ErrVolumeGet = errors.New("Failed to find specified Volume")

// Package Volume implements the volume management
// by calling child processes.
//
// XXX implement direct system calls/RPC calls, to mount etc.

// Individual mounted volumes.
type Volume struct {
	Name          string // docker volume name
	DatastorePath string // datastore NFS url
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

	volumes     map[string]Volume
	nfsServer   string
	mountBase   string
	routerHost  string
	initialized bool
	sync.Mutex
}

func New(routerHost string, nfsServer string, mountBase string) (m *VolumeMap, err error) {
	m = &VolumeMap{
		volumes:    make(map[string]Volume),
		mountBase:  mountBase,
		routerHost: routerHost,
		nfsServer:  nfsServer,
	}
	log.Printf("initializing volume driver with routerHost=%s and nfsServer=%s", m.routerHost, m.nfsServer)
	m.initialized = true
	return m, nil
}

func (m *VolumeMap) Create(name string) error {

	v := Volume{
		Name:          name,
		DatastorePath: m.nfsUrl(name),
		MountedPath:   m.mountPoint(name),
	}

	// Add the volume to our list.
	m.Lock()
	defer m.Unlock()
	m.volumes[v.Name] = v

	cmd := m.doCreate(&v)
	if err := doCommand(cmd); err != nil {
		v.Created = false
		return ErrVolumeCreate
	}

	v.Created = true

	return nil
}

func (m *VolumeMap) Remove(name string) error {
	m.Lock()
	defer m.Unlock()

	v, ok := m.volumes[name]
	if !ok {
		return ErrVolumeRemove
	}

	if v.Mounted {
		return ErrVolumeRemove
	}

	cmd := m.doRemove(&v)
	if err := doCommand(cmd); err != nil {
		return ErrVolumeRemove
	}

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

	if v.Mounted {
		return v.MountedPath, nil
	} else if !v.Created {
		return "", ErrVolumeMount
	}

	cmd := m.doMount(&v)
	if err := doCommand(cmd); err != nil {
		return "", nil
	}

	v.Mounted = true

	return v.MountedPath, nil
}

func (m *VolumeMap) Unmount(name string) error {
	m.Lock()
	defer m.Unlock()

	v, ok := m.volumes[name]

	if !ok {
		return ErrVolumeGet
	}

	if !v.Mounted {
		return ErrVolumeUnmount
	}

	cmd := m.doUmount(&v)
	if err := doCommand(cmd); err != nil {
		return ErrVolumeUnmount
	}

	v.Mounted = false

	return nil
}
