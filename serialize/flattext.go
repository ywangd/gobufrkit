package serialize

import (
    "github.com/ywangd/gobufrkit/bufr"
    "io"
    "fmt"
)

type FlatTextVisitor struct {
    w io.Writer

    ShowHidden bool
}

func NewFlatTextVisitor(w io.Writer) *FlatTextVisitor {
    return &FlatTextVisitor{w: w, ShowHidden: true}
}

func (v *FlatTextVisitor) VisitMessage(message *bufr.Message) error {
    for i, section := range message.Sections() {
        if _, err := fmt.Fprintf(v.w, "<<<<<< section %d >>>>>>\n", i); err != nil {
            return err
        }
        if err := section.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *FlatTextVisitor) VisitSection(section *bufr.Section) error {
    for _, field := range section.Fields() {
        if err := field.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *FlatTextVisitor) VisitField(field *bufr.Field) error {
    if field.Hidden && !v.ShowHidden {
        return nil
    }
    switch value := field.Value.(type) {
    case *bufr.Payload:
        return value.Accept(v)
    default:
        _, err := fmt.Fprintf(v.w, "%s = %v\n", field.Name, field.Value)
        return err
    }
}

func (v *FlatTextVisitor) VisitPayload(payload *bufr.Payload) error {
    subsets := payload.Subsets()
    for i, subset := range subsets {
        if _, err := fmt.Fprintf(v.w, "###### subset %d of %d ######\n", i+1, len(subsets)); err != nil {
            return err
        }
        if err := subset.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *FlatTextVisitor) VisitSubset(subset *bufr.Subset) error {
    for i, cell := range subset.Cells() {
        if _, err := fmt.Fprintf(v.w, "%5d ", i+1); err != nil {
            return err
        }
        if err := cell.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *FlatTextVisitor) VisitCell(cell *bufr.Cell) error {
    var s string
    value := cell.Value()
    switch value.(type) {
    case []byte:
        s = fmt.Sprintf("%q", string(value.([]byte)))
    default:
        s = fmt.Sprintf("%v", value)
    }

    _, err := fmt.Fprintf(v.w, "%-60s%v\n", cell.Node().Descriptor, s)
    return err
}

func (v *FlatTextVisitor) VisitValuelessNode(node *bufr.ValuelessNode) error {
    panic("implement me")
}

func (v *FlatTextVisitor) VisitValuedNode(node *bufr.ValuedNode) error {
    panic("implement me")
}

func (v *FlatTextVisitor) VisitBlock(block *bufr.Block) error {
    panic("implement me")
}
