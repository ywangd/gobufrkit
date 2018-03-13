package tdcfio_test

import (
    "testing"
    assert2 "github.com/seanpont/assert"
    "bytes"
    "gobufrkit/tdcfio"
    "os"
)

func TestBitWriter(t *testing.T) {
    assert := assert2.Assert(t)

    buffer := bytes.NewBufferString("")
    w := tdcfio.BitWriter(buffer)

    w.WriteBytes([]byte("BUFR"), 4)
    w.WriteUint(94, 24)

    w.WriteInt(4, 8)
    w.WriteUint(22, 24)
    w.WriteUint(0, 8)
    w.WriteUint(1, 16)
    w.WriteUint(0, 16)
    w.WriteUint(0, 8)

    w.WriteBool(false)

    binary, err := tdcfio.NewBinaryFromString("0000000")
    assert.Nil(err)
    w.WriteBinary(binary, 7)

    f, err := os.Open("../_testdata/contrived.bufr")
    defer f.Close()
    assert.Nil(err)

    b := make([]byte, 18)
    n, err := f.Read(b)
    assert.Nil(err)
    assert.Equal(n, 18)

    assert.Equal(b, buffer.Bytes())

}
