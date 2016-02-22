%{
package main

import (
  "fmt"
  "strconv"
  "gitlab.com/bantl23/python/syntree"
)

var root *syntree.Node
var savedName string
%}

%union {
  node *syntree.Node
  str string
}

%token <str> IF THEN ELSE END REPEAT UNTIL READ WRITE ASSIGN EQ LT PLUS MINUS TIMES OVER LPAREN RPAREN SEMI ID NUM ERROR
%type <node> program stmt_seq stmt if_stmt repeat_stmt assign_stmt read_stmt write_stmt error exp simple_exp term factor

%%

program     : stmt_seq                  {
                                          fmt.Printf("program0: %+v\n", $1)
                                          root = $1
                                        }

stmt_seq    : stmt_seq SEMI stmt        {
                                          fmt.Printf("stmt_seq0: %+v %+v %+v\n", $1, $2, $3)
                                          t := $1
                                          if (t != nil) {
                                            for (t.Sibling != nil) {
                                              t = t.Sibling
                                            }
                                            t.Sibling = $3
                                            $$ = $1
                                          } else {
                                            $$ = $3
                                          }
                                        }
            | stmt                      {
                                          fmt.Printf("stmt_seq1: %+v\n", $1)
                                          $$ = $1
                                        }
            ;

stmt        : if_stmt                   {
                                          fmt.Printf("stmt0: %+v\n", $1)
                                          $$ = $1
                                        }
            | repeat_stmt               {
                                          fmt.Printf("stmt1: %+v\n", $1)
                                          $$ = $1
                                        }
            | assign_stmt               {
                                          fmt.Printf("stmt2: %+v\n", $1)
                                          $$ = $1
                                        }
            | read_stmt                 {
                                          fmt.Printf("stmt3: %+v\n", $1)
                                          $$ = $1
                                        }
            | write_stmt                {
                                          fmt.Printf("stmt4: %+v\n", $1)
                                          $$ = $1
                                        }
            | error                     {
                                          fmt.Printf("stmt5: %+v\n", $1)
                                          $$ = nil
                                        }
            ;

if_stmt     : IF exp THEN stmt_seq END  {
                                          fmt.Printf("if_stmt0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.IF_KIND
                                          $$.Children = append($$.Children, $2)
                                          $$.Children = append($$.Children, $4)
                                        }
            | IF exp THEN stmt_seq ELSE stmt_seq END
                                        {
                                          fmt.Printf("if_stmt1: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $3, $5, $6, $7)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.IF_KIND
                                          $$.Children = append($$.Children, $2)
                                          $$.Children = append($$.Children, $4)
                                          $$.Children = append($$.Children, $6)
                                        }
            ;

repeat_stmt : REPEAT stmt_seq UNTIL exp {
                                          fmt.Printf("repeat_stmt0: %+v %+v %+v %+v\n", $1, $2, $3, $4)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.REPEAT_KIND
                                          $$.Children = append($$.Children, $2)
                                          $$.Children = append($$.Children, $4)
                                        }
            ;

assign_stmt : ID                        {
                                          fmt.Printf("assign_stmt0a: %+v\n", currText(yylex))
                                          savedName = currText(yylex)
                                        }
              ASSIGN exp
                                        {
                                          fmt.Printf("assign_stmt0b: %+v %+v %+v %+v $+v\n", $<node>1, $<node>2, $<node>3, $<node>4, currText(yylex))
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.ASSIGN_KIND
                                          $$.Children = append($$.Children, $<node>4)
                                          $$.Name = savedName
                                          $$.LineNumber = currLine(yylex)
                                        }
            ;

read_stmt   : READ ID                   {
                                          fmt.Printf("read_stmt0: %+v %+v\n", $1, $2)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.READ_KIND
                                          $$.Name = currText(yylex)
                                        }

write_stmt  : WRITE exp                 {
                                          fmt.Printf("write_stmt0: %+v %+v\n", $1, $2)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.WRITE_KIND
                                          $$.Children = append($$.Children, $2)
                                        }

exp         : simple_exp LT simple_exp  {
                                          fmt.Printf("exp0: %+v %+v %+v\n", $1, $2, $3)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.LT
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | simple_exp EQ simple_exp  {
                                          fmt.Printf("exp1: %+v %+v %+v\n", $1, $2, $3)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.EQ
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | simple_exp                {
                                          fmt.Printf("exp2: %+v\n", $1)
                                          $$ = $1
                                        }

simple_exp  : simple_exp PLUS term      {
                                          fmt.Printf("simple_exp0: %+v %+v %+v\n", $1, $2, $3)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.PLUS
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | simple_exp MINUS term     {
                                          fmt.Printf("simple_exp1: %+v %+v %+v\n", $1, $2, $3)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.MINUS
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | term                      {
                                          fmt.Printf("simple_exp2: %+v\n", $1)
                                          $$ = $1
                                        }
            ;

term        : term TIMES factor         {
                                          fmt.Printf("term0: %+v %+v %+v\n", $1, $2, $3)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.TIMES
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | term OVER factor          {
                                          fmt.Printf("term1: %+v %+v %+v\n", $1, $2, $3)
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.OVER
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | factor                    {
                                          fmt.Printf("term2: %+v\n", $1)
                                          $$ = $1
                                        }
            ;

factor      : LPAREN exp RPAREN         {
                                          fmt.Printf("factor0: %+v %+v %+v\n", $1, $2, $3)
                                          $$ = $2
                                        }
            | NUM                       {
                                          fmt.Printf("factor1: %+v %+v\n", $1, currText(yylex))
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.CONST_KIND
                                          $$.Value, _ = strconv.Atoi(currText(yylex))
                                        }
            | ID                        {
                                          fmt.Printf("factor2: %+v %+v\n", $1, currText(yylex))
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.ID_KIND
                                          $$.Name = currText(yylex)
                                        }
            | error                     {
                                          fmt.Printf("factor3: %+v\n", $1)
                                          $$ = nil
                                        }
            ;

%%
func currText(y yyLexer) string {
  if len(y.(*Lexer).stack) > 0 {
    return y.(*Lexer).stack[0].s
  }
  return ""
}
func currLine(y yyLexer) int {
  if len(y.(*Lexer).stack) > 0 {
    return y.(*Lexer).stack[0].line
  }
  return 0
}
func currCol(y yyLexer) int {
  if len(y.(*Lexer).stack) > 0 {
    return y.(*Lexer).stack[0].column
  }
  return 0
}
