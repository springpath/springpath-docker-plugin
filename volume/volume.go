package volume

import "sync"
import "errors"
import "log"
import "os"

var ErrVolumeNotCreated = errors.New("Volume is not created")
var ErrVolumeNotFound = errors.New("Volume does not exist")
var ErrVolumeNotMounted = errors.New("Volume is not mounted")
var ErrVolumeNotRemoved = errors.New("Failed to Remove Volume")
var ErrVolumeNotUnmounted = errors.New("Failed to Unmount Volume")
var ErrVolumeNotReady = errors.New("Volume is not ready")
var ErrVolumeInUse = errors.New("Volume is in use")

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

	volumes     map[string]*Volume
	nfsServer   string
	mountBase   string
	routerHost  string
	initialized bool
	sync.Mutex
}

// The Volume plugin interface does not allow
// to pass a size in. Lets create 10G volumes
// for now.
const DefaultVolumeSize = 10 * 1024 * 1024 * 1024

func New(routerHost string, nfsServer string, mountBase string) (m *VolumeMap, err error) {
	m = &VolumeMap{
		volumes:    make(map[string]*Volume),
		mountBase:  mountBase,
		routerHost: routerHost,
		nfsServer:  nfsServer,
	}

	log.Printf("initializing volume driver with routerHost=%s and nfsServer=%s", m.routerHost, m.nfsServer)
	m.initialized = true
	return m, nil
}

func (m *VolumeMap) Create(name string) error {
	var v *Volume

	m.Lock()
	defer m.Unlock()

	v, ok := m.volumes[name]
	if ok && v.Created {
		// volume already exists and is
		// created.
		return nil
	} else {
		v = &Volume{
			Name:          name,
			DatastorePath: m.nfsUrl(name),
			MountedPath:   m.mountPoint(name),
			Size:          DefaultVolumeSize,
		}
	}

	m.volumes[v.Name] = v

	cmd := m.doCreate(v)
	if err := doCommand(cmd); err != nil {
		v.Created = false
		return ErrVolumeNotCreated
	}

	v.Created = true

	return nil
}

func (m *VolumeMap) Remove(name string) error {
	m.Lock()
	defer m.Unlock()

	v, ok := m.volumes[name]
	if !ok {
		return ErrVolumeNotFound
	}

	if v.Mounted {
		return ErrVolumeInUse
	}

	cmd := m.doRemove(v)
	if err := doCommand(cmd); err != nil {
		return ErrVolumeNotRemoved
	}

	delete(m.volumes, name)

	return nil
}

func (m *VolumeMap) Path(name string) (mountpoint string, err error) {
	m.Lock()
	defer m.Unlock()

	v, ok := m.volumes[name]
	if !ok {
		return "", ErrVolumeNotFound
	}

	if err := os.MkdirAll(v.MountedPath, 0700); err != nil {
		return "", ErrVolumeNotReady
	}

	return v.MountedPath, nil
}

func (m *VolumeMap) Mount(name string) (mountpoint string, err error) {
	m.Lock()
	defer m.Unlock()

	v, ok := m.volumes[name]

	if !ok {
		return "", ErrVolumeNotFound
	}

	if v.Mounted {
		return v.MountedPath, nil
	} else if !v.Created {
		return "", ErrVolumeNotCreated
	}

	cmd := m.doMount(v)
	if err := doCommand(cmd); err != nil {
		return "", ErrVolumeNotMounted
	}

	v.Mounted = true

	return v.MountedPath, nil
}

func (m *VolumeMap) Unmount(name string) error {
	m.Lock()
	defer m.Unlock()

	v, ok := m.volumes[name]

	if !ok {
		return ErrVolumeNotFound
	}

	if !v.Mounted {
		return ErrVolumeNotUnmounted
	}

	cmd := m.doUmount(v)
	if err := doCommand(cmd); err != nil {
		return ErrVolumeNotMounted
	}

	v.Mounted = false

	return nil
}
