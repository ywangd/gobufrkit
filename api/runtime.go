package api

import (
    "github.com/ywangd/gobufrkit/tdcfio"
    "github.com/ywangd/gobufrkit/deserialize"
    "github.com/pkg/errors"
    "github.com/ywangd/gobufrkit/bufr"
)

type Config struct {
    DefinitionsPath string
    TablesPath      string

    // Only binary stream provides compressed data
    // in the format described by the BUFR Spec.
    InputType  tdcfio.InputType
    Compatible bool
    Verbose    bool
}

func (c *Config) toDeserializeConfig() *deserialize.Config {
    return &deserialize.Config{
        TablesPath: c.TablesPath,
        InputType:  c.InputType,
        Compatible: c.Compatible,
        Verbose:    c.Verbose,
    }
}

type Runtime struct {
    config   *Config
    scriptRt *ScriptRt
}

func NewRuntime(config *Config, pr tdcfio.PeekableReader) (*Runtime, error) {

    factory := deserialize.NewDefaultFactory(config.toDeserializeConfig(), pr)

    scriptRt := NewScriptRt(config.DefinitionsPath, factory)
    if err := scriptRt.Initialize(); err != nil {
        return nil, errors.Wrap(err, "cannot initialise script runtime")
    }

    return &Runtime{
        config:   config,
        scriptRt: scriptRt,
    }, nil
}

func (rt *Runtime) Run() (*bufr.Message, error) {
    return rt.scriptRt.RunDeserializer()
}
