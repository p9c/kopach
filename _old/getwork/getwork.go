package getwork

import (
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/parallelcointeam/pod/chaincfg/chainhash"
	"github.com/parallelcointeam/pod/wire"
)

// ToBlockHeader converts a getwork response data string into a BlockHeader
func ToBlockHeader(data string) (bh wire.BlockHeader) {
	bd, _ := hex.DecodeString(data)
	bh.Version = int32(binary.BigEndian.Uint32(bd[:4]))
	bh.PrevBlock.SetBytes(reverseByUint32(bd[4:36]))
	bh.MerkleRoot.SetBytes(reverseByUint32(bd[36:68]))
	bh.Timestamp = time.Unix(int64(binary.BigEndian.Uint32(bd[68:72])), 0)
	bh.Bits = binary.BigEndian.Uint32(bd[72:76])
	bh.Nonce = binary.BigEndian.Uint32(bd[76:80])
	return
}

func reverseByUint32(data []byte) (out []byte) {
	if len(data) != 32 {
		return
	}
	for i := 32; i > 0; i -= 4 {
		out = append(out, data[i-4:i]...)
	}
	out = reverse(out)
	return
}

// BytesToInt32 converts a byte slice to an int32, padding or truncating as required
func BytesToInt32(data []byte) (out int32) {
	if len(data) > 4 {
		data = data[:4]
	}
	if len(data) < 4 {
		d := make([]byte, 4)
		for i := range data {
			d[i] = data[i]
		}
		data = d
	}
	out += int32(data[0])
	out += int32(data[1]) << 8
	out += int32(data[1]) << 16
	out += int32(data[1]) << 24
	return
}

// BytesToUint32 converts a byte slice to an int32, padding or truncating as required
func BytesToUint32(data []byte) (out uint32) {
	if len(data) > 4 {
		data = data[:4]
	}
	if len(data) < 4 {
		d := make([]byte, 4)
		for i := range data {
			d[i] = data[i]
		}
		data = d
	}
	out += uint32(data[0])
	out += uint32(data[1]) << 8
	out += uint32(data[1]) << 16
	out += uint32(data[1]) << 24
	return
}

// BytesToChainhash converts 32 element byte slice to a chainhash.Hash
func BytesToChainhash(data []byte) (out chainhash.Hash) {
	if len(data) != 64 {
		return
	}
	d := make([]byte, 64)
	for i := range data {
		d[64-i] = data[i]
	}
	o, _ := chainhash.NewHash(data)
	return *o
}

func reverse(b []byte) (o []byte) {
	o = make([]byte, len(b))
	for i := range b {
		o[i] = b[len(b)-1-i]
	}
	return
}

/*
000000055fb30cdfe17572697da2668ffa12bbc822ed4d85f91b0c2cecbaa35e00000e41f85bee81614ec9b55012f7fe5f793f9340f7d4304164c3590e7352f212bd4a6f5c1a8b781e0fffff00000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000080020000

00000e41ecbaa35ef91b0c2c22ed4d85fa12bbc87da2668fe17572695fb30cdf
5fb30cdfe17572697da2668ffa12bbc822ed4d85f91b0c2cecbaa35e00000e41



00000e41 ecbaa35ef91b0c2c22ed4d85fa12bbc87da2668fe1757269 5fb30cdf
5fb30cdf e17572697da2668ffa12bbc822ed4d85f91b0c2cecbaa35e 00000e41

*/
