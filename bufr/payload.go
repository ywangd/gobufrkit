package bufr

import "encoding/json"

// Payload represents the meat of the data section of a BUFR message.
// It is comprised of a list of Subset.
type Payload struct {
    subsets    []*Subset
    Compressed bool
}

func NewPayload(compressed bool) *Payload {
    return &Payload{Compressed: compressed}
}

func (p *Payload) Accept(visitor Visitor) error {
    return visitor.VisitPayload(p)
}

func (p *Payload) Subsets() []*Subset {
    return p.subsets
}

func (p *Payload) Subset(i int) *Subset {
    return p.subsets[i]
}

// AddSubset adds a subsets to the Payload. It also sets a zero-based
// index to the added Subset.
func (p *Payload) AddSubset(subset *Subset) {
    subset.SetIndex(len(p.subsets))
    p.subsets = append(p.subsets, subset)
}

func (p *Payload) Add(root Node, cells []*Cell) {
    subset := &Subset{
        index: len(p.subsets),
        root:  root,
        cells: cells,
    }
    p.subsets = append(p.subsets, subset)
}

func (p *Payload) MarshalJSON() ([]byte, error) {
    return json.Marshal(p.subsets)
}

// Subset represents the data of a subset of the Payload.
type Subset struct {
    // Zero based index of the subset in relative to all subsets a Payload contains
    index int

    // A list of cells representing a structureless decoded node and its value from the source.
    cells []*Cell

    // The root/entry Node of the hierarchical data structure.
    root Node
}

func NewSubset(cells []*Cell, root Node) *Subset {
    return &Subset{cells: cells, root: root}
}

func (s *Subset) Accept(visitor Visitor) error {
    return visitor.VisitSubset(s)
}

func (s *Subset) Index() int {
    return s.index
}

func (s *Subset) SetIndex(index int) {
    s.index = index
}

func (s *Subset) Cells() []*Cell {
    return s.cells
}

func (s *Subset) Cell(i int) *Cell {
    return s.cells[i]
}

func (s *Subset) AddCell(cell *Cell) {
    s.cells = append(s.cells, cell)
}

func (s *Subset) Root() Node {
    return s.root
}

func (s *Subset) MarshalJSON() ([]byte, error) {
    subsetValues := make([]interface{}, len(s.cells))
    for i, cell := range s.cells {
        subsetValues[i] = cell.Value()
    }
    return json.Marshal(subsetValues)
}
