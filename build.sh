#!/usr/bin/env bash

go get github.com/blynn/nex
nex tiny.nex
go tool yacc tiny.y
go build
