# callmemaybe
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

<expr>        := <calculation>
<expr>        := <function>
<expr>        := <call>
<expr>        := <list>
<expr>        := <string>

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
<string>      := "\"" <stringChar>* "\""
<rawChar>     := "\"" | "\\" | [^"\]

<function>    := <argList> <type>? "{" <seq> "}"
<argList>     := "|" "me"? "|" | "|" "me,"? (<arg> ",")* <arg> "|"
<arg>         := <identifier> <type>

<list>        := "<" <type> "," <num> ">" "[" (<expr> ",") <expr> "]"
<list>        := "<" <type> ">" "[" "]"

<identifier>  := (letter|_)[letter|[0-9]|_]*

<type>        := <rawType> | "list" "<" <type> "," <num> "> | "func" | "func" "<" <type>+ ">"
<rawType>     := "int" | "char" | "bool"
```

### TODO
- 14.12 Loop
- 14.12 Write tests
- 15.12 Improve IO
- 15.12 Free heap allocated memory when out of scope (?)
- 15.12 Handle division-by-zero and out-of-bounds errors
- 16.12 Structs
- 16.12 Module system
- 16.12 Write tests
