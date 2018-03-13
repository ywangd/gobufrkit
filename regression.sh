#!/bin/bash

# Build the project
go build
if [[ "$?" != "0" ]]; then
    exit 1
fi

function doTest {
    diff -qs <(./gobufrkit d -jx _testdata/$1.bufr) _regression/$1.bufr.json
    if [[ "$?" != "0" ]]; then
        echo ERROR
        exit 1
    fi
}

doTest 207003
doTest ISMD01_OKPR
doTest IUSK73_AMMC_040000
doTest IUSK73_AMMC_182300
doTest amv2_87
doTest asr3_190
doTest b002_95
doTest b005_89
doTest contrived
doTest g2nd_208
doTest jaso_214
doTest mpco_217
doTest profiler_european
doTest rado_250
doTest uegabe
