package appender

type rotationFile struct {
	LogDir      string
	Backups     int
	LotateSize  int
	RotateDaily bool
}

func NewRotationFile() Appender {
	return &rotationFile{}
}

func init() {
	appenders["RotationFile"] = NewRotationFile
}
