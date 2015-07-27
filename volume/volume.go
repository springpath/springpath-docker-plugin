package volume

// Individual mounted volumes.
type Volume struct {
	Name          string // docker volume name
	DatastorePath string // springpath datastore name + path
	Size          uint64 // size of the datastore in bytes.
	Mounted       bool   // locally mounted.
	Alive         bool   // scvmclient is reachable.
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
	m map[string]Volume
}

func (v *VolumeMap) Create(name string) error {
	return nil
}

func (v *VolumeMap) Remove(name string) error {
	return nil
}

func (v *VolumeMap) Path(name string) (mountpoint string, err error) {
	return
}

func (v *VolumeMap) Mount(name string) (mountpoint string, err error) {
	return
}

func (v *VolumeMap) Unmount(name string) error {
	return nil
}
