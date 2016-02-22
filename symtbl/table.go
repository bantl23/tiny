package symtbl

import (
	"fmt"
	"gitlab.com/bantl23/python/syntree"
)

type Bucket struct {
	Lines  []int
	MemLoc int
}

type SymTbl map[string]*Bucket

var Location int = 0

func (s *SymTbl) Insert(name string, line int, memLoc int) {
	table := *s
	fmt.Printf("%+v %+v %+v\n", name, line, memLoc)
	_, ok := table[name]
	if ok == true {
		table[name].Lines = append(table[name].Lines, line)
	} else {
		table[name] = new(Bucket)
		table[name].Lines = append(table[name].Lines, line)
		table[name].MemLoc = memLoc
	}
}

func (s *SymTbl) Obtain(name string) int {
	table := *s
	_, ok := table[name]
	if ok == true {
		return table[name].MemLoc
	}
	return -1
}

func (t *SymTbl) InsertNode(node *syntree.Node) {
	table := *t
	switch node.NodeKind {
	case syntree.STATEMENT_KIND:
		switch node.StmtKind {
		case syntree.ASSIGN_KIND:
			if table.Obtain(node.Name) == -1 {
				table.Insert(node.Name, node.LineNumber, Location)
				Location = Location + 1
			} else {
				table.Insert(node.Name, node.LineNumber, 0)
			}
		case syntree.READ_KIND:
			if table.Obtain(node.Name) == -1 {
				table.Insert(node.Name, node.LineNumber, Location)
				Location = Location + 1
			} else {
				table.Insert(node.Name, node.LineNumber, 0)
			}
		}
	case syntree.EXPRESSION_KIND:
		switch node.ExpKind {
		case syntree.ID_KIND:
			if table.Obtain(node.Name) == -1 {
				table.Insert(node.Name, node.LineNumber, Location)
				Location = Location + 1
			} else {
				table.Insert(node.Name, node.LineNumber, 0)
			}
		}
	}
}

func (t *SymTbl) BuildTable(node *syntree.Node) {
	syntree.Traverse(node, t.InsertNode, syntree.Nothing)
}

func (t *SymTbl) PrintTable() {
	table := *t
	for i, e := range table {
		fmt.Printf("%+v: %+v %+v\n", i, e.MemLoc, e.Lines)
	}
}

func (t *SymTbl) CheckNode(node *syntree.Node) {
	switch node.NodeKind {
	case syntree.EXPRESSION_KIND:
		switch node.ExpKind {
		case syntree.OP_KIND:
			if node.Children[0].ExpType != syntree.INTEGER_TYPE || node.Children[0].ExpType != syntree.INTEGER_TYPE {
				fmt.Println("error op applied to non integer")
			}
			if node.TokenType == syntree.EQ || node.TokenType == syntree.LT {
				node.ExpType = syntree.BOOLEAN_TYPE
			} else {
				node.ExpType = syntree.INTEGER_TYPE
			}
		case syntree.CONST_KIND:
			node.ExpType = syntree.INTEGER_TYPE
		case syntree.ID_KIND:
			node.ExpType = syntree.INTEGER_TYPE
		}
	case syntree.STATEMENT_KIND:
		switch node.StmtKind {
		case syntree.IF_KIND:
			if node.Children[0].ExpType == syntree.INTEGER_TYPE {
				fmt.Println("error if test non-boolean")
			}
		case syntree.ASSIGN_KIND:
			if node.Children[0].ExpType != syntree.INTEGER_TYPE {
				fmt.Println("error assignment non-integer")
			}
		case syntree.WRITE_KIND:
			if node.Children[0].ExpType != syntree.INTEGER_TYPE {
				fmt.Println("error write non-integer")
			}
		case syntree.REPEAT_KIND:
			if node.Children[0].ExpType == syntree.INTEGER_TYPE {
				fmt.Println("error repeat non-boolean")
			}
		}
	}
}

func (t *SymTbl) CheckTable(node *syntree.Node) {
	syntree.Traverse(node, syntree.Nothing, t.CheckNode)
}
