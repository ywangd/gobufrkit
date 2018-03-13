package table

import (
    "sync"
    "path/filepath"
    "strconv"
)

// The tableManager represents an entry point to all tables. It provides internal
// cache for all loaded tables and only read from files when necessary.
type tableManager struct {
    // Locks preventing concurrent changes to the cache of tables.
    bmu sync.RWMutex
    dmu sync.RWMutex

    // Cache for Table B and D. The keys are input file path of each table.
    bs map[string]*B
    ds map[string]*D
}

// manager is the singleton tableManager shared by all table groups
var manager = &tableManager{
    bs: make(map[string]*B),
    ds: make(map[string]*D),
}

// Get a Table B from the path calculated using the given arguments. The retrieval
// is first attempted from the cache. It then tries to load the table from disk
// if no cached version is available. Any newly loaded table will be saved in
// the cache.
func (tm *tableManager) getTableB(tablesBasePath string,
    masterTableNumber, centreNumber, subCentreNumber, versionNumber int) (*B, error) {

    tablePath := composeTablePath(tablesBasePath,
        masterTableNumber, centreNumber, subCentreNumber, versionNumber,
        "TableB.csv")

    tm.bmu.RLock()
    b, ok := tm.bs[tablePath]
    tm.bmu.RUnlock()
    if ok {
        return b, nil
    }

    tm.bmu.Lock()
    defer tm.bmu.Unlock()
    b, ok = tm.bs[tablePath]
    if ok {
        return b, nil
    }
    b, err := LoadTableB(tablePath)
    if err != nil {
        return nil, err
    }
    tm.bs[tablePath] = b
    return b, nil
}

// Get a Table D from the path calculated using the given arguments. See also
// getTableB
func (tm *tableManager) getTableD(tablesBasePath string,
    masterTableNumber, centreNumber, subCentreNumber, versionNumber int) (*D, error) {
    tablePath := composeTablePath(tablesBasePath,
        masterTableNumber, centreNumber, subCentreNumber, versionNumber,
        "TableD.csv")

    tm.dmu.RLock()
    d, ok := tm.ds[tablePath]
    tm.dmu.RUnlock()
    if ok {
        return d, nil
    }

    tm.dmu.Lock()
    defer tm.dmu.Unlock()
    d, ok = tm.ds[tablePath]
    if ok {
        return d, nil
    }
    d, err := LoadTableD(tablePath)
    if err != nil {
        return nil, err
    }
    tm.ds[tablePath] = d
    return d, nil
}

// Calculate the table path using the given arguments.
func composeTablePath(tablesBasePath string,
    masterTableNumber, centreNumber, subCentreNumber, versionNumber int, tableName string) string {

    return filepath.Join(
        tablesBasePath,
        strconv.Itoa(masterTableNumber),
        strconv.Itoa(centreNumber),
        strconv.Itoa(subCentreNumber),
        strconv.Itoa(versionNumber),
        tableName,
    )
}
