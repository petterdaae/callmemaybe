# lang
A simple compiler.

### The current grammar that I am trying to implement
```
<seq>         := <stmt>*

<stmt>        := <assign>
<stmt>        := <println>
<stmt>        := <call>
<stmt>        := <return>
<stmt>        := <if>

<assign>      := <identifier> "=" <exp>
<println>     := "println" <exp>
<return>      := "return" <exp>
<loop>        := "loop" <expr> "{" <seq> "}"
<if>          := "if" <expr> "{" <seq> "}" 
<call>        := "call" <identifier> [ "with" (<expr> ",")* <expr> ]
<append>      := "append" <expr> "to" <list>

<expr>        := <calculation>
<expr>        := <function>
<expr>        := <call>
<expr>        := <list>

<calculation> := <val> (<bop> <val>)*

<val>         := "(" <calulation> ")"
<val>         := <num>
<val>         := <identifier>
<val>         := <bool>
<val>         := <uop> <num>
<val>         := <char>
<val>         := "get" <exp> "from" <expr>

<bop>         := "+" | "*" | "<" | ">" | "==" | "-" | "/" | "%"
<uop>         := "-"
<num>         := [0-9]+
<bool>        := "true" | "false"
<char>        := [^\'] | "\" "'" | "\" "\"

<function>    := <argList> "=>" <type>? "{" <seq> "}"
<argList>     := "<" <identifier>? ">" | "<" (<recurse>,)? (<arg> ",")* <arg> ">"
<arg>         := <identifier> <type>

<list>        := "[" (<expr> ",") <expr> "]" ":" <type> ":" <num> | "[" "]" ":" <type> ":" <num>

<identifier>  := (letter|_)[letter|[0-9]|_]*

<type>        := "int" | "empty" | "bool" | "char"
```

### TODO
- Comparison for characters
- Expand type signatures to handle lists
- Strings (syntactic sugar for lists of characters)
- Print lists
- Input
- Loop
- Syntax highlighting  
- Make it possible to return functions and pass functions as arguments
- Structs
- Handle division-by-zero and out-of-bounds errors
- Free heap allocated memory when out of scope (?)
- Module system (github?)
