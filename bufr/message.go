package bufr

import (
    "fmt"
    "encoding/json"
)

// Message represents a single BUFR message which is comprised of
// multiple sections. It is the top level structure of a BUFR message.
type Message struct {
    // The metadata map is a place to store arbitrary information about a
    // message object.
    metadata map[string]interface{}

    // A list of Sections which form the message
    sections []*Section

    // The proxyFields is a way to allow easier access to fields from different
    // sections, e.g. BUFR edition number. So the consumer of the field does
    // not have to know the provider of the field (loose coupling)
    proxyFields map[string]*Field
}

func (m *Message) Accept(visitor Visitor) error {
    return visitor.VisitMessage(m)
}

func (m *Message) SetMetadata(name string, value interface{}) {
    m.metadata[name] = value
}

func (m *Message) Metadata(name string) interface{} {
    value, ok := m.metadata[name]
    if !ok {
        return nil
    }
    return value
}

func (m *Message) Sections() []*Section {
    return m.sections
}

func (m *Message) NewSection(number int, description string) *Section {
    section := NewSection(number, description)
    m.sections = append(m.sections, section)
    return section
}

func (m *Message) SetProxyField(field *Field) {
    m.proxyFields[field.Name] = field
}

func (m *Message) ProxyField(name string) (*Field, error) {
    p, ok := m.proxyFields[name]
    if ok {
        return p, nil
    }
    return nil, fmt.Errorf("no proxy field named: %s", name)
}

func (m *Message) MarshalJSON() ([]byte, error) {
    return json.Marshal(m.sections)
}

// NewMessage initialise a new Message object and set the inputPath metadata.
// Note the inputPath is only for recording purpose. The function does NOT
// try to open or read anything from the inputPath.
func NewMessage(inputPath string) *Message {
    message := Message{
        proxyFields: make(map[string]*Field),
        metadata: map[string]interface{}{
            "inputPath": inputPath,
        },
    }
    return &message
}
