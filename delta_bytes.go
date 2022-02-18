package delta

import (
	"bytes"
	"encoding/binary"
)

// Bytes converts the Delta structure to a byte array
// (for serializing to a file, etc.)
func (ob *Delta) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	//
	writeInt := func(i int) error {
		err := binary.Write(buf, binary.BigEndian, int32(i))
		if err != nil {
			return mod.Error("writeInt(", i, ") failed:", err)
		}
		return nil
	}
	writeBytes := func(data []byte) error {
		err := writeInt(len(data))
		if err != nil {
			return mod.Error("writeBytes([", len(data), "]) failed @1:", err)
		}
		var n int
		n, err = buf.Write(data)
		if err != nil {
			return mod.Error("writeBytes([", len(data), "]) failed @2:", err)
		}
		if n != len(data) {
			return mod.Error("writeBytes([", len(data), "]) failed @3:",
				"wrote wrong number of bytes:", n)
		}
		return nil
	}
	// write the header
	writeInt(ob.sourceSize)
	writeBytes(ob.sourceHash)
	writeInt(ob.targetSize)
	writeBytes(ob.targetHash)
	writeInt(ob.newCount)
	writeInt(ob.oldCount)
	writeInt(len(ob.parts))
	//
	// write the parts
	for _, part := range ob.parts {
		writeInt(part.sourceLoc)
		if part.sourceLoc == -1 {
			writeBytes(part.data)
			continue
		}
		writeInt(part.size)
	}
	// compress the delta

	// based on compression algo we need to perform the compression
	switch ob.compressionAlgo {
	case zlibCompression:
		return compressZlibBytes(buf.Bytes())
	case snappyCompression:
		return  compressSnappyBytes(buf.Bytes())
	case noCompression:
		return  buf.Bytes()
	default:
		return  compressZlibBytes(buf.Bytes())
	}
} //                                                                       Bytes

// end
