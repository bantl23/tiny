/if/         { println("if"); return IF }
/then/       { println("then"); return THEN }
/else/       { println("else"); return ELSE }
/end/        { println("end"); return END }
/repeat/     { println("repeat"); return REPEAT }
/until/      { println("until"); return UNTIL }
/read/       { println("read"); return READ }
/write/      { println("write"); return WRITE }
/:=/         { println(":="); return ASSIGN }
/=/          { println("="); return EQ }
/</          { println("<"); return LT }
/\+/         { println("+"); return PLUS }
/-/          { println("-"); return MINUS }
/\*/         { println("*"); return TIMES }
/\//         { println("/"); return OVER }
/\(/         { println("("); return LPAREN }
/\)/         { println(")"); return RPAREN }
/;/          { println(";"); return SEMI }
/[0-9]+/     { println(yylex.Text()); return NUM }
/[a-zA-Z]+/  { println(yylex.Text()); return ID }
/[ \t\n]+/   { /* white space */ }
/{[^\{\}]*}/ { /* comments */ }
/./          { return ERROR }
//
package main
