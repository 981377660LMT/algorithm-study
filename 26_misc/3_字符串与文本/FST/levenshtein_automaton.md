下面这段代码来自 [blevesearch/vellum/levenshtein](https://github.com/blevesearch/vellum/tree/master/levenshtein) 项目，实现了一个构建与使用 **Levenshtein Automaton** 的工具，用来在给定的最大可容忍编辑距离（`maxDistance`）下，对字符串进行近似匹配、计算编辑距离等操作。

为了帮助你深入理解，该代码包含多个主要部分，我们可以从以下几个方面来解析：

1. **Levenshtein Automaton 简介**
2. **核心数据结构：NFA、ParametricDFA、DFA**
3. **Levenshtein NFA 构建**
4. **ParametricDFA（带参数的自动机）**
5. **Builder / BuildDfa** 流程
6. **距离的定义与处理**
7. **字母表（Alphabet）与特征向量（Characteristic Vector）**
8. **DFA 查询/匹配逻辑**
9. **代码的整体流程与关键要点**

---

## 1. Levenshtein Automaton 简介

**Levenshtein Automaton** 是一种基于有穷自动机的结构，可以识别与某个字符串在 **Levenshtein 编辑距离** 不超过 \(d\) 的所有字符串。通常，它对应了一个“在线计算编辑距离”的过程，只需要对输入字符串执行一次状态转移即可判断是否满足“编辑距离 \(\le d\)”。

- 其中，**编辑操作**（插入、删除、替换，以及可选的转置）每次消耗一次“距离”，当总消耗超过 \(d\) 时自动机拒绝该字符串。
- 该项目中还支持可选的**Damerau** 转置（`transposition=true`），即把相邻字符交换也算距离 1。

然而，直接构建传统的“Levenshtein DFA”可能比较庞大，尤其当 \(d\) 较大时。这里使用了一个“**Parametric DFA**”的思想，将其做了一层抽象：先在 NFA（不确定性自动机）中进行通用的“转移预计算”，再映射到更为紧凑的 “ParametricDFA”，最后可针对不同查询字符串构建出针对性的 “DFA”。

---

## 2. 核心数据结构：NFA、ParametricDFA、DFA

### 2.1 NFA（LevenshteinNFA）

- `LevenshteinNFA` 保存了 `mDistance`（最大编辑距离）以及是否允许 `damerau`（转置操作）。
- 它可以针对某个 pattern 构建出相应的多重状态（MultiState），并在读取字符时进行 NFA 转移，累积编辑距离。

### 2.2 ParametricDFA

- 将 LevenshteinNFA 的每个**多重状态 (MultiState)** 进行“归一化”后哈希去重，得到一系列**形状 (shape)**。
- 对于输入字符集的各种“特征向量 (Characteristic Vector)”进行跳转表 `transitions` 预先存储，形成一个“**带参数**”的自动机。
- 当真正要匹配具体字符串时，再结合此 ParametricDFA 和查询串，生成/执行具体的 DFA。

### 2.3 最终 DFA

- 最后得到的 `DFA` 就是对 UTF-8 或普通字节进行匹配的确定自动机，其中 `DFA.transitions[fromState][byte] = toState`。
- 调用 `IsMatch(state)`、`Accept(state, b)` 等方法可以在线判断某字符串的编辑距离是否 \(\le d\)。

---

## 3. Levenshtein NFA 构建

```go
type LevenshteinNFA struct {
    mDistance uint8
    damerau   bool
}
```

- `mDistance`: 最大可允许的编辑距离 \(d\)。
- `damerau`: 是否支持转置（Damerau）。

NFA 的状态类型：`NFAState`，其中包含：

- `Offset`: 当前处理到的模式串位置
- `Distance`: 已经消耗的编辑距离
- `InTranspose`: 标记此状态下是否正处于转置上下文。

**核心**是 `func (la *LevenshteinNFA) transition(cState *MultiState, dState *MultiState, scv uint64)`:

- `cState` 是当前多重状态集合；
- `scv` 是针对下一输入字符构造的“特征向量” (bitmask)，表示哪些位置字符匹配；
- 根据规则（插入/删除/替换/转置），把下一步可能到达的所有 `NFAState` 加到 dState 里。

最后 `la.multistateDistance(ms, offset)` 可以得到某个 `MultiState` 到结尾的最小距离。

---

## 4. ParametricDFA（带参数的自动机）

`func fromNfa(nfa *LevenshteinNFA) (*ParametricDFA, error)`

1. **哈希表 (`hash`)**：把 NFA 中出现的所有 `MultiState` 收集起来，并为它们分配一个 ID（`lookUp.getOrAllocate`）。
2. **构造 transitions**：对每个 `MultiState`，对所有可能的特征向量 (chi) 做 NFA 转移得到 `destMs`，再通过 `normalize()` 使之有序/最小化 offset，获得 `destID`。就形成了 `Transition{destShapeID, deltaOffset}`：
   - `deltaOffset` 表示 offset 的移动量，用于后续计入 ParametricDFA 状态下的 `offset`。
3. **记录 distance**：对每个 MultiState，不同 offset 下的编辑距离存入 `pdfa.distance[]`。

这样就得到一张通用“形状 + offset”的跳转表 (`pdfa.transitions[]`)。

**ParametricState** = \(\text{shapeID} + \text{offset}\)。对不同输入串时，会动态更新 offset 并根据 shapeID 去做转移。

---

## 5. Builder / BuildDfa 流程

```go
func (lab *LevenshteinAutomatonBuilder) BuildDfa(query string, fuzziness uint8) (*DFA, error) {
    return lab.pDfa.buildDfa(query, fuzziness, false)
}
```

- `NewLevenshteinAutomatonBuilder(maxDistance, transposition)`：先通过 `newLevenshtein(...)` 构建 NFA，再 `fromNfa(...)` 获得 `ParametricDFA`。这是预处理，每次只做一次，耗时可能较大，但可复用。
- 当要针对一个具体 query + fuzziness 构建实际的 DFA 时，调用 `buildDfa(query, distance, prefix)`:
  1. 取 `qLen = len(query)`；
  2. 初始化一个 `ParametricStateIndex`（psi），用来记录形状 + offset -> 具体的 stateID；
  3. BFS/DFS 式地扩展 state：
     - 读取特征向量 “\(\chi\)” 对应的 transition -> next ParametricState；
     - offset 可能随之增加（`deltaOffset`）。
     - 分配新的 stateID 并存到 DFA builder 里。
  4. 若状态数超过 `StateLimit = 10000` 就报错。
  5. `dfaBuilder.build(distance)` 最终得到 `*DFA`。

---

## 6. 距离的定义与处理

代码中定义了 `Distance` 接口，包含两种实现：

- `Exact{d uint8}`：表示确切距离
- `Atleast{d uint8}`：表示距离至少是多少，如果大于 `maxDistance`，就无法再区分更大距离了，等价于“超限”。

在 `ParametricDFA` 中，`getDistance(state, qLen)` 会根据 `state.offset`、`pdfa.distance[...]` 来获得某个状态对应的距离。如果超过 `maxDistance`，就返回 `Atleast`。

查询时，如果 `Distance` 是 `Exact{d}` 并 `d <= maxDistance`，则在终止点视为可接受（匹配）；否则不可接受。

---

## 7. 字母表（Alphabet）与特征向量（Characteristic Vector）

```go
type Alphabet struct {
    charset []tuple
    index   uint32
}
type tuple struct {
    char rune
    fcv  FullCharacteristicVector
}
```

- 构建时，会将 query 中的字符去重+排序，给每个字符 c 生成一个特征向量 `fcv`，表明在 query 各位置上出现了 c 的哪些位置。
- 例如，若 `query = "cats"`, `fcv('c')` 表示 bitmask：第 0 位是 1，其它是 0； `fcv('a')` 第1位是1等。
- 这样当读取字符 c 时，可以快速得到“匹配位置 bitmask”，再推进 `ParametricDFA` 中的 offset。
- `fcv.shiftAndMask(offset, mask)` 会把 bitmask 向右 shift 一段，以便对应当前 offset 的对齐。

在 `buildDfa()` 中：

```go
chr, cv, err := alphabet.next()
chi := cv.shiftAndMask(state.offset, mask)
transition := pdfa.transition(state, chi)
```

这一步便是关键：根据 **(offset + 特征向量)** 推断下个 state。

---

## 8. DFA 查询 / 匹配逻辑

生成的 `DFA`:

```go
type DFA struct {
    transitions [][256]uint32
    distances   []Distance
    initState   int
    ed          uint8
}
```

- `transitions[s][b] = s'`：状态 `s` 在输入字节 `b` 时转到状态 `s'`；
- `distances[s]` 是状态 `s` 对应的编辑距离（或 Atleast）。
- `IsMatch(s)`: 若 `distances[s]` 是 `Exact{d}` 且 `d <= ed`，则可匹配。
- `Accept(s, b)`: 做 `transitions[s][b]`;
- `MatchAndDistance(input string)`: 从 `initState` 出发，对每个字符 `b` 做一步转移，如果中途到达 `SinkState` 就表示编辑距离超限；若最终在末尾状态 `IsMatch()`，返回 `(true, dist)`。

---

## 9. 代码的整体流程与关键要点

1. **NFA 构建**

   - `newLevenshtein(d, transposition)` -> `LevenshteinNFA`。
   - 它定义了编辑操作规则（插入/删除/替换/转置），可计算多重状态 `MultiState` 的演变。

2. **fromNfa**

   - 通过哈希表记录所有可能的 `MultiState` 形态。
   - 对所有特征向量 (chi) 枚举转移，用 `(destShapeID, deltaOffset)` 的形式记录到 `ParametricDFA.transitions` 中。
   - 也将 `distance[]` 填充，以便后续查询时读取状态对应的距离。

3. **ParametricDFA.buildDfa**(query, fuzziness, prefix)

   - 解析 query 的字符集，构建 `Alphabet`；
   - 从 `initialState` 出发，不断对特征向量做转移，并在 “**形状ID + offset**” 与 “DFA 状态ID” 之间做映射。
   - 最终返回一个对普通字节进行匹配的 `*DFA`。

4. **DFA**

   - `transitions[s][byte]` 给出下一个状态下标，`distances[s]` 表示此状态到终点编辑距离；
   - 提供 `IsMatch(s)` / `Accept(s,b)` / `MatchAndDistance(str)` 等方法。

5. **工厂模式**

   - `NewLevenshteinAutomatonBuilder(maxDist, transposition)` 构建一个可复用的 `LevenshteinAutomatonBuilder`。
   - 之后多次 `BuildDfa(query, fuzziness)` 可快速生成**针对不同 query** 的定制 DFA。
   - 如此避免重复构建 NFA / ParametricDFA。

6. **主要用途**
   - 实现近似字符串匹配：给定一个最大距离 d、一个 query，生成 DFA；对候选字符串逐个喂给 DFA，看是否到达可接受状态。
   - 在 vellum 中，也可与 FST 结合，对**大规模**键集合做“编辑距离 \(\le d\)”搜索。

---

### 结论

**整段代码**的使命是：

1. **构建**一个可复用的（Parametric）Levenshtein 自动机，用于指定的最大距离 `maxDistance`。
2. **针对具体的查询字符串**快速生成确定性 DFA；
3. **利用 DFA** 在线、一次性地判断任何候选字符串的编辑距离是否 \(\le d\)。

**关键思想**：

- 利用 Levenshtein NFA 表达所有可能的编辑操作路径；
- 以哈希方式去重多重状态，得到紧凑表示（ParametricDFA）；
- 构建时通过 bitmask（特征向量）减少重复；
- 查询时只需 O(length of string) 的状态跳转来判断匹配性和计算距离。

这份实现相对复杂，但非常适合在搜索引擎、模糊查找、自动补全等场景中，用有限状态机（DFA）来**高效判断**是否在给定编辑距离范围内。

如果你要理解或改进这段代码，可以重点关注：

- **NFA** 中的编辑操作如何编程实现；
- **ParametricDFA** 的 `transitions` 构造；
- **BuildDfa** 时如何基于 query 生成最终用字节跳转的 DFA；
- **Characteristic Vector** 的构造与 shift 用途；
- **distance** 与 `Atleast/Exact` 的判断，防止搜索过度扩张。

这些要点串在一起，就构成了完整的 “Levenshtein Automaton” 体系。
