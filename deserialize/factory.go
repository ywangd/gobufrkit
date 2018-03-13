//go:generate stringer -type=DataType
package deserialize

import (
    "fmt"
    "io"
    "github.com/pkg/errors"
    "github.com/ywangd/gobufrkit/table"
    "github.com/ywangd/gobufrkit/tdcfio"
    "github.com/ywangd/gobufrkit/bufr"
    "github.com/ywangd/gobufrkit/deserialize/payload"
    "github.com/ywangd/gobufrkit/deserialize/parser"
    "github.com/ywangd/gobufrkit/deserialize/ast"
    "os"
)

type DataType int

const (
    INT    DataType = iota
    UINT
    FLOAT
    BOOL
    BYTES
    BINARY
)

// Factory provides enabling operations for building BUFR object.
// The factory object is to be driven by Lua scripts.
type Factory interface {
    // Message returns the message that's currently being built.
    Message() *bufr.Message

    // NewMessage initialises a new message that will be built by subsequent operations.
    // The new message becomes the current message
    NewMessage(inputPath string) *bufr.Message

    // NewSection adds a new section to the current message that is being built.
    // The new section becomes the current section
    NewSection(number int, description string) *bufr.Section

    // InitTableGroup initialise the table group using the given information
    InitTableGroup(masterTableNo, centreNo, subCentreNo, wmoVersion, localVersion int) error

    // NewField creates a new Field and returns it (does not add to the current section)
    NewField(name string, dataType DataType, nbits int, proxy bool) (*bufr.Field, error)

    // NewTemplateField creates a new field holding the template values and returns it.
    // The number of bits for F, X, Y is passed as fbits, xbits and ybits.
    // The section length is to ensure that the read will stop at the section boundary
    // (as the template field is the last field of a section).
    NewTemplateField(name string, fbits, xbits, ybits int, sectionLengthInBytes uint) (*bufr.Field, error)

    // NewPayloadField creates a new field holding payload data and returns it.
    NewPayloadField(name string, nsubsets int, compressed bool) (*bufr.Field, error)

    // Padding creates a new field by reading remaining bits in the current section and returns it.
    Padding(sectionLengthInBytes uint) (*bufr.Field, error)

    // CheckEOF checks whether the EOF is reached.
    CheckEOF() (bool, error)

    // PeekEditionNumber peeks the edition number from the input stream without advancing the read head.
    PeekEditionNumber() (uint, error)

    // SeekStartSignature read the input stream until the start signature is found.
    SeekStartSignature() error
}

type DefaultFactory struct {
    config *Config
    // For reading data
    r tdcfio.PeekableReader
    // current message being built
    message *bufr.Message
    // current section being built
    section *bufr.Section
    // current unexpanded template
    ut *table.UnexpandedTemplate

    // table group for lookup descriptors
    tableGroup table.TableGroup
}

func NewDefaultFactory(config *Config, r tdcfio.PeekableReader) *DefaultFactory {
    return &DefaultFactory{config: config, r: r}
}

func (fac *DefaultFactory) Message() *bufr.Message {
    return fac.message
}

func (fac *DefaultFactory) NewMessage(inputPath string) *bufr.Message {
    fac.message = bufr.NewMessage(inputPath)
    return fac.message
}

func (fac *DefaultFactory) NewSection(number int, description string) *bufr.Section {
    fac.section = fac.message.NewSection(number, description)
    fac.section.StartByteIndex = fac.r.Pos() / tdcfio.NBITS_PER_BYTE
    return fac.section
}

func (fac *DefaultFactory) InitTableGroup(masterTableNo, centreNo, subCentreNo, wmoVersion, localVersion int) error {
    ctg := table.NewChainingTableGroup(fac.config.TablesPath)
    if err := ctg.AddLocalAndWmoTableGroups(
        masterTableNo, centreNo, subCentreNo, wmoVersion, localVersion); err != nil {
        return err
    }
    fac.tableGroup = ctg
    return nil
}

func (fac *DefaultFactory) NewField(name string, dataType DataType, nbits int, proxy bool) (*bufr.Field, error) {

    var (
        value interface{}
        err   error
    )
    switch dataType {
    case UINT:
        value, err = fac.r.ReadUint(nbits)
    case INT:
        value, err = fac.r.ReadInt(nbits)
    case FLOAT:
        value, err = fac.r.ReadFloat32()
    case BOOL:
        value, err = fac.r.ReadBool()
    case BYTES:
        value, err = fac.r.ReadBytes(nbits / tdcfio.NBITS_PER_BYTE)
    case BINARY:
        value, err = fac.r.ReadBinary(nbits)
    default:
        err = fmt.Errorf("unsupported data type: %s", dataType)
    }

    if err != nil {
        return nil, err
    }

    field := bufr.NewField(name, value, nbits)
    fac.section.AddField(field)
    if proxy {
        fac.message.SetProxyField(field)
    }

    return field, nil
}

