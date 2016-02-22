package syntree

import (
	"fmt"
)

type NodeKind int

const (
	STATEMENT_KIND = iota
	EXPRESSION_KIND
)

type StatementKind int

const (
	IF_KIND = iota
	REPEAT_KIND
	ASSIGN_KIND
	READ_KIND
	WRITE_KIND
)

type ExpressionKind int

const (
	OP_KIND = iota
	CONST_KIND
	ID_KIND
)

type ExpressionType int

const (
	VOID_TYPE = iota
	INTEGER_TYPE
	BOOLEAN_TYPE
)

type TokenType int

const (
	ENDFILE = iota
	ERROR
	IF
	THEN
	ELSE
	END
	REPEAT
	UNTIL
	READ
	WRITE
	ID
	NUM
	ASSIGN
	EQ
	LT
	PLUS
	MINUS
	TIMES
	OVER
	LPAREN
	RPAREN
	SEMI
)

type Node struct {
	NodeKind   NodeKind
	StmtKind   StatementKind
	ExpKind    ExpressionKind
	ExpType    ExpressionType
	TokenType  TokenType
	Name       string
	Value      int
	LineNumber int
	Sibling    *Node
	Children   []*Node
}

func NewNode() *Node {
	n := new(Node)
	return n
}

func Print(node *Node, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
	for node != nil {
		if node.NodeKind == STATEMENT_KIND {
			switch node.StmtKind {
			case IF_KIND:
				fmt.Println("if")
			case REPEAT_KIND:
				fmt.Println("repeat")
			case ASSIGN_KIND:
				fmt.Println("assign:", node.Name)
			case READ_KIND:
				fmt.Println("read:", node.Name)
			case WRITE_KIND:
				fmt.Println("write")
			default:
				fmt.Println("unknown statement kind")
			}
		} else if node.NodeKind == EXPRESSION_KIND {
			switch node.ExpKind {
			case OP_KIND:
				fmt.Println("operator:", node.TokenType)
			case CONST_KIND:
				fmt.Println("constant:", node.Value)
			case ID_KIND:
				fmt.Println("id:", node.Name)
			default:
				fmt.Println("unknown expression kind")
			}
		} else {
			fmt.Println("unknown node kind")
		}
		for _, v := range node.Children {
			Print(v, indent+2)
		}
		node = node.Sibling
	}
}

type Proc func(class *Node)

func Traverse(node *Node, pre Proc, post Proc) {
	if node != nil {
		pre(node)
		for _, n := range node.Children {
			Traverse(n, pre, post)
		}
		post(node)
		Traverse(node.Sibling, pre, post)
	}
}

func Nothing(node *Node) {
	return
}
