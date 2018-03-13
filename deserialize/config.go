package deserialize

import (
    "github.com/ywangd/gobufrkit/tdcfio"
    "github.com/ywangd/gobufrkit/deserialize/payload"
)

type Config struct {
    TablesPath string
    InputType  tdcfio.InputType
    Compatible bool
    Verbose    bool
}

func (c *Config) toDesVisitorConfig(compressed bool) *payload.DesVisitorConfig {
    return &payload.DesVisitorConfig{
        Compressed: compressed,
        InputType:  c.InputType,
        Compatible: c.Compatible,
        Verbose:    c.Verbose,
    }
}
