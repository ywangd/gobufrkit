package ast

import (
    "io"
    "fmt"
)

type dumpVisitor struct {
    w      io.Writer
    prefix string
}

func DumpVisitor(w io.Writer) *dumpVisitor {
    return &dumpVisitor{w: w}
}

func (v *dumpVisitor) VisitNode(node Node) error {
    if err := v.printf("%v (%T)\n", node.Descriptor(), node); err != nil {
        return err
    }
    return v.visitMembers(node.Members())
}

func (v *dumpVisitor) VisitElementNode(node *ElementNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitE031021Node(node *E031021Node) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitFixedReplicationNode(node *FixedReplicationNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitDelayedReplicationNode(node *DelayedReplicationNode) error {
    if err := v.printf("%v (%T)\n", node.Descriptor(), node); err != nil {
        return err
    }
    if err := v.printf("....%v\n", node.members[0].Descriptor()); err != nil {
        return err
    }
    return v.visitMembers(node.Members()[1:])
}

func (v *dumpVisitor) VisitSequenceNode(node *SequenceNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpNbitsOffsetNode(node *OpNbitsOffsetNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpScaleOffsetNode(node *OpScaleOffsetNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpNewRefvalNode(node *OpNewRefvalNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpAssocFieldNode(node *OpAssocFieldNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpInsertStringNode(node *OpInsertStringNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpSkipLocalNode(node *OpSkipLocalNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpModifyPackingNode(node *OpModifyPackingNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpSetStringLengthNode(node *OpSetStringLengthNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpCancelBackRefNode(node *OpCancelBackRefNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpDataNotPresentNode(node *OpDataNotPresentNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpAssessmentNode(node *OpAssessmentNode) error {
    if err := v.printf("%v (%T)\n", node.Descriptor(), node); err != nil {
        return err
    }
    v.indent()
    if err := node.Bitmap.Accept(v); err != nil {
        return err
    }
    v.dedent()
    if err := v.visitMembers(node.Attrs); err != nil {
        return err
    }
    return v.visitMembers(node.Members())
}

func (v *dumpVisitor) VisitOpMarkerNode(node *OpMarkerNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitOpCancelBitmapNode(node *OpCancelBitmapNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) VisitBitmapNode(node *BitmapNode) error {
    return v.VisitNode(node)
}

func (v *dumpVisitor) visitMembers(members []Node) error {
    v.indent()
    defer v.dedent()
    for _, m := range members {
        if err := m.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

// printf prefix the given output with indentation
func (v *dumpVisitor) printf(format string, args ...interface{}) error {
    if _, err := fmt.Fprint(v.w, v.prefix); err != nil {
        return err
    }
    _, err := fmt.Fprintf(v.w, format, args...)
    return err
}

func (v *dumpVisitor) indent() {
    v.prefix += "    "
}

func (v *dumpVisitor) dedent() {
    v.prefix = v.prefix[4:]
}
