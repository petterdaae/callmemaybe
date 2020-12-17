# Grammar

```
<seq>             := <stmt>*

<stmt>            := <assign> | <println> | <return> | <if> | <loop> | <structType> | <update>

<assign>          := <identifier> "=" <exp>
<println>         := "println" <exp>
<return>          := "return" <exp>
<if>              := "if" <expr> "{" <seq> "}" 
<loop>            := "loop" <expr> "{" <seq> "}"
<structType>      := "struct" <identifier> "{" (<identifier> <type>)* "}"
<update>          := <reference> "=" <exp>

<exp>             := <val> (<bop> <val>)
<val>             := <num> | <bool> | <char> | <function> | <call> | <list> | <string> | 
                     <structValue> | <reference> | "length" "(" <exp> ")" | <uop> <exp> |
                     <identifier> | "(" <exp> ")"
<bop>             := "+" | "*" | "<" | ">" | "==" | "-" | "/" | "%" | "!="
<uop>             := "-"
                     
<num>             := [0-9]+
<bool>            := "true" | "false"
<char>            := TODO
<function>        := "|" (("me"|<identifier><type>)(","<identifier><type>)*)? "|" <type>? "{" <seq> "}"
<call>            := "#"<exp>("("(<exp> ",")*<exp>")")?
<list>            := "<"<type>","<num>">" "[" (<exp> (","<exp>)*)? "]"
<string>          := TODO
<structValue>     := "@" <identifier> "{" (<identifier> ":" <type>)* "}"
<length>          := "length" "(" <exp> ")"

<reference>       := "?" <exp> ( "." <identifier> | "[" <exp> "]" )+
<identifier>      := TODO

<type>            := "@" <identifier>
<type>            := "int" | "char" | "bool" | "string" | "func"
<type>            := "list" "<" <type> ">" 
<type>            := "func" "<" <type>+ ">"
```
