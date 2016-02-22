%{
package main

import (
  "strconv"
  "github.com/bantl23/tiny/syntree"
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
                                          root = $1
                                        }

stmt_seq    : stmt_seq SEMI stmt        {
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
                                          $$ = $1
                                        }
            ;

stmt        : if_stmt                   {
                                          $$ = $1
                                        }
            | repeat_stmt               {
                                          $$ = $1
                                        }
            | assign_stmt               {
                                          $$ = $1
                                        }
            | read_stmt                 {
                                          $$ = $1
                                        }
            | write_stmt                {
                                          $$ = $1
                                        }
            | error                     {
                                          $$ = nil
                                        }
            ;

if_stmt     : IF exp THEN stmt_seq END  {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.IF_KIND
                                          $$.Children = append($$.Children, $2)
                                          $$.Children = append($$.Children, $4)
                                        }
            | IF exp THEN stmt_seq ELSE stmt_seq END
                                        {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.IF_KIND
                                          $$.Children = append($$.Children, $2)
                                          $$.Children = append($$.Children, $4)
                                          $$.Children = append($$.Children, $6)
                                        }
            ;

repeat_stmt : REPEAT stmt_seq UNTIL exp {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.REPEAT_KIND
                                          $$.Children = append($$.Children, $2)
                                          $$.Children = append($$.Children, $4)
                                        }
            ;

assign_stmt : ID                        {
                                          savedName = currText(yylex)
                                        }
              ASSIGN exp
                                        {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.ASSIGN_KIND
                                          $$.Children = append($$.Children, $<node>4)
                                          $$.Name = savedName
                                          $$.LineNumber = currLine(yylex)
                                        }
            ;

read_stmt   : READ ID                   {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.READ_KIND
                                          $$.Name = currText(yylex)
                                        }

write_stmt  : WRITE exp                 {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.STATEMENT_KIND
                                          $$.StmtKind = syntree.WRITE_KIND
                                          $$.Children = append($$.Children, $2)
                                        }

exp         : simple_exp LT simple_exp  {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.LT
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | simple_exp EQ simple_exp  {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.EQ
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | simple_exp                {
                                          $$ = $1
                                        }

simple_exp  : simple_exp PLUS term      {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.PLUS
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | simple_exp MINUS term     {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.MINUS
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | term                      {
                                          $$ = $1
                                        }
            ;

term        : term TIMES factor         {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.TIMES
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | term OVER factor          {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.OP_KIND
                                          $$.TokenType = syntree.OVER
                                          $$.Children = append($$.Children, $1)
                                          $$.Children = append($$.Children, $3)
                                        }
            | factor                    {
                                          $$ = $1
                                        }
            ;

factor      : LPAREN exp RPAREN         {
                                          $$ = $2
                                        }
            | NUM                       {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.CONST_KIND
                                          $$.Value, _ = strconv.Atoi(currText(yylex))
                                        }
            | ID                        {
                                          $$ = syntree.NewNode()
                                          $$.NodeKind = syntree.EXPRESSION_KIND
                                          $$.ExpKind = syntree.ID_KIND
                                          $$.Name = currText(yylex)
                                        }
            | error                     {
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
