package table

import (
    "fmt"
    "os"
    "log"
)

// TableGroup is a group of related tables, e.g. Tables of the same version number.
type TableGroup interface {
    // Get a descriptor for the given ID. The descriptor can from any table
    // of the group.
    Lookup(id ID) (Descriptor, error)
}

// SingleTableGroup is a group of tables that belong to the same centre and
// have the same version number.
type SingleTableGroup struct {
    b *B
    d *D
}

// NewSingleTableGroup creates a new table group using the given arguments to
// request corresponding tables. The tables are requested through the singleton
// TableManager object.
func NewSingleTableGroup(tablesBasePath string,
    masterTableNumber, centreNumber, subCentreNumber, versionNumber int) (TableGroup, error) {
    // Load Table B
    b, err := manager.getTableB(tablesBasePath,
        masterTableNumber, centreNumber, subCentreNumber, versionNumber)
    if err != nil {
        return nil, err
    }
    // Load Table D
    d, err := manager.getTableD(tablesBasePath,
        masterTableNumber, centreNumber, subCentreNumber, versionNumber)
    if err != nil {
        return nil, err
    }

    return &SingleTableGroup{b: b, d: d}, nil
}

func (tg *SingleTableGroup) Lookup(id ID) (Descriptor, error) {
    switch id.F() {
    case F_ELEMENT:
        // TODO: handle error and allow local descriptor (should be handled by chaining group)
        return tg.b.Lookup(id)
    case F_REPLICATION:
        // TODO: entry info
        return &ReplicationDescriptor{BaseDescriptor{id, &Rentry{name: id.String()}}}, nil
    case F_OPERATOR:
        // TODO: proper entry
        return &OperatorDescriptor{BaseDescriptor{id, &Centry{name: id.String()}}}, nil
    case F_SEQUENCE:
        return tg.d.Lookup(id)
    default:
        return nil, fmt.Errorf("unknown ID: %d", id)
    }
}

// ChainingTableGroup a meta TableGroup that is support by a list of member
// TableGroup. All member groups must be constructed from tables located in
// the same base path.
type ChainingTableGroup struct {
    tablesBasePath string
    groups         []TableGroup
}

func (ctg *ChainingTableGroup) Lookup(id ID) (Descriptor, error) {
    // Loop through the groups and return the first descriptor that can be retrieved
    for _, g := range ctg.groups {
        descriptor, err := g.Lookup(id)
        if err == nil {
            return descriptor, nil
        }
    }
    return nil, fmt.Errorf("ID not found: %v", id)
}

// Add a SingleTableGroup as a member
func (ctg *ChainingTableGroup) AddSingleTableGroup(
    masterTableNumber, centreNumber, subCentreNumber, versionNumber int) error {
    g, err := NewSingleTableGroup(
        ctg.tablesBasePath,
        masterTableNumber, centreNumber, subCentreNumber, versionNumber,
    )
    if err != nil {
        return err
    }

    ctg.groups = append(ctg.groups, g)
    return nil
}

// Add two SingleTableGroup with one being used locally by the centre and the
// other being used by WMO. The two groups share the same masterTableNumber.
// If a local table group does not exist, only the WMO table group will be added.
func (ctg *ChainingTableGroup) AddLocalAndWmoTableGroups(
    masterTableNumber, centreNumber, subCentreNumber, wmoVersionNumber, localVersionNumber int) error {
    if localVersionNumber != 0 {
        err := ctg.AddSingleTableGroup(masterTableNumber, centreNumber, subCentreNumber, localVersionNumber)
        if os.IsNotExist(err) {
            log.Println("Warning: ", err.Error(), "Fallback to subCentreNumber 0")
            if err := ctg.AddSingleTableGroup(masterTableNumber, centreNumber,
                0, localVersionNumber); err != nil {
                return err
            }
        } else if err != nil {
            return err
        }
    }
    err := ctg.AddSingleTableGroup(masterTableNumber, 0, 0, wmoVersionNumber)
    return err
}

// ResetGroups clears the internal list of member groups.
func (ctg *ChainingTableGroup) ResetGroups() {
    ctg.groups = []TableGroup{}
}

func NewChainingTableGroup(tablesBasePath string) *ChainingTableGroup {
    return &ChainingTableGroup{tablesBasePath: tablesBasePath, groups: []TableGroup{}}
}
