# callmemaybe
A simple compiler.

### The current grammar that I am trying to implement
```
<seq>             := <stmt>*

<stmt>            := <assign>
<stmt>            := <println>
<stmt>            := <call>
<stmt>            := <return>
<stmt>            := <if>
<stmt>            := <struct>

<assign>          := <identifier> "=" <exp>
<assignStruct>    := <identifier> "." <identifier> ( "." <identifier> )* "=" <exp>
<assignList>      := <identifier> "[" <exp> "]" ( "[" <exp> "]" )* "=" <exp>
<println>         := "println" <exp>
<return>          := "return" <exp>
<loop>            := "loop" <expr> "{" <seq> "}"
<if>              := "if" <expr> "{" <seq> "}" 

<call>            := "!" <expr> "(" (<expr> ",")* <expr> ] ")"

<expr>            := <calculation>
<expr>            := <function>
<expr>            := <call>
<expr>            := <list>
<expr>            := <string>
<expr>            := <structExp>

<calculation>     := <val> (<bop> <val>)*

<val>             := "(" <calulation> ")"
<val>             := <num>
<val>             := <identifier>
<val>             := <bool>
<val>             := <uop> <num>
<val>             := <char>
<val>             := "?" <expr> "[" <expr> "]" ( "[" <expr> "]" )*
<val>             := "?" <expr>  "." <identifier> ( "." <identifier> )*

<bop>             := "+" | "*" | "<" | ">" | "==" | "-" | "/" | "%"
<uop>             := "-"
<num>             := [0-9]+
<bool>            := "true" | "false"
<char>            := [^\'] | "\" "'" | "\" "\"
<string>          := "\"" <stringChar>* "\""
<rawChar>         := "\"" | "\\" | [^"\]

<function>        := <argList> <type>? "{" <seq> "}"
<argList>         := "|" "me"? "|" | "|" "me,"? (<arg> ",")* <arg> "|"
<arg>             := <identifier> <type>

<list>            := "<" <type> "," <num> ">" "[" (<expr> ",") <expr> "]"
<list>            := "<" <type> ">" "[" "]"

<identifier>      := (letter|_)[letter|[0-9]|_]*

<type>            := "@" <identifier>
<type>            := "int" | "char" | "bool" | "string"
<type>            := "list" "<" <type> ">" 
<type>            := "func" 
<type>            := "func" "<" <type>+ ">"

<struct>          := "struct" <identifier> "{" <structMember>* "}"
<structMember>    := <identifier> <type>

<structExp>       := "@" <identifier> "{" <structExpMember>* "}"
<structExpMember> := <identifier> ":" <exp>
```

### TODO
- 16.12 Module system
- 16.12 Input
- 16.12 Free heap allocated memory when out of scope (?)
- 16.12 Handle division-by-zero and out-of-bounds errors
- 16.12 Write tests
