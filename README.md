# GoBufrKit

An unfinished project for implementing WMO [BUFR](https://en.wikipedia.org/wiki/BUFR) 
decoder in [Go](https://golang.org/). Build the binary with `go build` or directly 
run with `go run main.go`.

The current code is able to decode most BUFR messages. The output format is plain
text only. JSON output is almost there, as well as binary output, i.e. encoder.
The intention was to make a faster alternative to [PyBufrKit](https://github.com/ywangd/pybufrkit).
But I cannot see myself working on this project anytime soon. Adoptions are welcome.