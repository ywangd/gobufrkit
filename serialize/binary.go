package serialize

import (
    "io"
    "github.com/ywangd/gobufrkit/tdcfio"
    "github.com/ywangd/gobufrkit/bufr"
    "github.com/ywangd/gobufrkit/table"
    "fmt"
    "github.com/ywangd/gobufrkit/serialize/pack"
)

type BinaryVisitor struct {
    w      tdcfio.Writer
    packer pack.Packer
}

func NewBinaryVisitor(writer io.Writer) *BinaryVisitor {
    return &BinaryVisitor{
        w: tdcfio.BitWriter(writer),
    }
}

func (v *BinaryVisitor) VisitMessage(message *bufr.Message) error {
    for _, section := range message.Sections() {
        if err := section.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *BinaryVisitor) VisitSection(section *bufr.Section) error {
    for _, field := range section.Fields() {
        if err := field.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *BinaryVisitor) VisitField(field *bufr.Field) error {
    var err error
    switch value := field.Value.(type) {
    case uint:
        err = v.w.WriteUint(value, field.Nbits)
    case int:
        err = v.w.WriteInt(value, field.Nbits)
    case float64:
        err = v.w.WriteFloat32(value)
    case bool:
        err = v.w.WriteBool(value)
    case []byte:
        err = v.w.WriteBytes(value, field.Nbits/tdcfio.NBITS_PER_BYTE)
    case *tdcfio.Binary:
        err = v.w.WriteBinary(value, field.Nbits)
    case *table.UnexpandedTemplate:
        for _, id := range value.Ids() {
            v.w.WriteUint(uint(id.F()), value.Fbits())
            v.w.WriteUint(uint(id.X()), value.Xbits())
            v.w.WriteUint(uint(id.Y()), value.Ybits())
        }
    case *bufr.Payload:
        value.Accept(v)
    default:
        err = fmt.Errorf("unsupported value type: %T", value)
    }
    return err
}

func (v *BinaryVisitor) VisitPayload(payload *bufr.Payload) error {
    if payload.Compressed {
        v.packer = pack.NewCompressedPacker(v.w)
        nsubsets := len(payload.Subsets())
        subset0 := payload.Subset(0)
        ncells := len(subset0.Cells())
        for i := 0; i < ncells; i++ {
            values := make([]interface{}, nsubsets)
            node := subset0.Cell(i).Node()
            for j := 0; j < nsubsets; j++ {
                values[j] = payload.Subset(j).Cell(i).Value()
            }
            if err := v.packer.Pack(node, values); err != nil {
                return err
            }
        }

    } else {
        v.packer = pack.NewUncompressedPacker(v.w)
        for _, subset := range payload.Subsets() {
            if err := subset.Accept(v); err != nil {
                return err
            }
        }
    }
    return nil
}

func (v *BinaryVisitor) VisitSubset(subset *bufr.Subset) error {
    for _, cell := range subset.Cells() {
        if err := cell.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *BinaryVisitor) VisitCell(cell *bufr.Cell) error {
    return v.packer.Pack(cell.Node(), cell.Value())
}

func (v *BinaryVisitor) VisitValuelessNode(node *bufr.ValuelessNode) error {
    panic("implement me")
}

func (v *BinaryVisitor) VisitValuedNode(node *bufr.ValuedNode) error {
    panic("implement me")
}

func (v *BinaryVisitor) VisitBlock(block *bufr.Block) error {
    panic("implement me")
}
