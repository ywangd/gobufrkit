package unpack

import (
    "github.com/ywangd/gobufrkit/tdcfio"
    "github.com/ywangd/gobufrkit/table"
    "fmt"
    "github.com/ywangd/gobufrkit/bufr"
)

type UncompressBitUnpacker struct {
    r tdcfio.Reader
}

func (up *UncompressBitUnpacker) Unpack(info *bufr.PackingInfo) (interface{}, error) {
    switch info.Unit {
    case table.STRING:
        return up.r.ReadBytes(info.Nbits / tdcfio.NBITS_PER_BYTE)

    case table.CODE:
        return up.r.ReadInt(info.Nbits)

    case table.NONNEG_CODE, table.FLAG:
        if info.Nbits == 0 {
            return uint(0), nil
        }
        vuint, err := up.r.ReadUint(info.Nbits)
        if err != nil {
            return nil, err
        }
        return unpackMissing(info, vuint)

    case table.NUMERIC:
        vuint, err := up.r.ReadUint(info.Nbits)
        if err != nil {
            return nil, err
        }
        return unpackNumeric(info, vuint)

    case table.BINARY:
        return up.r.ReadBinary(info.Nbits)

    default:
        return nil, fmt.Errorf("unrecognised unit: %v", info.Unit)
    }
}

type JsonUnpacker struct {
    r tdcfio.Reader
}

func (up *JsonUnpacker) Unpack(info *bufr.PackingInfo) (interface{}, error) {
    switch info.Unit {
    case table.STRING:
        return up.r.ReadBytes(info.Nbits / tdcfio.NBITS_PER_BYTE)

    case table.CODE:
        return up.r.ReadInt(info.Nbits)

    case table.NONNEG_CODE, table.FLAG:
        return up.r.ReadUint(info.Nbits)

    case table.NUMERIC:
        return up.r.ReadNumber(info.Nbits)

    case table.BINARY:
        return up.r.ReadBinary(info.Nbits)

    default:
        return nil, fmt.Errorf("unrecognised unit: %v", info.Unit)
    }
}
