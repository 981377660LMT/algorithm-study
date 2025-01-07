下面是一份**详细分析**该代码的讲解。它实现了一个将 Go 的正则语法（通过 `regexp/syntax` 包解析）编译成**确定性有限自动机 (DFA)** 的过程，并提供了一个 `Regexp` 结构来匹配输入字符串。它也使用了一些与 UTF-8 编码、字节范围处理相关的技巧。

阅读本解读前，需要了解以下概念：

- **正则解析 (regexp/syntax) 的 AST**：Go 内置的 `regexp/syntax` 库可将正则表达式解析为一个语法树（`syntax.Regexp`）；
- **NFA/DFA 编译**：把正则转换成指令集 (类似于“Thompson 的 NFA”，再到 DFA 的过程)；
- **UTF-8 字符拆分**：用“byte ranges”实现对 Unicode 字符（可能占多字节）进行匹配；
- **稀疏集合 (sparseSet)**：高效管理活跃指令索引集合；
- **有限状态机**：`dfa.states` 数组，每个 `state` 保存 “下一跳表 + 是否 match”；

下面我们会按**模块**进行详细解读：

---

## 1. `Regexp` 结构

```go
type Regexp struct {
    orig string
    dfa  *dfa
}
```

- `Regexp` 是对外公开的结构，保存原始正则表达式字符串 (`orig`) 和编译好的 DFA (`dfa`)。
- 构造方式：
  1. `NewRegexp(expr string) (*Regexp, error)`：创建一个正则自动机，默认限制**编译后**大小上限 (10 MB)；
  2. `NewRegexpWithLimit(expr, size)`：自定义上限。
  3. 内部会调用 `syntax.Parse(expr, syntax.Perl)` 获取正则语法树，然后 `compile` => `dfaBuilder.build()` 得到 `dfa`。

### 1.1 `Start()/IsMatch()/Accept()`

这是对外实现的**状态机接口** (符合 vellum.Automaton 规范)：

- `Start() int`：起始状态为 `1`（注意代码里 state 0 是无效状态）；
- `IsMatch(s int) bool`：判断 `dfa.states[s].match`;
- `CanMatch(s int) bool`：只要 `s` 在有效范围内 (`0 < s < len(states)`) 即可；
- `Accept(s int, b byte) int`：从当前状态 `s` 出发，读入字节 `b`，转移到下一个状态 `dfa.states[s].next[b]`；若越界则返回 `0`；
- `MatchesRegex(input string)`：对外提供的简化调用：从 `Start()` 出发依次 `Accept(...)`，看最终是否到达 `IsMatch()` 的状态。

这样就把正则表达式变成了**字节级** DFA！

---

## 2. 整体编译流程

### 2.1 正则解析

```go
parsed, err := syntax.Parse(expr, syntax.Perl)
```

- 这是 Go 自带的正则解析，会把 `expr` 转成 `syntax.Regexp` AST，支持 Perl-like 语法。

### 2.2 `newCompiler(sizeLimit)`

```go
compiler := newCompiler(size)
insts, err := compiler.compile(parsed)
```

- `compiler.compile(ast *syntax.Regexp)`: 递归遍历 `ast`，生成一系列“指令” (type `inst`)。这部分类似于“Thompson 构造法”：
  - `OpSplit`, `OpJmp`, `OpRange`, `OpMatch` 四种指令；
  - “split” 指令用于分支 (正则 OR / Kleene star 结构)，“jmp” 指令用于回环，“range” 指令表示匹配指定字节区间，“match” 表示可接受。

结果是一个“程序” (`prog`) = `[]*inst`，长度可能不能超过用户限制 (`sizeLimit`)。

### 2.3 构造 DFA

```go
dfaBuilder := newDfaBuilder(insts)
dfa, err := dfaBuilder.build()
```

- 先有一张“指令表” `prog`；
- `dfaBuilder` 通过**子集构造**或类似的方式把指令集合编成**确定性**状态机。
- `build()` 里，会创建一个初始集合（包含指令入口），然后对 256 个字节逐个尝试，生成下一个状态的指令集，再缓存起来。最后形成 `dfa.states[]`，每个 state 有一个 `[256]int` 的 `next` 跳转表。

**得到**：`dfa{ insts, states }`，`states` 是一个 slice，每个 `state{ insts []uint, next []int, match bool }`。

- `insts []uint`：该状态活跃的指令集下标；
- `next[b]`：读入字节 `b` 后转移到的状态编号；
- `match`：是否包含 `OpMatch` 指令，表示正则成功匹配。

---

## 3. 代码细节解析

### 3.1 `compiler.compile(...)`

```go
func (c *compiler) compile(ast *syntax.Regexp) (prog, error) {
    err := c.c(ast)
    // ...
    inst := c.allocInst()
    inst.op = OpMatch
    c.insts = append(c.insts, inst)
    return c.insts, nil
}
```

