package cmd

import (
    "os"
    "log"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/ywangd/gobufrkit/api"
    "github.com/ywangd/gobufrkit/tdcfio"
    "path/filepath"
    "github.com/ywangd/gobufrkit/serialize"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
    Use:     "decode [filename]",
    Short:   "Decode from a BUFR file or STDIN if no file is given.",
    Long:    `Decode from a BUFR file or STDIN if no file is given.`,
    Aliases: []string{"d"},
    Args:    cobra.MaximumNArgs(1),
    Run:     runDecode,
}

func init() {
    RootCmd.AddCommand(decodeCmd)
    decodeCmd.Flags().BoolP("first-message", "1", false, "Decode only the first message")
    decodeCmd.Flags().BoolP("attributed", "a", false, "Output attributed hierarchical structure")
    decodeCmd.Flags().BoolP("json", "j", false, "Output as bare JSON format")
    decodeCmd.Flags().BoolP("show-hidden-fields", "x", false, "Show hidden fields, e.g. padding")
}

func runDecode(cmd *cobra.Command, args []string) {

    // Command line argument processing
    firstMessage := cmd.Flag("first-message").Changed

    // Open the input BUFR file
    var (
        ins *os.File
        err error
    )
    if len(args) > 0 {
        ins, err = os.Open(args[0])
        if err != nil {
            log.Fatal(err.Error())
        }
        defer ins.Close()
    } else {
        ins = os.Stdin
    }

    pr := tdcfio.NewPeekableBitReader(ins)
    definitionsPath := viper.GetString("definitions_path")
    tablesPath := filepath.Join(definitionsPath, "tables")

    config := &api.Config{
        DefinitionsPath: definitionsPath,
        TablesPath:      tablesPath,
        InputType:       tdcfio.BinaryInput,
        Compatible:      cmd.Flag("compatible").Changed,
        Verbose:         cmd.Flag("debug").Changed,
    }

    rt, err := api.NewRuntime(config, pr)
    if err != nil {
        log.Fatal(err.Error())
    }

    flatText := serialize.NewFlatTextSerializer(os.Stdout)

    for i := 0; ; i++ {
        message, err := rt.Run()
        if err != nil {
            log.Fatal("Lua error:", err.Error())
        }
        message.SetMetadata("number", i+1)

        flatText.Serialize(message)

        break // TODO: debug
        if firstMessage {
            break
        }
    }
}
