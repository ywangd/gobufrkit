package unpack

import (
    "github.com/ywangd/gobufrkit/tdcfio"
    "github.com/ywangd/gobufrkit/table"
    "fmt"
    "github.com/ywangd/gobufrkit/bufr"
)

const NBITS_FOR_NBITS_DIFF = 6

type CompressedVal struct {
    MinValue  interface{}
    NbitsDiff int
    Values    []interface{}
}

type CompressedBitUnpacker struct {
    r        tdcfio.Reader
    nsubsets int
}

func (up *CompressedBitUnpacker) Unpack(info *bufr.PackingInfo) (interface{}, error) {
    switch info.Unit {
    case table.STRING:
        return up.unpackString(info)

    case table.CODE:
        return up.unpackCode(info)

    case table.NONNEG_CODE, table.FLAG:
        // TODO: handle compatible 237000 etc
        if info.Nbits == 0 {
            values := make([]interface{}, up.nsubsets)
            for i := 0; i < up.nsubsets; i++ {
                values[i] = uint(0)
            }
            return &CompressedVal{
                MinValue:  uint(0),
                NbitsDiff: 0,
                Values:    values,
            }, nil
        }
        return up.unpackOthers(info, unpackMissing)

    case table.NUMERIC:
        return up.unpackOthers(info, unpackNumeric)

    case table.BINARY:
        return up.unpackOthers(info, unpackBinary)

    default:
        return nil, fmt.Errorf("unrecognised unit: %v", info.Unit)
    }
}

func (up *CompressedBitUnpacker) unpackString(info *bufr.PackingInfo) (interface{}, error) {

    var (
        ret CompressedVal
        err error
    )

    nbytes := info.Nbits / tdcfio.NBITS_PER_BYTE
    ret.MinValue, err = up.r.ReadBytes(nbytes)
    if err != nil {
        return nil, err
    }
    // TODO: Is this actually nbytesDiff for string?
    ret.NbitsDiff, err = readNbitsDiff(up.r)
    if err != nil {
        return nil, err
    }

    ret.Values = make([]interface{}, up.nsubsets)

    if ret.NbitsDiff == 0 {
        for i := 0; i < up.nsubsets; i++ {
            ret.Values[i] = ret.MinValue
        }
    } else {
        if string(ret.MinValue.([]byte)) != string(make([]byte, nbytes)) {
            return nil, fmt.Errorf(
                "different string must be compressed with empty minimum value")
        }
        for i := 0; i < up.nsubsets; i++ {
            if ret.Values[i], err = up.r.ReadBytes(ret.NbitsDiff); err != nil {
                return nil, err
            }
        }
    }

    return &ret, nil
}

func (up *CompressedBitUnpacker) unpackOthers(info *bufr.PackingInfo,
    unpackFunc func(*bufr.PackingInfo, uint) (interface{}, error)) (interface{}, error) {

    var (
        ret CompressedVal
        err error
    )

    ret.MinValue, err = up.r.ReadUint(info.Nbits)
    if err != nil {
        return nil, err
    }
    ret.NbitsDiff, err = readNbitsDiff(up.r)
    if err != nil {
        return nil, err
    }

    ret.Values = make([]interface{}, up.nsubsets)

    if ret.NbitsDiff == 0 {
        v, err := unpackFunc(info, ret.MinValue.(uint))
        if err != nil {
            return nil, err
        }
        for i := 0; i < up.nsubsets; i++ {
            ret.Values[i] = v
        }
    } else {
        if bufr.IsMissing(ret.MinValue, info.Nbits) {
            return nil, fmt.Errorf("missing value must have zero valued nbitsDiff")
        }
        for i := 0; i < up.nsubsets; i++ {
            diff, err := up.r.ReadUint(ret.NbitsDiff)
            if err != nil {
                return nil, err
            }
            if bufr.IsMissing(diff, ret.NbitsDiff) {
                ret.Values[i] = nil
            } else {
                if ret.Values[i], err = unpackFunc(info, ret.MinValue.(uint)+diff); err != nil {
                    return nil, err
                }
            }
        }
    }
    return &ret, nil
}

func (up *CompressedBitUnpacker) unpackCode(info *bufr.PackingInfo) (interface{}, error) {
    var (
        ret CompressedVal
        err error
    )
    ret.MinValue, err = up.r.ReadInt(info.Nbits)
    if err != nil {
        return nil, err
    }
    ret.NbitsDiff, err = readNbitsDiff(up.r)
    if err != nil {
        return nil, err
    }

    ret.Values = make([]interface{}, up.nsubsets)

    if ret.NbitsDiff == 0 {
        for i := 0; i < up.nsubsets; i++ {
            ret.Values[i] = ret.MinValue
        }
    } else {
        for i := 0; i < up.nsubsets; i++ {
            diff, err := up.r.ReadUint(ret.NbitsDiff)
            if err != nil {
                return nil, err
            }
            if bufr.IsMissing(diff, ret.NbitsDiff) {
                ret.Values[i] = nil
            } else {
                // TODO: in theory the uint diff could be out of the int range
                ret.Values[i] = ret.Values[i].(int) + int(diff)
            }
        }
    }
    return &ret, nil
}

func readNbitsDiff(reader tdcfio.Reader) (int, error) {
    NbitsDiff, err := reader.ReadUint(NBITS_FOR_NBITS_DIFF)
    if err != nil {
        return 0, err
    }
    return int(NbitsDiff), nil
}
