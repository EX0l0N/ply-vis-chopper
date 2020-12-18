package cloudcompare

type PointCloud interface {
	Elements() int
	GetPointAt(int) interface{}
	GetPosition(interface{}) (int, bool)
}
