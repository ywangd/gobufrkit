package tdcfio

import (
    "github.com/dgryski/go-bitstream"
    "fmt"
    "math"
    "io"
)

type bitWriter struct {
    w *bitstream.BitWriter
}

// BitWriter returns a pointer to bitWriter
//
// bitWriter implements tdcfio.Writer interface for writing to a binary stream
func BitWriter(writer io.Writer) *bitWriter {
    return &bitWriter{w: bitstream.NewWriter(writer)}
}

// WriteUint writes n bits least significant bits of given uint, most-significant-bit first
func (w *bitWriter) WriteUint(v uint, n int) error {
    return w.w.WriteBits(uint64(v), n)
}

// WriteInt writes a signed integer using n bits.
//
// The first bit is zero if value is positive or one if value is negative.
// The absolute value is written the same way as an uint.
func (w *bitWriter) WriteInt(v int, n int) error {
    if v < 0 {
        return w.WriteUint((1<<uint(n-1))|uint(-v), n)
    }
    return w.WriteUint(uint(v), n)
}

func (w *bitWriter) WriteBool(v bool) error {
    return w.w.WriteBit(v == true)
}

// Write a byte array, first n byte only
func (w *bitWriter) WriteBytes(v []byte, n int) error {
    var err error
    for i := 0; i < n; i++ {
        if i < len(v) {
            err = w.w.WriteByte(v[i])
        } else {
            err = w.w.WriteByte(0)
        }
        if err != nil {
            return err
        }
    }
    return nil
}

// Write Bin most-significant-bit first
func (w *bitWriter) WriteBinary(v *Binary, nbits int) error {
    if nbits != v.nbits {
        return fmt.Errorf("inconsistent number of bits for writing Bin data: %v != %v",
            nbits, v.nbits)
    }
    for i, b := range v.b {
        if i == len(v.b)-1 && v.nbits%NBITS_PER_BYTE != 0 {
            if err := w.w.WriteBits(uint64(b), nbits%NBITS_PER_BYTE); err != nil {
                return err
            }
        } else {
            if err := w.w.WriteByte(b); err != nil {
                return err
            }
        }
    }
    return nil
}

func (w *bitWriter) WriteFloat32(v float64) error {
    return w.WriteUint(uint(math.Float32bits(float32(v))), 32)
}
