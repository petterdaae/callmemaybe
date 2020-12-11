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

<calculation> := <val> (<bop> <val>)*

<val>         := "(" <calulation> ")"
<val>         := <num>
<val>         := <identifier>
<val>         := <bool>
<val>         := <uop> <num>
<val>         := <char>
<val>         := <list>
<val>         := <list> ? <num>

<bop>         := "+" | "*" | "<" | ">" | "==" | "-" | "/" | "%"
<uop>         := "-"
<num>         := [1-9][0-9]+
<bool>        := "true" | "false"
<char>        := [^\'] | "\" "'" | "\" "\"

<function>    := <argList> "=>" <type>? "{" <seq> "}"
<argList>     := "<" <identifier>? ">" | "<" (<recurse>,)? (<arg> ",")* <arg> ">"
<arg>         := <identifier> <type>

<list>        := "[" (<expr> ",") <expr> "]" : <type> | "[" "]" : <type>

<identifier>  := (letter|_)[letter|[0-9]|_]*

<type>        := "int" | "empty" | "bool" | "char"
```

### TODO
- Lists and strings
- Loop
- Make it possible to return functions and pass functions as arguments
- Structs
- Input
- Handle division-by-zero and out-of-bounds errors
- Make grammar more formal

### Resources
- The structure of the tokenizer and parser is inspired by this blog post: https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
- This site was very helpful when constructing the grammar: http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm.
- Nice x86 assembly reference: http://www.cs.virginia.edu/~evans/cs216/guides/x86.html.
