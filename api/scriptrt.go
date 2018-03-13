package api

import (
    "io"
    "os"
    "path/filepath"
    "github.com/Shopify/go-lua"
    "github.com/ywangd/gobufrkit/deserialize"
    "github.com/ywangd/gobufrkit/bufr"
    "github.com/ywangd/gobufrkit/tdcfio"
)

const (
    RUNTIME_METATABLE = "gobufrkit.runtime"
    DESERIALIZER      = "DESERIALIZER"
)

type ScriptRt struct {
    definitionsPath string
    state           *lua.State
    factory         deserialize.Factory
}

func NewScriptRt(definitionsPath string, factory deserialize.Factory) *ScriptRt {
    state := lua.NewState()
    return &ScriptRt{definitionsPath: definitionsPath, state: state, factory: factory}
}

// Initialize the runtime. This will setup the bindings and perform other
// initialization stuff.
func (r *ScriptRt) Initialize() error {
    lua.OpenLibraries(r.state)
    r.initConstants()
    r.initPkgPath()
    if err := r.initScripts(); err != nil {
        return err
    }
    r.initLibs()
    r.initTypes()
    return nil
}

// Add a new global user data using the given name and metatable
func (r *ScriptRt) NewGlobalUserData(name string, value interface{}, metaName string) {
    r.state.PushUserData(value)
    lua.SetMetaTableNamed(r.state, metaName)
    r.state.SetGlobal(name)
}

// Add a new global simple value using the given name
func (r *ScriptRt) NewGlobalSimpleValue(name string, value interface{}) {
    pushSimpleValue(r.state, value)
    r.state.SetGlobal(name)
}

// Load a script and save as a field to the runtime metatable
func (r *ScriptRt) Load(reader io.Reader, chunkName, keyName string) error {
    lua.NewMetaTable(r.state, RUNTIME_METATABLE)
    err := r.state.Load(reader, "="+chunkName, "")
    if err != nil {
        return err
    }
    r.state.SetField(-2, keyName)
    r.state.Pop(1)
    return nil
}

// Get the decoder func from the runtime metatable and execute it
func (r *ScriptRt) RunDeserializer() (*bufr.Message, error) {
    lua.MetaTableNamed(r.state, RUNTIME_METATABLE)
    r.state.Field(-1, DESERIALIZER)
    r.state.Remove(-2)
    // TODO: check Start signature (probably in lua part)
    if err := r.state.ProtectedCall(0, 1, 0); err != nil {
        return nil, err
    }
    message := r.state.ToUserData(-1).(*bufr.Message)
    r.state.Pop(1)
    return message, nil
}

// Add useful global constant values
func (r *ScriptRt) initConstants() {
    // Data types
    r.NewGlobalSimpleValue("INT", int(deserialize.INT))
    r.NewGlobalSimpleValue("UINT", int(deserialize.UINT))
    r.NewGlobalSimpleValue("FLOAT", int(deserialize.FLOAT))
    r.NewGlobalSimpleValue("BOOL", int(deserialize.BOOL))
    r.NewGlobalSimpleValue("BYTES", int(deserialize.BYTES))
    r.NewGlobalSimpleValue("BINARY", int(deserialize.BINARY))

    r.NewGlobalSimpleValue("BITS_PER_BYTE", tdcfio.NBITS_PER_BYTE)
}

// Initialise the package search path
func (r *ScriptRt) initPkgPath() {
    r.state.Global("package")
    r.state.Field(-1, "path")
    path, _ := r.state.ToString(-1)
    r.state.Pop(1)
    r.state.PushString(filepath.Join(r.definitionsPath, "?.lua;") + path)
    r.state.SetField(-2, "path")
    r.state.Pop(1)
}

// Load the lua scripts for deserializer
func (r *ScriptRt) initScripts() error {
    scriptFileName := filepath.Join(r.definitionsPath, "boot.lua")
    ins, err := os.Open(scriptFileName)
    if err != nil {
        return err
    }
    defer ins.Close()
    return r.Load(ins, scriptFileName, DESERIALIZER)
}

func (r *ScriptRt) initTypes() {
    (&ObjMessage{}).registerGribMessageType(r.state)
    (&ObjSection{}).registerGribSectionType(r.state)
    (&ObjField{}).registerSectionFieldType(r.state)
}

// initialise local libraries
func (r *ScriptRt) initLibs() {
    lua.Require(r.state, "factory",
        (&LibDeserializer{factory: r.factory}).Register, true)
}
