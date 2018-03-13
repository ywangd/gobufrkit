package serialize

import (
    "io"
    "github.com/ywangd/gobufrkit/bufr"
    "github.com/ywangd/gobufrkit/table"
    "encoding/json"
)

type FlatJsonVisitor struct {
    w   io.Writer
    enc *json.Encoder
}

func NewFlatJsonVisitor(w io.Writer) *FlatJsonVisitor {
    return &FlatJsonVisitor{w: w, enc: json.NewEncoder(w)}
}

func (v *FlatJsonVisitor) VisitMessage(message *bufr.Message) error {
    for _, section := range message.Sections() {
        if err := section.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *FlatJsonVisitor) VisitSection(section *bufr.Section) error {
    for _, field := range section.Fields() {
        if err := field.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *FlatJsonVisitor) VisitField(field *bufr.Field) error {
    switch value := field.Value.(type) {
    case *bufr.Payload:
        return value.Accept(v)
    case *table.UnexpandedTemplate:
        v.enc.Encode(value.Ids())
    }
    panic("implement me")
}

func (v *FlatJsonVisitor) VisitPayload(payload *bufr.Payload) error {
    panic("implement me")
}

func (v *FlatJsonVisitor) VisitSubset(subset *bufr.Subset) error {
    panic("implement me")
}

func (v *FlatJsonVisitor) VisitCell(cell *bufr.Cell) error {
    panic("implement me")
}

func (v *FlatJsonVisitor) VisitValuelessNode(node *bufr.ValuelessNode) error {
    panic("implement me")
}

func (v *FlatJsonVisitor) VisitValuedNode(node *bufr.ValuedNode) error {
    panic("implement me")
}

func (v *FlatJsonVisitor) VisitBlock(block *bufr.Block) error {
    panic("implement me")
}
