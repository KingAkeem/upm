package registry

type Package interface {
	Type() string
	Name() string
	Author() string
	Description() string
	License() string
	Version() string
	//Created() time.Time
	//LastModified() time.Time
}
