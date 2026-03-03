## Context-Free Grammar Rules
Describes the rules of golox's CFG used to parse tokens. 

expression → literal | unary | binary | grouping ;

literal → NUMBER | STRING | "true" | "false" | "none" ;
grouping → "(" expression ")" ;
unary → ("-" | "!") expression ;
binary → expression operator expression ;
operator → "==" | "!=" | "<" | "<=" | ">" | ">=" | "+" | "-" | "*" | "/" ;
