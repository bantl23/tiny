#!/usr/bin/env bash

nex tiny.l
go tool yacc tiny.y
go build
