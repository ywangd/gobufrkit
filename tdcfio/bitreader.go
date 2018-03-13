package tdcfio

import (
    "io"
    "github.com/dgryski/go-bitstream"
    "math"
    "bufio"
    "bytes"
    "fmt"
    "github.com/pkg/errors"
)

type BitReader struct {
    r *bitstream.BitReader

    // Bit position of current read
    pos int
}

// NewBitReader returns a pointer to BitReader.
//
// BitReader implements the tdcfio.Reader interface for reading from a binary stream.
func NewBitReader(r io.Reader) *BitReader {
    return &BitReader{
        r:   bitstream.NewReader(r),
    }
}

func (r *BitReader) Pos() int {
    return r.pos
}

func (r *BitReader) ReadNumber(nbits int) (float64, error) {
    x, err := r.ReadUint(nbits)
    if err != nil {
        return 0, err
    }
    return float64(x), nil
}

func (r *BitReader) ReadUint(nbits int) (uint, error) {
    r.pos += nbits
    value, err := r.r.ReadBits(nbits)
    return uint(value), err
}

func (r *BitReader) ReadInt(nbits int) (int, error) {
    negative, err := r.ReadBool()
    if err != nil {
        return 0, err
    }
    r.pos += nbits - 1
    value, err := r.r.ReadBits(nbits - 1)

    if negative {
        return -int(value), err
    } else {
        return int(value), err
    }
}

func (r *BitReader) ReadBool() (bool, error) {
    r.pos += 1
    value, err := r.r.ReadBit()
    return bool(value), err
}

func (r *BitReader) ReadBytes(nbytes int) ([]byte, error) {
    r.pos += nbytes * NBITS_PER_BYTE
    bs := make([]byte, nbytes)
    for i := 0; i < nbytes; i++ {
        b, err := r.r.ReadByte()
        if err != nil {
            return bs, err
        }
        bs[i] = b
    }
    return bs, nil
}

func (r *BitReader) ReadBinary(nbits int) (*Binary, error) {
    r.pos += nbits

    var (
        currentValue uint64
        buffer       []byte
        err          error
    )

    // TODO: possible optimisation?
    nbitsSaved := nbits
    for nbits > 0 {
        if nbits < NBITS_PER_BYTE {
            currentValue, err = r.r.ReadBits(nbits)
        } else {
            currentValue, err = r.r.ReadBits(NBITS_PER_BYTE)
        }
        if err != nil {
            return nil, nil
        }
        buffer = append(buffer, byte(currentValue))
        nbits -= NBITS_PER_BYTE
    }
    return &Binary{b: buffer, nbits: nbitsSaved}, nil
}

func (r *BitReader) ReadFloat32() (float64, error) {
    v, err := r.ReadUint(32)
    if err != nil {
        return 0, err
    }
    n := uint32(v)
    f := math.Float32frombits(n)
    return float64(f), nil
}

// PeekableBitReader implements PeekableReader for reading from binary stream.
type PeekableBitReader struct {
    *BitReader
    br *bufio.Reader
}

// NewPeekableBitReader returns a pointer to a PeekableBitReader.
//
// The implementation of PeekableBitReader employs a BitReader for normal read
// operations. It realises the peeking functions by utilising a bufio.Reader object.
// The BitReader and
// bufio.Reader share a single underlying io.Reader object. When a peeking
// function is called, bufio.Reader peeks enough bits from the underlying
// io.Reader. A new BitReader is created for the peeked bits and it is then
// used to get the required value.
// Peek functions always operate in integer number of bytes, i.e. peek must
// both begin and end at the boundary of a byte.
func NewPeekableBitReader(r io.Reader) *PeekableBitReader {
    br := bufio.NewReader(r)
    return &PeekableBitReader{NewBitReader(br), br}
}

// PeekUint returns an uint value of n bits by skipping number of skip bytes
func (r *PeekableBitReader) PeekUint(skip int, n int) (uint, error) {
    nbytes := n / NBITS_PER_BYTE
    if n%NBITS_PER_BYTE > 0 {
        nbytes += 1
    }
    b, err := r.peek(skip, nbytes)
    if err != nil {
        return 0, errors.Wrap(err, "cannot peek bytes")
    }
    v, err := NewBitReader(bytes.NewReader(b)).ReadUint(n)
    if err != nil {
        return 0, errors.Wrap(err, "cannot read peeked bytes")
    }
    return v, nil
}

func (r *PeekableBitReader) PeekBytes(skip int, nbytes int) ([]byte, error) {
    b, err := r.peek(skip, nbytes)
    if err != nil {
        return []byte{}, errors.Wrap(err, "cannot peek bytes")
    }
    v, err := NewBitReader(bytes.NewReader(b)).ReadBytes(nbytes)
    if err != nil {
        return []byte{}, errors.Wrap(err, "cannot read peeked bytes")
    }
    return v, nil
}

func (r *PeekableBitReader) peek(skip int, n int) ([]byte, error) {
    if r.pos%NBITS_PER_BYTE != 0 {
        return nil, fmt.Errorf("can only peek at complete byte boundary: (%v, %v)",
            r.pos, n)
    }
    b, err := r.br.Peek(skip + n)
    if err != nil {
        if len(b) > skip {
            return b[skip:], err
        }
        return []byte{}, err
    }
    return b[skip:], nil
}
