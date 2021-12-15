package archive

import (
	"encoding/binary"
)

type FileMeta struct {
	Name string
	Size int64
}

func (f *FileMeta) Serialize() []byte {
	data := []byte{}
	// Store name
	data = append(data, []byte(f.Name)...)
	data = append(data, byte(0))
	// Store file size
	b := make([]byte, 8)
	binary.PutVarint(b, f.Size)
	// Store data
	data = append(data, b...)
	return data
}

func readString(data []byte, n int) (string, int) {
	name := ""
	var i int
	for i = n; i < len(data); i++ {
		if data[i] == 0 {
			name = string(data[n:i])
			break
		}
	}
	return name, i
}

func (f *FileMeta) Deserialize(data []byte) {
	var n int
	f.Name, n = readString(data, n)
	f.Size = int64(binary.LittleEndian.Uint64(data[n : n+8]))
	n += 8
}
