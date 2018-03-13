package ast

import (
    "fmt"
    "github.com/ywangd/gobufrkit/table"
)

type LintError struct {
    message string
}

func (e *LintError) Error() string {
    return e.message
}

func lintError(format string, args ...interface{}) *LintError {
    return &LintError{message: fmt.Sprintf(format, args...)}
}

type LintVisitor struct {
}

func (v *LintVisitor) VisitNode(node Node) error {
    panic("implement me")
}

func (v *LintVisitor) VisitElementNode(node *ElementNode) error {
    if node.Descriptor().F() != 0 {
        return lintError("not an element descriptor: %v", node.Descriptor())
    }
    return nil
}

func (v *LintVisitor) VisitE031021Node(node *E031021Node) error {
    panic("implement me")
}

func (v *LintVisitor) VisitFixedReplicationNode(node *FixedReplicationNode) error {
    if node.Descriptor().Y() != len(node.Members()) {
        return lintError("incorrect number of replicated members: expect: %v, got: %v",
            node.descriptor.Y(), len(node.Members()))
    }
    return nil
}

func (v *LintVisitor) VisitDelayedReplicationNode(node *DelayedReplicationNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitSequenceNode(node *SequenceNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpNbitsOffsetNode(node *OpNbitsOffsetNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpScaleOffsetNode(node *OpScaleOffsetNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpNewRefvalNode(node *OpNewRefvalNode) error {
    if node.Descriptor().Y() != 255 {
        for _, m := range node.Members() {
            if m.Descriptor().F() != 0 {
                return lintError(
                    "non-element descriptor appears in a new refval definition session: %v",
                    m.Descriptor())
            }
        }
    }
    return nil
}

func (v *LintVisitor) VisitOpAssocFieldNode(node *OpAssocFieldNode) error {
    if node.Descriptor().Y() == 0 {
        return nil
    }
    members := node.Members()
    if len(members) != 1 {
        return lintError("incorrect number of associated field significance")
    }
    if members[0].Descriptor().Id() != table.ID_031021 {
        return lintError("invalid associated field significance: expected 031021, got %v",
            members[0].Descriptor().Id())
    }
    return nil
}

func (v *LintVisitor) VisitOpInsertStringNode(node *OpInsertStringNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpSkipLocalNode(node *OpSkipLocalNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpModifyPackingNode(node *OpModifyPackingNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpSetStringLengthNode(node *OpSetStringLengthNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpDataNotPresentNode(node *OpDataNotPresentNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpAssessmentNode(node *OpAssessmentNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpMarkerNode(node *OpMarkerNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpCancelBackRefNode(node *OpCancelBackRefNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitOpCancelBitmapNode(node *OpCancelBitmapNode) error {
    panic("implement me")
}

func (v *LintVisitor) VisitBitmapNode(node *BitmapNode) error {
    panic("implement me")
}
