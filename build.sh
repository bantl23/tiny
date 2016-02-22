#!/usr/bin/env bash

go get github.com/blynn/nex
nex tiny.l
go tool yacc tiny.y
go build
