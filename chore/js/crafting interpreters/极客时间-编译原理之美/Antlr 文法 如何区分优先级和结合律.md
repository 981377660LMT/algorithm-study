**在 ANTLR 文法中区分优先级和结合律**

在使用 ANTLR 编写语法时，正确设置运算符的优先级（Precedence）和结合律（Associativity）对于解析表达式至关重要。以下是如何在 ANTLR 文法中实现这一点的方法。

---

### **1. 运算符优先级**

运算符优先级决定了在没有括号指定的情况下，哪个运算符先被解析。例如，乘法 (`*`) 通常优先于加法 (`+`)。

### **2. 结合律**

结合律决定了相同优先级的运算符是从左到右解析（左结合）还是从右到左解析（右结合）。例如，减法 (`-`) 通常是左结合的，而赋值运算符 (`=`) 通常是右结合的。

---

### **在 ANTLR 中实现优先级和结合律**

ANTLR 通过规则的顺序和递归方式来处理运算符的优先级和结合律。以下是一个示例语法，展示了如何区分不同优先级和结合律的运算符。

#### **示例：简单的算术表达式语法**

```antlr
grammar Expr;

// 入口规则
expr: assignment;

// 赋值表达式（右结合）
assignment
    : IDENTIFIER '=' assignment       # Assign
    | addition                        # ToAddition
    ;

// 加法和减法（左结合， 内部处理了左递归）
addition
    : addition ('+' | '-') multiplication   # AddSub
    | multiplication                        # ToMultiplication
    ;

// 乘法和除法（左结合）
multiplication
    : multiplication ('*' | '/') atom       # MulDiv
    | atom                                   # ToAtom
    ;

// 原子（标识符、数字或括号表达式）
atom
    : NUMBER            # Number
    | IDENTIFIER        # Identifier
    | '(' expr ')'      # Parens
    ;

// 词法规则
IDENTIFIER: [a-zA-Z_][a-zA-Z_0-9]*;
NUMBER: [0-9]+;
WS: [ \t\r\n]+ -> skip;
```

---

### **解析过程说明**

1. **优先级设置**

   - **最高优先级**：`atom` 规则，处理基本的数字、标识符和括号表达式。
   - **中等优先级**：`multiplication` 规则，处理乘法 (`*`) 和除法 (`/`)。
   - **较低优先级**：`addition` 规则，处理加法 (`+`) 和减法 (`-`)。
   - **最低优先级**：`assignment` 规则，处理赋值 (`=`)。

   通过这种递归的方式，ANTLR 会先解析 `atom`，然后是 `multiplication`，再是 `addition`，最后是 `assignment`，从而确保运算符的优先级正确。

2. **结合律设置**

   - **左结合**：`addition` 和 `multiplication` 规则通过左递归的方式实现左结合。例如，`a - b - c` 会被解析为 `(a - b) - c`。
   - **右结合**：`assignment` 规则通过右递归的方式实现右结合。例如，`a = b = c` 会被解析为 `a = (b = c)`。

---

### **使用 Operator Precedence 格式**

除了递归规则，ANTLR 还支持使用 **operator precedence** 格式，通过直接在规则中指定优先级和结合律。这种方式更简洁，适用于较为复杂的表达式语法。

#### **示例：使用 Precedence 和 Associativity**

```antlr
grammar Expr;

expr
    : expr '^' expr       # Power        // 右结合
    | expr '*' expr       # MulDiv        // 左结合
    | expr '/' expr       # MulDiv        // 左结合
    | expr '+' expr       # AddSub        // 左结合
    | expr '-' expr       # AddSub        // 左结合
    | '(' expr ')'        # Parens
    | NUMBER              # Number
    ;

NUMBER: [0-9]+;
WS: [ \t\r\n]+ -> skip;
```

**注意**：上述方法在处理复杂的优先级和结合律时可能导致歧义，建议使用递归规则来明确指定优先级和结合律。

---

### **总结**

- **优先级**：通过规则的递归层次来控制运算符的优先级，先解析高优先级的规则。
- **结合律**：通过左递归（左结合）或右递归（右结合）来控制相同优先级运算符的解析顺序。

使用递归规则方式能够更清晰地定义优先级和结合律，避免语法歧义，适用于构建可靠的解析器。

---

希望以上内容能帮助你在 ANTLR 文法中正确区分和实现运算符的优先级与结合律。
