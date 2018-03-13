package bufr

type Visitor interface {
    VisitMessage(message *Message) error
    VisitSection(section *Section) error
    VisitField(field *Field) error
    VisitPayload(payload *Payload) error
    VisitSubset(subset *Subset) error
    VisitCell(cell *Cell) error
    VisitValuelessNode(node *ValuelessNode) error
    VisitValuedNode(node *ValuedNode) error
    VisitBlock(block *Block) error
}

type Acceptor interface {
    Accept(visitor Visitor) error
}