- 调用 `c.c(ast)`: 递归处理 AST；
- 末尾再加一个 `OpMatch` 指令，表示**可接受**；

#### 3.1.1 `compiler.c(ast)`

内部一个大 `switch ast.Op`：

- `OpLiteral`：插入与特定字符（或大小写折叠集）匹配的 byte ranges；
- `OpAnyChar` / `OpAnyCharNotNL`：转换成 `OpCharClass` 处理；
- `OpCharClass`：遍历 `[start, end]` 区间，一一编译；
- `OpConcat` / `OpAlternate` / `OpStar` / `OpPlus` / `OpQuest` / `OpRepeat` 等，对应正则中 `(...|...)`、`(...)*`、`(...)+`、`(...)?` 等操作。
- `OpCapture`：忽略捕获组，仅编译子表达式；
- `OpEmptyMatch`：表示可以匹配空字符串，这里可能无需插入新指令。

每次遇到 “分支” 会插入 `split` 指令、遇到“循环”会配合 `jmp` 指令来实现。

---

### 3.2 指令格式 (`type inst`)

```go
type inst struct {
    op         instOp
    to         uint
    splitA     uint
    splitB     uint
    rangeStart byte
    rangeEnd   byte
}
```

- `op`：`OpMatch`, `OpJmp`, `OpSplit`, `OpRange`;
- `to`：如果是 `OpJmp`，跳转到 `inst[to]`；
- `splitA/splitB`：如果是 `OpSplit`，表示分支 A/B 要跳到的指令索引；
- `rangeStart/rangeEnd`：若是 `OpRange`，在 `[rangeStart..rangeEnd]` 内匹配成功后，继续到**下一条指令** (i.e. `ip+1`)。

### 3.3 编译后 `prog []*inst`

`c.insts` 就是一段“Thompson NFA”式的指令流，这在下一个阶段由 `dfaBuilder` 转换成真正的**DFA**。

---

## 4. `dfaBuilder.build()`：将指令集转为 DFA

**核心思路**：

1. 给定一组活跃指令索引 (比如 `[2, 3, 7]`)，通过 `run(cur, next, b)` 看读入字节 `b` 后会激活哪些新指令；
2. 新的活跃指令集合 -> 这是下一个 DFA 状态；
3. 用一个 `cache` (`map[string]int`) 来去重（相同指令集合 = 同一状态）；
4. 依次 BFS（或 DFS）探索 256 个字节可能性，直到没有新的状态出现。

### 4.1 初始状态

```go
cur := newSparseSet(uint(len(d.dfa.insts)))
next := newSparseSet(uint(len(d.dfa.insts)))
d.dfa.add(cur, 0)
ns, _ := d.cachedState(cur, nil)
// ...
```

- `d.dfa.add(cur, 0)`：激活指令集从第 0 条指令开始 (因为 index=0 常常是 OpSplit/OpJmp leading to real start states)。
- `cachedState`：将活跃指令集合序列化为字符串，查 `cache`。若不存在，就分配一个新的 state idx。
- state 0 在 `dfa.states` 中是保留的“无效”/“死”状态，所以下一个创建的 stateID 就是 1 了 => `Start = 1`。

### 4.2 转移

```go
for b := 0; b < 256; b++ {
    ns, instsReuse = d.runState(cur, next, s, byte(b), instsReuse)
    if ns != 0 {
        ...
    }
}
```

- 依次枚举所有可能的字节 `b`；
- `runState(...)` 先复制当前状态 `s` 的活跃指令到 `cur`，然后对 `b` 做匹配(`d.dfa.run(cur, next, b)`)，得到 `next`；最后 `cachedState(next, instsReuse)` => 下一个 stateID；并填 `dfa.states[s].next[b] = nextState`。
- 若 `nextState != 0` 且没见过就加入 `seen` 待处理队列。

### 4.3 `dfa.run(cur, to, b)`

```go
func (d *dfa) run(from, to *sparseSet, b byte) bool {
    to.Clear()
    for i := 0; i < from.Len(); i++ {
        ip := from.Get(i)
        switch d.insts[ip].op {
        case OpMatch: ...
        case OpRange: if in range => d.add(to, ip+1)
        }
    }
}
```

- 对每个存活的指令下标 `ip`，若是 `OpRange` 并且 `rangeStart <= b <= rangeEnd`，则激活下个指令 `ip+1`; 若是 `OpMatch`，此时 `isMatch=true`（DFA 状态可接受，但不影响激活哪些指令）。
- `d.add(to, someIP)` 再递归地处理 `OpJmp/Split` 指令 (在 `add(set, ip)` 里).

---

## 5. `dfa` 结构

