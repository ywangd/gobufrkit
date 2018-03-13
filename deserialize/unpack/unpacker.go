package unpack

import (
    "github.com/ywangd/gobufrkit/tdcfio"
    "math"
    "fmt"
    "github.com/ywangd/gobufrkit/bufr"
)

type Unpacker interface {
    Unpack(info *bufr.PackingInfo) (interface{}, error)
}

// NewUnpacker returns an appropriate Unpacker implementation based on given reader type
func NewUnpacker(reader tdcfio.Reader, nsubsets int, compressed bool, inputType tdcfio.InputType) (Unpacker, error) {
    switch inputType {
    case tdcfio.BinaryInput:
        if compressed {
            return &CompressedBitUnpacker{r: reader, nsubsets: nsubsets}, nil
        }
        return &UncompressBitUnpacker{r: reader}, nil
    case tdcfio.FlatJsonInput:
        return &JsonUnpacker{r: reader}, nil
    default:
        return nil, fmt.Errorf("no candidate unpacker for %v", inputType)
    }
}

func unpackMissing(info *bufr.PackingInfo, x uint) (interface{}, error) {
    if bufr.IsMissing(x, info.Nbits) {
        return nil, nil
    }
    return x, nil
}

// Numeric values are always unpacked to float64
func unpackNumeric(info *bufr.PackingInfo, x uint) (interface{}, error) {
    // No need to apply any packing for missing value
    if y, _ := unpackMissing(info, x); y == nil {
        return nil, nil
    }
    xfloat := float64(x)

    if info.Refval != 0 || info.Scale != 0 {
        if info.Refval != 0 {
            xfloat += info.Refval
        }
        if info.Scale != 0 {
            xfloat /= math.Pow10(info.Scale)
        }
        return xfloat, nil
    }

    return xfloat, nil
}

func unpackBinary(info *bufr.PackingInfo, x uint) (interface{}, error) {
    return tdcfio.NewBinaryFromUint(x, info.Nbits)
}