func (fac *DefaultFactory) NewTemplateField(
    name string, fbits, xbits, ybits int, sectionLengthInBytes uint) (*bufr.Field, error) {
    remainingBits := int(sectionLengthInBytes)*tdcfio.NBITS_PER_BYTE -
        (fac.r.Pos() - fac.section.StartByteIndex*tdcfio.NBITS_PER_BYTE)

    idBits := fbits + xbits + ybits               // number of bits required for a Template ID
    ids := make([]table.ID, remainingBits/idBits) // number of IDs to read

    for i := 0; i < len(ids); i++ {
        f, err := fac.r.ReadUint(fbits)
        if err != nil {
            return nil, err
        }
        x, err := fac.r.ReadUint(xbits)
        if err != nil {
            return nil, err
        }
        y, err := fac.r.ReadUint(ybits)
        if err != nil {
            return nil, err
        }
        ids[i] = table.ID(f*100000 + x*1000 + y)
    }
    fac.ut = table.NewUnexpandedTemplate(ids, fbits, xbits, ybits)
    field := bufr.NewField(name, fac.ut, len(ids)*idBits)
    fac.section.AddField(field)
    fac.message.SetProxyField(field)

    return field, nil
}

func (fac *DefaultFactory) NewPayloadField(name string, nsubsets int, compressed bool) (*bufr.Field, error) {
    spos := fac.r.Pos()
    tree, err := parser.NewParser(fac.tableGroup).Parse(fac.ut)
    if err != nil {
        return nil, errors.Wrap(err, "cannot parse template")
    }

    if fac.config.Verbose {
        v := ast.DumpVisitor(os.Stdout)
        tree.Accept(v)
    }

    desvis, err := payload.NewDeserializeVisitor(fac.config.toDesVisitorConfig(compressed), fac.r, nsubsets)
    if err != nil {
        return nil, err
    }

    n := nsubsets
    if compressed && fac.config.InputType == tdcfio.BinaryInput {
        n = 1 // values from all subsets are packed together for compressed binary data
    }
    pay := &bufr.Payload{Compressed: compressed}
    for i := 0; i < n; i++ {
        if err := tree.Accept(desvis); err != nil {
            return nil, err
        }
        if err := desvis.Produce(pay); err != nil {
            return nil, err
        }
    }
    field := bufr.NewField(name, pay, fac.r.Pos()-spos)
    fac.section.AddField(field)
    fac.message.SetProxyField(field)
    return field, nil
}

func (fac *DefaultFactory) Padding(sectionLengthInBytes uint) (*bufr.Field, error) {
    bitsTotal := int(sectionLengthInBytes) * tdcfio.NBITS_PER_BYTE
    bitsRead := fac.r.Pos() - fac.section.StartByteIndex*tdcfio.NBITS_PER_BYTE
    bitsPadding := int(bitsTotal - bitsRead)

    switch {
    case bitsPadding == 0:
        return nil, nil
    case bitsPadding < 0:
        return nil, fmt.Errorf("read beyond section boundary by %d bits", -bitsPadding)
    }

    fac.section.Padding = bitsPadding
    binary, err := fac.r.ReadBinary(int(bitsPadding))
    if err != nil {
        return nil, errors.Wrap(err, "cannot read padding")
    }
    field := bufr.NewHiddenField("padding", binary, binary.Nbits())
    fac.section.AddField(field)
    return field, nil

}

func (fac *DefaultFactory) CheckEOF() (bool, error) {
    _, err := fac.r.PeekUint(0, 8)
    if err == io.EOF {
        return true, nil
    }
    return false, err
}

// TODO: BUFR structural info leak
func (fac *DefaultFactory) PeekEditionNumber() (uint, error) {
    v, err := fac.r.PeekUint(7, 8)
    if err != nil {
        return 0, err
    }
    return v, nil
}

func (fac *DefaultFactory) SeekStartSignature() error {
    for {
        bs, err := fac.r.PeekBytes(0, 4)
        if err != nil {
            return err
        }
        if len(bs) < 4 {
            return io.EOF
        } else if string(bs) == "BUFR" {
            return nil
        }
        fac.r.ReadBytes(1)
    }
}
