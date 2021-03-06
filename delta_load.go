package delta

import (
	"bytes"
	"encoding/binary"
)

// Load fills a new Delta structure from a byte
// array previously returned by Delta.Bytes().
func Load(data []byte, compressionAlgo string) (Delta, error) {
	//
	// uncompress the delta
	switch compressionAlgo {
	case zlibCompression:
		data = uncompressZlibBytes(data)
	case snappyCompression:
		data = uncompressSnappyBytes(data)
	case noCompression:
	default:
		data = uncompressZlibBytes(data)
	}
	buf := bytes.NewBuffer(data)
	readInt := func() int {
		var i int32
		err := binary.Read(buf, binary.BigEndian, &i)
		if err != nil {
			mod.Error("readInt() failed:", err)
			return -1
		}
		return int(i)
	}
	readBytes := func() []byte {
		var size int32
		err := binary.Read(buf, binary.BigEndian, &size)
		if err != nil {
			mod.Error("readBytes() failed @1:", err)
		}
		ar := make([]byte, size)
		var nread int
		nread, err = buf.Read(ar)
		if err != nil {
			mod.Error("readBytes() failed @2:", err)
		}
		if nread != int(size) {
			mod.Error("readBytes() failed @3: size:", size, "nread:", nread)
		}
		return ar
	}
	// read the header
	ret := Delta{
		sourceSize: readInt(),
		sourceHash: readBytes(),
		targetSize: readInt(),
		targetHash: readBytes(),
		newCount:   readInt(),
		oldCount:   readInt(),
	}
	// read the parts
	count := readInt()
	if count < 1 {
		return Delta{},
			mod.Error("readBytes() failed @4: invalid number of parts:", count)
	}
	ret.parts = make([]deltaPart, count)
	for i := range ret.parts {
		pt := &ret.parts[i]
		pt.sourceLoc = readInt()
		if pt.sourceLoc == -1 {
			pt.data = readBytes()
			pt.size = len(pt.data)
			continue
		}
		pt.size = readInt()
	}
	return ret, nil
} //                                                                        Load

// end
