package gen

import (
	"fmt"
	"github.com/bantl23/tiny/symtbl"
	"github.com/bantl23/tiny/syntree"
	"os"
)

var PC int = 7
var MP int = 6
var GP int = 5
var AC int = 0
var AC1 int = 1
var TMP int = 0

var HighEmitLoc int = 0
var EmitLoc int = 0

type Gen struct {
	filename string
	file     *os.File
}

func (g *Gen) emitComment(c string) {
	fmt.Printf("* %s", c)
}

func (g *Gen) emitRO(op string, r int, s int, t int, c string) {
	buffer := fmt.Sprintf("%3d:  %5s  %d,%d,%d \n", EmitLoc, op, r, s, t)
	EmitLoc = EmitLoc + 1
	if HighEmitLoc < EmitLoc {
		HighEmitLoc = EmitLoc
	}
	fmt.Printf("%s", buffer)
	g.file.WriteString(buffer)
}

func (g *Gen) emitRM(op string, r int, d int, s int, c string) {
	buffer := fmt.Sprintf("%3d:  %5s  %d,%d(%d) \n", EmitLoc, op, r, d, s)
	EmitLoc = EmitLoc + 1
	if HighEmitLoc < EmitLoc {
		HighEmitLoc = EmitLoc
	}
	fmt.Printf("%s", buffer)
	g.file.WriteString(buffer)
}

func (g *Gen) emitSkip(howMany int) int {
	i := EmitLoc
	EmitLoc = EmitLoc + howMany
	if HighEmitLoc < EmitLoc {
		HighEmitLoc = EmitLoc
	}
	return i
}

func (g *Gen) emitBackup(loc int) {
	if loc > HighEmitLoc {
		g.emitComment("BUG in emitBackup")
	}
	EmitLoc = loc
}

func (g *Gen) emitRestore() {
	EmitLoc = HighEmitLoc
}

func (g *Gen) emitRM_Abs(op string, r int, a int, c string) {
	buffer := fmt.Sprintf("%3d:  %5s  %d,%d(%d) \n", EmitLoc, op, r, a-(EmitLoc+1), PC)
	EmitLoc = EmitLoc + 1
	if HighEmitLoc < EmitLoc {
		HighEmitLoc = EmitLoc
	}
	g.file.WriteString(buffer)
}

func (g *Gen) cGen(node *syntree.Node, table *symtbl.SymTbl) {
	if node != nil {
		switch node.NodeKind {
		case syntree.STATEMENT_KIND:
			g.getStmt(node, table)
		case syntree.EXPRESSION_KIND:
			g.getExp(node, table)
		}
		g.cGen(node.Sibling, table)
	}
}

func (g *Gen) getStmt(node *syntree.Node, table *symtbl.SymTbl) {
	var p1 *syntree.Node
	var p2 *syntree.Node
	var p3 *syntree.Node
	var savedLoc1 int
	var savedLoc2 int
	var currLoc int
	var loc int

	switch node.StmtKind {
	case syntree.IF_KIND:
		g.emitComment("-> if")
		p1 = node.Children[0]
		p2 = node.Children[1]
		if len(node.Children) == 3 {
			p3 = node.Children[2]
		}

		g.cGen(p1, table)
		savedLoc1 = g.emitSkip(1)
		g.emitComment("if: jump to else belongs here")

		g.cGen(p2, table)
		savedLoc2 = g.emitSkip(1)
		g.emitComment("if: jump to end belongs here")

		currLoc = g.emitSkip(0)
		g.emitBackup(savedLoc1)
		g.emitRM_Abs("JEQ", AC, currLoc, "if: jmp to else")
		g.emitRestore()

		g.cGen(p3, table)
		currLoc = g.emitSkip(0)
		g.emitBackup(savedLoc2)
		g.emitRM_Abs("LDA", PC, currLoc, "jmp to end")
		g.emitRestore()
	case syntree.REPEAT_KIND:
		g.emitComment("-> repeat")
		p1 = node.Children[0]
		p2 = node.Children[1]
		savedLoc1 = g.emitSkip(0)

		g.emitComment("repeat: jump after body comes back here")
		g.cGen(p1, table)
		g.cGen(p2, table)

		g.emitRM_Abs("JEQ", AC, savedLoc1, "repeat: jmp back to body")

	case syntree.ASSIGN_KIND:
		g.emitComment("-> assign")
		g.cGen(node.Children[0], table)

		loc = table.Obtain(node.Name)
		g.emitRM("ST", AC, loc, GP, "assign: store value")
	case syntree.READ_KIND:
		g.emitRO("IN", AC, 0, 0, "read integer value")
		loc = table.Obtain(node.Name)
		g.emitRM("ST", AC, loc, GP, "read: store value")
	case syntree.WRITE_KIND:
		g.cGen(node.Children[0], table)
		g.emitRO("OUT", AC, 0, 0, "write ac")
	}
}

func (g *Gen) getExp(node *syntree.Node, table *symtbl.SymTbl) {
	var p1 *syntree.Node
	var p2 *syntree.Node
	var loc int

	switch node.ExpKind {
	case syntree.CONST_KIND:
		g.emitComment("-> Const")
		g.emitRM("LDC", AC, node.Value, 0, "load const")
	case syntree.ID_KIND:
		g.emitComment("-> Id")

		loc = table.Obtain(node.Name)
		g.emitRM("LD", AC, loc, GP, "load id value")

	case syntree.OP_KIND:
		g.emitComment("-> Op")
		p1 = node.Children[0]
		p2 = node.Children[1]

		g.cGen(p1, table)
		g.emitRM("ST", AC, TMP, MP, "op: push left")
		TMP = TMP - 1

		g.cGen(p2, table)
		TMP = TMP + 1
		g.emitRM("LD", AC1, TMP, MP, "op: load left")

		switch node.TokenType {
		case syntree.PLUS:
			g.emitRO("ADD", AC, AC1, AC, "op =")
		case syntree.MINUS:
			g.emitRO("SUB", AC, AC1, AC, "op -")
		case syntree.TIMES:
			g.emitRO("MUL", AC, AC1, AC, "op *")
		case syntree.OVER:
			g.emitRO("DIV", AC, AC1, AC, "op /")
		case syntree.LT:
			g.emitRO("SUB", AC, AC1, AC, "op <")
			g.emitRM("JLT", AC, 2, PC, "br if true")
			g.emitRM("LDC", AC, 0, AC, "false case")
			g.emitRM("LDA", PC, 1, PC, "unconditional jump")
			g.emitRM("LDC", AC, 1, AC, "true case")
		case syntree.EQ:
			g.emitRO("SUB", AC, AC1, AC, "op ==")
			g.emitRM("JEQ", AC, 2, PC, "br if true")
			g.emitRM("LDC", AC, 0, AC, "false case")
			g.emitRM("LDA", PC, 1, PC, "unconditional jump")
			g.emitRM("LDC", AC, 1, AC, "true case")
		default:
			g.emitComment("BUG: unknown operator")
		}
	}
}

func (g *Gen) Generate(node *syntree.Node, table *symtbl.SymTbl, file string) {
	g.filename = file
	g.file, _ = os.Create(g.filename)

	g.emitRM("LD", MP, 0, AC, "load maxaddress from location 0")
	g.emitRM("ST", AC, 0, AC, "clear location 0")
	g.cGen(node, table)

	g.emitRO("HALT", 0, 0, 0, "")
}
