适合搭建解析器原型

- `lib.ts`:解析库
- `grammar.ts`:语法定义
- `parseGen.ts`:parse 生成器
- `grammar.grammar`:语法定义文件
  需要注意的是，

  - lexical(词法)使用`:`来分割 head 和 body, 使用 `RegExp` 而非自然语言描述词，且正则表达式必须以`^`开头
  - syntax(语法)使用`->`来分割 head 和 body

  ````
  Prologue : /^```.*```/s ;

  Regex : /^[/][^].*[/](?=\s*[;\.])/ ;

  MapCode : /^\.map\(.*?\)(?=\s*;)/s ;

  Terminal : /^[A-Z][A-Za-z_]*/ ;

  NonTerminal : /^[a-z][A-Za-z_]*/ ;

  Literal : /^"[^"]*"/ ;

  primary -> Terminal | NonTerminal | Literal | "(" choice ")" ;

  qualified -> primary "?" | primary "*" | primary "+" | primary ;

  sequence -> qualified+ ;

  choice -> sequence ("|" sequence)* ;

  syntax -> NonTerminal "->" choice MapCode? ";" ;

  lexical -> Terminal ":" Regex MapCode? ";" ;

  grammar -> Prologue? (syntax | lexical)* ;
  ````

grammar.ts 基于 lib.ts 工作，parseGen.ts 使用 grammar.ts ，读取语法定义文件(后缀名一般是`.grammar`)，生成解析器

---

表达式解析

https://craftinginterpreters.com/appendix-i.html
https://craftinginterpreters.com/representing-code.html#enhancing-our-notation
https://github.com/981377660LMT/parserGen

---

- [parser combinators 的实现](https://qszhu.github.io/2021/08/22/parser-combinators.html)
- [parser combinators 的使用](https://qszhu.github.io/2021/09/07/parsing-misc.html)
