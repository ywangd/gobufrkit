//go:generate stringer -type=Unit
package table

import (
    "fmt"
    "os"
    "encoding/csv"
    "strconv"
    "strings"
)

// TODO: move this to bufr
type Unit int

const (
    NUMERIC Unit = iota
    STRING
    NONNEG_CODE  // non-negative code, this covers all the descriptors of code unit
    CODE // full range code
    FLAG
    BINARY
)

// Convert the unit description in table records to Normalised Unit.
func unitOf(s string) (unit Unit) {
    switch {
    case s == "CCITT IA5" || s == "Character":
        unit = STRING
    case s == "FLAG TABLE":
        unit = FLAG
    case s == "CODE TABLE" || strings.HasPrefix(s, "Common CODE TABLE"):
        unit = NONNEG_CODE
    default:
        unit = NUMERIC
    }
    return
}

// Entry is a basic interface representing additional information about a descriptor
// that is available in tables. It is useful to provide a common interface to different
// meanings that different type of descriptors may have.
type Entry interface {
    Name() string
}

// Bentry is the information about an Elemental descriptor that can be found in table B.
type Bentry struct {
    name string

    UnitString string
    Unit       Unit
    Scale      int
    Refval     int
    Nbits      int

    CrexUnitString string
    CrexUnit       Unit
    CrexScale      int
    CrexNchars     int
}

func (e *Bentry) Name() string {
    return e.name
}

// Rentry is a placeholder entry for replication descriptor
type Rentry struct {
    name string
}

func (e *Rentry) Name() string {
    return e.name
}

// Centry represents the entry for operator descriptor
type Centry struct {
    name string
}

func (e *Centry) Name() string {
    return e.name
}

// Dentry represents the information about Sequence Descriptor that is
// available in Table D.
type Dentry struct {
    name string

    // The expanded member IDs of the sequence descriptor
    Members []ID
}

func (e *Dentry) Name() string {
    return e.name
}

// Table is the basic interface representing a BUFR table.
type Table interface {
    // Get back a descriptor given its ID.
    Lookup(id ID) (Descriptor, error)
}

// B represents a BUFR Table B deserialised from an input CSV file.
type B struct {
    // Path to the input file.
    path string

    // Each non-comment line in the input file is converted to one Bentry
    // with the key being the corresponding ID.
    entries map[ID]*Bentry
}

func (b *B) Lookup(id ID) (Descriptor, error) {
    entry, ok := b.entries[id]
    if ok {
        return NewElementDescriptor(id, entry), nil
    }
    return nil, fmt.Errorf("ID not found: %s", id)
}

// LoadTableB build a Table B by reading the given input file.
func LoadTableB(tablePath string) (*B, error) {
    ins, err := os.Open(tablePath)
    defer ins.Close()
    if err != nil {
        return nil, err
    }

    r := csv.NewReader(ins)
    r.Comment = '#'

    records, err := r.ReadAll()
    if err != nil {
        return nil, err
    }

    entries := make(map[ID]*Bentry, len(records))

    for _, record := range records {
        id, err := strconv.Atoi(record[0])
        if err != nil {
            return nil, err
        }
        entry, err := recordToBentry(record)
        if err != nil {
            return nil, err
        }
        entries[ID(id)] = entry
    }
    return &B{path: tablePath, entries: entries}, nil
}

// D represents a single BUFR Table D deserialised from an input CSV file.
type D struct {
    // Path to the input file
    path string

    // Each non-commment line from the input file is converted to a Dentry
    // indexed by the corresponding ID.
    entries map[ID]*Dentry
}

func (d *D) Lookup(id ID) (Descriptor, error) {
    entry, ok := d.entries[id]
    if ok {
        return NewSequenceDescriptor(id, entry), nil
    }
    return nil, fmt.Errorf("ID not found: %s", id)
}

// Build a Table D from the given input file.
func LoadTableD(tablePath string) (*D, error) {
    ins, err := os.Open(tablePath)
    defer ins.Close()
    if err != nil {
        return nil, err
    }

    r := csv.NewReader(ins)
    r.Comment = '#'

    records, err := r.ReadAll()
    if err != nil {
        return nil, err
    }

    entries := make(map[ID]*Dentry, len(records))

    for _, record := range records {
        id, err := strconv.Atoi(record[0])
        if err != nil {
            return nil, err
        }
        entry, err := recordToDentry(record)
        if err != nil {
            return nil, err
        }
        entries[ID(id)] = entry
    }
    return &D{path: tablePath, entries: entries}, nil
}

// recordToBentry is a helper function that convert a list of string
// reading from input CSV file to a Bentry.
func recordToBentry(record []string) (*Bentry, error) {

    scale, err := strconv.Atoi(record[3])
    if err != nil {
        return nil, err
    }
    refval, err := strconv.Atoi(record[4])
    if err != nil {
        return nil, err
    }
    nbits, err := strconv.Atoi(record[5])
    if err != nil {
        return nil, err
    }

    crexScale, err := strconv.Atoi(record[7])
    if err != nil {
        return nil, err
    }
    crexNchars, err := strconv.Atoi(record[8])
    if err != nil {
        return nil, err
    }
    return &Bentry{
        name:           record[1],
        UnitString:     record[2],
        Unit:           unitOf(record[2]),
        Scale:          scale,
        Refval:         refval,
        Nbits:          nbits,
        CrexUnitString: record[6],
        CrexUnit:       unitOf(record[6]),
        CrexScale:      crexScale,
        CrexNchars:     crexNchars,
    }, nil
}

// recordToDentry is a helper function that converts a list of string
// reading form input CSV file to a Dentry.
func recordToDentry(record []string) (*Dentry, error) {
    fields := strings.Split(record[2], ",")
    ids := make([]ID, len(fields))

    for i, idstring := range fields {
        id, err := strconv.Atoi(idstring)
        if err != nil {
            return nil, err
        }
        ids[i] = ID(id)
    }
    return &Dentry{name: record[1], Members: ids}, nil
}
