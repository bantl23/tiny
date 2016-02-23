Tiny Compiler
=============

## Overview

### Files

* tiny.nex
  * Lex Definitions
* tiny.y
  * Yacc Definitions
* main.go
  * main rountine

### Directories

* syntree
  * Syntax Tree Implementation
    * Types
    * Tree traversal
* symtbl
  * Symbol Tree Implementation
    * Table implemented as a HashMap
    * Checks semantics
* gen
  * Code Generation Implementation
    * Emits code
    * Generates code

## TODO

* Cleanup code
* Handle echo source code flag
* Make trace optional