```go
type dfa struct {
    insts  prog
    states []state
}
type state struct {
    insts []uint
    next  []int
    match bool
}
```

- `states[0]` 是保留的死状态。
- `states[1..]` 由编译过程动态创建，每个 state 存：
  - `insts`: 这里主要存当前活跃指令的**下标**；
  - `next[b]`: 读字节 `b` 的转移到另一个 state；
  - `match`: 是否包含 `OpMatch`。

查询时只要 `s != 0` 即可 `CanMatch(s)=true`，`s=0` 表示退出或无法匹配。

---

## 6. UTF-8 / 代码点处理

在 `compiler.compileClassRange / compileUtf8Ranges` 等处，看到**大量**对 `rune -> byte sequences` 的处理：

- 如果一个字符是 `'a'`~`'z'`，很好处理；但若是中文/emoji，需要多字节 UTF-8；
- 这里就把 `rune` 拆成多个 `Range` 指令 (`OpRange`)，让DFA 在“多字节模式”下也能正确匹配这个 Unicode 字符。
- `NewSequencesPrealloc(...)` / `SequenceFromEncodedRange(...)` 都是将 `[startR, endR]` 的 Unicode range 转成多个**字节序** ranges。
- 结果是在匹配时，对输入串**逐字节**喂给 Automaton，能正确识别多字节 UTF-8。

---

## 7. 其它辅助部分

### 7.1 稀疏集合 (sparseSet)

```go
type sparseSet struct {
    dense  []uint
    sparse []uint
    size   uint
}
```

- 常见优化结构，用 `dense[i]` 表示第 `i` 个元素；`sparse[x]` 存储 `x` 在 `dense` 中的下标；
- `Add(x)`, `Contains(x)`, `Clear()` 方便在 NFA/DFA 算法中快速插入和检查。

### 7.2 `ErrNoEmpty/ErrNoWordBoundary/...`

- 提示用户：若正则表达式里出现零宽断言(`^`, `$`, `\b`)等，就报错：**不支持**，因为无法在简单的字节级自动机上处理这些断言。

### 7.3 大小限制

```go
if uint(len(c.insts)*instSize) > c.sizeLimit {
    return ErrCompiledTooBig
}
```

- 若指令过多，就报错 `ErrCompiledTooBig`。默认限制 10MB。

---

## 8. 运行时匹配

1. 通过 `MatchesRegex(input string) bool`：

   ```go
   currentState := r.Start() // 1
   index := 0
   for r.CanMatch(currentState) && index < len(input) {
       currentState = r.Accept(currentState, input[index])
       index++
   }
   return index == len(input) && r.IsMatch(currentState)
   ```

   - 依次读 `input[index]`，调用 `Accept(currentState, b)` = `dfa.states[currentState].next[b]` => nextState。
   - 若到 0 则失败 (死状态)；
   - 最后若读完 `input` 并且 `IsMatch(currentState)=true`，匹配成功。

2. 和 Go 原生的 `regexp` 相比，这里只处理一部分 Perl 语法特性，不支持分组捕获等**运行时**功能；也对一些“零宽断言”做了拒绝；好处是**可将模式编译成一个独立的 DFA**，极大地提高匹配速度，适合**大量重复匹配**场景。

---

## 9. 小结与要点

1. **总体架构**：先把正则表达式解析成**指令** (NFA)；再**子集构造**成**DFA**；最后查询时只需 O(n) (n=输入长度) 的状态转移。
2. **UTF-8** 处理：通过将 `[startRune..endRune]` 拆成若干个 “byte range” (`Range {Start, End}`)，用 `OpRange` 指令在字节层面完成匹配；
3. **限制**：不支持**零宽断言**(`^`, `$`)、`\b`、懒量词 `?+` 等，以及**最大 10000 states** (StateLimit)。
4. **性能**：编译后查询速度快；若使用大量或复杂表达式可能导致指令数/状态数过大，触发 `ErrCompiledTooBig` 或 `ErrTooManyStates`。
5. **代码结构**：
   - `compile` 阶段：`compiler` 递归遍历 `syntax.Regexp` => `prog (insts)`；
   - `dfaBuilder.build()` 阶段：BFS 或类似方式在指令集上跑，构建 `dfa.states[]`；
   - 运行时：`Regexp.MatchesRegex(input)` 只做简单的字节扫描 + 数组索引。

这是一个可行的**“纯 DFA”正则实现**，更适合**静态场景**（编译一次、多次匹配）。在实际项目中，若需要捕获组或更复杂特性，仍需要 Go 标准 `regexp` 包；若只需要快速判断**某字符串是否匹配**固定模式，则可使用本库来取得**高性能**和**可移植性**（因为它序列化后的 `dfa` 可独立存储/加载）。
