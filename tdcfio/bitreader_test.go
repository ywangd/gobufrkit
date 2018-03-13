package tdcfio_test

import (
    "testing"
    assert2 "github.com/seanpont/assert"
    "gobufrkit/tdcfio"
    "os"
)

func TestBitReader(t *testing.T) {
    assert := assert2.Assert(t)

    var (
        err error
        b []byte
        u uint
        i int
        q bool
        s *tdcfio.Binary
    )

    f, err := os.Open("../_testdata/contrived.bufr")
    defer f.Close()
    assert.Nil(err)
    r := tdcfio.NewBitReader(f)
    b, err = r.ReadBytes(4)
    assert.Nil(err)
    assert.Equal(string(b), "BUFR")

    u, err = r.ReadUint(24)
    assert.Nil(err)
    assert.Equal(u, uint(94))

    i, err = r.ReadInt(8)
    assert.Nil(err)
    assert.Equal(i, 4)

    r.ReadUint(24)
    r.ReadUint(8)
    r.ReadUint(16)
    r.ReadUint(16)
    r.ReadUint(8)

    q, err = r.ReadBool()
    assert.Nil(err)
    assert.False(q, "value should be false")

    s, err = r.ReadBinary(7)
    assert.Nil(err)
    assert.Equal(s.String(), "0000000")
}

func TestPeekableBitReader(t *testing.T) {
    assert := assert2.Assert(t)

    var (
        err error
        b []byte
        u uint
        i int
        q bool
        s *tdcfio.Binary
    )

    f, err := os.Open("../_testdata/contrived.bufr")
    defer f.Close()
    assert.Nil(err)
    r := tdcfio.NewPeekableBitReader(f)

    b, err = r.PeekBytes(0, 4)
    assert.Nil(err)
    assert.Equal(string(b), "BUFR")

    u, err = r.PeekUint(8, 24)
    assert.Nil(err)
    assert.Equal(u, uint(22))

    b, err = r.ReadBytes(4)
    assert.Nil(err)
    assert.Equal(string(b), "BUFR")

    u, err = r.ReadUint(24)
    assert.Nil(err)
    assert.Equal(u, uint(94))

    i, err = r.ReadInt(8)
    assert.Nil(err)
    assert.Equal(i, 4)

    r.ReadUint(24)
    r.ReadUint(8)
    r.ReadUint(16)
    r.ReadUint(16)
    r.ReadUint(8)

    u, err = r.PeekUint(4, 8)
    assert.Nil(err)
    assert.Equal(u, uint(18))

    q, err = r.ReadBool()
    assert.Nil(err)
    assert.False(q, "value should be false")

    s, err = r.ReadBinary(7)
    assert.Nil(err)
    assert.Equal(s.String(), "0000000")

    u, err = r.PeekUint(7, 8)
    assert.Nil(err)
    assert.Equal(u, uint(2))

}