# lang
A simple compiler.

### The current grammar that I am trying to implement
```
<seq>         := <stmt>*
<stmt>        := <assign> | <prinln> | <call> | <return> | <if>
<assign>      := <identifier> "=" <exp>
<println>     := "println" <exp>
<return>      := "return" <exp>
<loop>        := "loop" <expr> "{" <seq> "}"
<if>          := "if" <expr> "{" <seq> "}" 
<call>        := "call" <identifier> [ "with" (<expr> ",")* <expr> ]

<expr>        := <calculation> | <function> | <call>
<calculation> := <val> (<bop> <val>)*
<val>         := "(" <calulation> ")" | <num> | <identifier> | <bool> | <uop> <num>
<bop>         := "+" | "*" | "<" | ">" | "==" | "-" | "/" | "%"
<uop>         := "-"
<num>         := sequence of digits
<bool>        := "true" | "false"
<function>    := <argList> "=>" <type>? "{" <seq> "}"
<identifier>  := words consisting og letters, digits and underscores, 
                 starting with a letter or underscore

<argList>     := "<" <identifier>? ">" | "<" (<recurse>,)? (<arg> ",")* <arg> ">"
<arg>         := <identifier> <type>
<type>        := "int" | "empty" | "bool"
```

### TODO
- Negative numbers, subtraction, division and modulo
- Add an extra parameter to functions to make recursion possible
- Loop
- Refactor
- Characters, lists and strings
- IO
- Make it possible to return functions and pass functions as arguments
- Structs

### Resources
- The structure of the tokenizer and parser is inspired by this blog post: https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
- This site was very helpful when constructing the grammar: http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm.
- Nice x86 assembly reference: http://www.cs.virginia.edu/~evans/cs216/guides/x86.html.
