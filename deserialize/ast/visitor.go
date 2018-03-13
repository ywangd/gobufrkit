package ast

// Visitor is the interface that work with all nodes of a parsed template tree.
// Different implementations of this interface can provide difference services.
// A few builtin implementations are: dump the tree structure, lint the template
// tree for semantics, deserialization based on the tree.
type Visitor interface {
    VisitNode(node Node) error
    VisitElementNode(node *ElementNode) error
    VisitE031021Node(node *E031021Node) error
    VisitFixedReplicationNode(node *FixedReplicationNode) error
    VisitDelayedReplicationNode(node *DelayedReplicationNode) error
    VisitSequenceNode(node *SequenceNode) error

    VisitOpNbitsOffsetNode(node *OpNbitsOffsetNode) error
    VisitOpScaleOffsetNode(node *OpScaleOffsetNode) error
    VisitOpNewRefvalNode(node *OpNewRefvalNode) error

    VisitOpAssocFieldNode(node *OpAssocFieldNode) error
    VisitOpInsertStringNode(node *OpInsertStringNode) error
    VisitOpSkipLocalNode(node *OpSkipLocalNode) error
    VisitOpModifyPackingNode(node *OpModifyPackingNode) error
    VisitOpSetStringLengthNode(node *OpSetStringLengthNode) error

    VisitOpDataNotPresentNode(node *OpDataNotPresentNode) error
    VisitOpAssessmentNode(node *OpAssessmentNode) error
    VisitOpMarkerNode(node *OpMarkerNode) error
    VisitOpCancelBackRefNode(node *OpCancelBackRefNode) error
    VisitOpCancelBitmapNode(node *OpCancelBitmapNode) error

    VisitBitmapNode(node *BitmapNode) error
}
