package bufr

import (
    "log"
    "encoding/json"
)

// Section represents a section in a BUFR message. It is the second level
// structure of a BUFR message. It is comprised of a list of Fields.
type Section struct {
    // The zero based index of the start byte of the section.
    StartByteIndex int

    // Binary padding at the end of the section as an integer
    Padding int

    // Zero based section number. This is defined by the BUFR manual.
    number int

    // A place to store arbitrary key/value pair information.
    metadata map[string]interface{}

    // The list of Field that form the section.
    fields []*Field
}

func (s *Section) Accept(visitor Visitor) error {
    return visitor.VisitSection(s)
}

func (s *Section) Number() int {
    return s.number
}

func (s *Section) SetMetadata(name string, value interface{}) {
    s.metadata[name] = value
}

func (s *Section) Metadata(name string) interface{} {
    value, ok := s.metadata[name]
    if !ok {
        // TODO: return an error instead
        log.Fatalf("No metadata of name: %s\n", name)
    }
    return value
}

func (s *Section) Fields() []*Field {
    return s.fields
}

func (s *Section) AddField(field *Field) {
    s.fields = append(s.fields, field)
}

// Get a field using the given name. Returns the first match or nil if no match
// is found.
func (s *Section) FieldByName(name string) *Field {
    for _, p := range s.fields {
        if p.Name == name {
            return p
        }
    }
    return nil
}

func (s *Section) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.fields)
}

func NewSection(number int, description string) *Section {
    return &Section{
        number: number,
        metadata: map[string]interface{}{
            "description": description,
        },
    }
}
