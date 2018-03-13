package tdcfio_test

import (
    "testing"
    "fmt"
    "gobufrkit/tdcfio"
    assert2 "github.com/seanpont/assert"
)

func TestBinary_String(t *testing.T) {
    assert := assert2.Assert(t)

    var binary *tdcfio.Binary

    binary, _ = tdcfio.NewBinary([]byte{1}, 1)
    assert.Equal(fmt.Sprint(binary), "1")

    binary, _ = tdcfio.NewBinary([]byte{1}, 2)
    assert.Equal(binary.String(), "01")

    binary, _ = tdcfio.NewBinary([]byte{1}, 8)
    assert.Equal(binary.String(), "00000001")

    binary, _ = tdcfio.NewBinary([]byte{127, 1}, 10)
    assert.Equal(binary.String(), "0111111101")
}

func TestBinary_AtBit(t *testing.T) {
    assert := assert2.Assert(t)

    var binary *tdcfio.Binary

    binary, _ = tdcfio.NewBinary([]byte{85}, 8)
    expected := false
    for i := 0; i < 8; i++ {
        assert.Equal(binary.Bit(i), expected)
        expected = !expected
    }

    binary, _ = tdcfio.NewBinary([]byte{85, 85}, 16)
    expected = false
    for i := 8; i < 16; i++ {
        assert.Equal(binary.Bit(i), expected)
        expected = !expected
    }

    binary, _ = tdcfio.NewBinary([]byte{80, 2}, 10)
    assert.Equal(binary.Bit(8), true)
    assert.Equal(binary.Bit(9), false)
}

func TestNewBinary(t *testing.T) {
    assert := assert2.Assert(t)
    var err error

    _, err = tdcfio.NewBinary([]byte{1}, 1)
    assert.Nil(err)

    _, err = tdcfio.NewBinary([]byte{1}, 8)
    assert.Nil(err)

    _, err = tdcfio.NewBinary([]byte{1}, 9)
    assert.NotNil(err)

    _, err = tdcfio.NewBinary([]byte{1, 2}, 9)
    assert.Nil(err)

    _, err = tdcfio.NewBinary([]byte{1, 2}, 16)
    assert.Nil(err)

    _, err = tdcfio.NewBinary([]byte{1, 2}, 8)
    assert.NotNil(err)

    _, err = tdcfio.NewBinary([]byte{1, 2}, 17)
    assert.NotNil(err)
}

func TestBinary_UnmarshalJSON(t *testing.T) {
    assert := assert2.Assert(t)

    bin := &tdcfio.Binary{}

    bin.UnmarshalJSON([]byte("\"00110000\""))

    assert.Equal(bin.String(), "00110000")
}
