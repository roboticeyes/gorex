package rex

// File represents a complete valid REX file which can
// either be stored locally or sent to an arbirary writer with
// the Encoder.
type File struct {
	Header    Header
	PointList []PointList
}
