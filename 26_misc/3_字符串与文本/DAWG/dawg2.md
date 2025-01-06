下面这段代码实现了一个带“模糊搜索”功能的 **DAWG（Directed Acyclic Word Graph）**。整体上，可以把它分为几个主要部分：

1. **核心数据结构**：`DAWG`、`state`、`letter`
2. **构建 & 压缩（Trie -> DAWG）**：`addWord`、`compressTrie`、`analyseSubTrie`
3. **查询 / 模糊搜索**：`Search`、`searchSubString`
4. **序列化 / 反序列化**：`SaveToFile`、`LoadDAWGFromFile`
5. **辅助函数**：如 `equals`、`containsLetter`、`getletter` 等

下面会按功能模块进行分析，并对代码中的关键点做详细解释和点评。

---

## 1. 数据结构

### 1.1. `DAWG` 结构

```go
type DAWG struct {
    initialState *state
    nodesCount   uint64
}
```

- `initialState`: 整个图（或说 Trie / DAG）的初始节点（根节点）。
- `nodesCount`: 在构建过程中统计的节点数（注意，这里是最终合并后的节点数）。

### 1.2. `state` 结构

```go
type state struct {
    final bool

    letters      *letter // 该节点的子边（根指针）
    lettersCount int     // 子边数量

    next   *state  // 用于将相同“层级”的 state 串成链表，后续合并时会用到
    letter *letter // 指向从哪个 letter 来到这个 state（在合并时使用）
    number uint64  // 用于序列化/反序列化时的编号
}
```

- `final`: 表示这个状态是否是某个单词的结束节点（类似 Trie 中的“终止标志”）。
- `letters`: 指向一棵由 `letter` 组成的**二叉搜索树**（BST），同时也通过 `next` 串成**链表**。这在查找时可以 O(log n) 搜索，也可以 O(n) 遍历。
- `lettersCount`: 当前 `state` 下有多少个子边（子节点）。
- `next`: 在“压缩”阶段，用于把**同一层**上的所有 `state` 串成链表，以便进行重复检测和合并。
- `letter`: 回指“是哪条边”到达了当前 `state`。在合并时，如果合并成功，需要回溯修改指针。
- `number`: 序列化/反序列化用的唯一标识。

### 1.3. `letter` 结构

```go
type letter struct {
    char  rune
    state *state

    // BST 相关
    left  *letter
    right *letter

    // 把同一个 state 下的 letter 串到一个单链表里
    next  *letter
}
```

- `char`: 该字母边上存储的字符（`rune` 以兼容 Unicode）。
- `state`: 这条边指向的下一个状态。
- `left` / `right`: 使得当前节点在 BST 中可以进行插入/查找（按 `char` 大小左右分支）。
- `next`: 用来串联同一个 `state` 所有的 `letter`，形成一个单链表，以便在合并、遍历时能方便一次性遍历所有子边。

> **特别之处**：一方面使用二叉搜索树做“精确查找”子字符，另一方面又用单链表方便“顺序遍历”。这是一个相对独特的折衷实现。

---

## 2. Trie 构建与 DAWG 压缩

### 2.1. 从输入单词构建 Trie

- **核心函数**：`addWord(initialState, word) -> (newEndState bool, wordSize int, createdNodes uint64)`

```go
func addWord(initialState *state, word string) (bool, int, uint64) {
    curState := initialState
    var createdNodes uint64
    var wordSize int

    for _, l := range word {
        // 1) 在 curState.letters（BST）里找到对应字符 'l' 的 letter
        //    若无则新建 letter
        // 2) 若 letter.state == nil 表示还没有下一层 state，就新建一个
        //    并记录 createdNodes++
        // 3) 移动 curState = letter.state
        // 4) wordSize++ (统计当前单词的字符数)
    }
    curState.final = true
    return
}
```

- **过程**：相当于传统的 Trie 插入：
  1. 从根节点（`initialState`）开始，逐个字符向下走；
  2. 如果对应的 `letter` 不存在，就新建 `letter`；如果 `letter` 的 `state` 不存在，也要新建；
  3. 最终把最后一个 `state` 标记为 `final` = `true`。
- **BST 插入**：它通过 `curLetter.left / right` 来判断要走左子树还是右子树，直到找到 `curLetter.char == l` 或创建一个。
- **统计**：`createdNodes` 记录新建的 `state` 数；`wordSize` 记录单词长度。

### 2.2. 压缩：`compressTrie(initialState, maxWordSize)`

构建完 Trie 后，需要进行最重要的一步：**合并重复子树**（压缩成一个 DAWG）。

- 大致思路：
  1. 先做一次 DFS/递归（`analyseSubTrie`）来**收集同一“层”的所有 `state`**。
     - “层”这里指的是离根节点的深度，比如第 0 层是 `initialState`，第 1 层是它的子节点……
     - 不同层的状态显然不能互相合并（合并条件要求两者结构、后缀相同）。
  2. 在同一层里，尝试两两比较 `state.equals(...)`，如果发现两个 state 是完全等价的，就把其中一个合并到另一个上，`deletedNodes++`。

```go
func compressTrie(initialState *state, maxWordSize int) (deletedNodes uint64) {
    levels := make([]*state, maxWordSize)
    ...
    // 1) analyseSubTrie: 把各层节点串起来
    // 2) 在每层里，合并等价的 state
    //    for curState := levels[i]; curState != nil && curState.next != nil; curState = curState.next {
    //        for previousState, sameState := curState, curState.next; sameState != nil; sameState = sameState.next {
    //            if curState.equals(sameState) {
    //                // 合并
    //                previousState.next = sameState.next
    //                sameState.letter.state = curState
    //                deletedNodes++
    //            } else {
    //                previousState = sameState
    //            }
    //        }
    //    }
    return
}
```

#### 2.2.1. `analyseSubTrie(curState, levels, channels)`

- 这是一个**递归**函数，用来遍历 `curState` 的全部子节点，计算最大的深度（`curLevel`）；
- 每访问到当前节点，就把它挂到 `levels[curLevel]` 链表里去 (`curState.next = levels[curLevel]; levels[curLevel] = curState`)；
- 代码里还使用了 `channels[curLevel]` 做同步，可能是作者尝试对不同分支进行**并发**遍历，然后再合并结果。不过实现比较简易，可能存在一些竞态或等待逻辑上的疑问。

#### 2.2.2. `state.equals(otherState)`

- 判断两个状态是否可以被合并的核心：
  1. `final` 标志相同（是否都为单词结束）；
  2. `lettersCount` 相同；
  3. 对每个 `letter`，在 `otherState` 中都存在**相同字符**、且**指向的子状态是同一个**的 `letter`。

```go
func (state *state) equals(otherState *state) bool {
    if state.final != otherState.final || state.lettersCount != otherState.lettersCount {
        return false
    }
    for curLetter := state.letters; curLetter != nil; curLetter = curLetter.next {
        if !otherState.containsLetter(curLetter) {
            return false
        }
    }
    return true
}
```

- `otherState.containsLetter(curLetter)`: 在 BST 中二分搜索 `curLetter.char`，并且还检查 `curLetter.state == letter.state`。

> **关键**：如果两者 `lettersCount` 相等，但其指向的子节点在之前**已经合并**过，可能会让多个子节点都指向同一个状态，这也被认为是**“相同结构”**。

#### 2.2.3. 合并过程

- 当 `curState.equals(sameState)` 判断为真时：
  - `previousState.next = sameState.next` 将 `sameState` 从链表中摘除；
  - `sameState.letter.state = curState`：把 `sameState` 原先指向自己的父 letter 改成 `curState`，从而**合并**了这两个节点的后续路径。
  - `deletedNodes++`：统计合并次数。

---

## 3. 查询与模糊搜索

### 3.1. 精确查询

- 代码中并没有一个“最纯粹”的精确查询函数，但可以用 `getletter` 来逐字符往下走，最后判断 `state.final` 来确定单词是否存在。
- `getletter(letter rune)`：在 BST 里按照 `char` 大小左右移动，直到找到匹配的 `letter` 或 `nil`。这就是 Trie 中的查找过程。

### 3.2. 模糊搜索：`Search`

```go
func (dawg *DAWG) Search(
    word string,
    levenshteinDistance int,
    maxResults int,
    allowAdd bool,
    allowDelete bool,
) (words []string, err error) {
    wordsFound, _, wordsSize, err := searchSubString(
        dawg.initialState,
        *bytes.NewBufferString(""),
        *bytes.NewBufferString(word),
        levenshteinDistance, maxResults, allowAdd, allowDelete, 0,
    )
    ...
}
```

- 调用 `searchSubString` 去做递归的模糊匹配。
- 参数含义：
  - `levenshteinDistance`: 最大可允许的编辑距离（Levenshtein）。
  - `maxResults`: 最多返回多少个结果，避免海量结果超时或消耗过大。
  - `allowAdd`: 是否允许插入字符（对比原 word 时多加一个字符）。
  - `allowDelete`: 是否允许删除字符（对比原 word 时少一个字符）。

#### 3.2.1. `searchSubString` 函数

这是一个**核心的递归搜索**：

```go
func searchSubString(
    state *state,
    start bytes.Buffer,
    end bytes.Buffer,
    levenshteinDistance int,
    maxResults int,
    allowAdd bool,
    allowDelete bool,
    ignoreChar rune,
) (words *word, lastWord *word, wordsSize int, er error)
```

- `start`: 当前已匹配的前缀（在递归过程中不断增长或回退）。
- `end`: 还没匹配的剩余字符（从 `word` 中读取）。
- `levenshteinDistance`: 还剩多少可“容错”编辑的机会。
- `ignoreChar`: 递归里处理“替换字符”或“跳过字符”时，用来记住被替换/跳过的字符。
- 核心逻辑：
  1. **读一个字符** `char` = `end.ReadRune()`
     - 如果在当前 `state` 能找到 `letter.char == char`，就正常匹配下去；
     - 否则尝试执行编辑操作（替换/删除/新增），如果还有剩余的 `levenshteinDistance`。
  2. 如果 `end` 已经读完，且 `state.final`，说明这是一个匹配结尾，就把当前 `start` 的内容记录到结果里。
  3. 如果允许 `allowAdd`、`allowDelete` 并且还有编辑距离剩余，则分别尝试**插入一个新字符**或**删除当前字符**。
  4. 通过 `mergeWords` 把不同分支递归得到的结果链表合并在一起，记录返回结果。

> 这个模糊搜索的实现比较直接，属于“回溯 + 递归”策略，不同分支都要遍历，可能性能较高（指数级），但在小规模或有限 `levenshteinDistance`、`maxResults` 时也足够。

---

## 4. 序列化 & 反序列化

### 4.1. `SaveToFile(fileName string)`

- 先写入 `dawg.nodesCount`；
- 然后深度遍历（`saveSubTrieToFile`）每个 `state`，依次为其分配一个 `number`（递增），并将 `final` 标志、子 `letter`（字符 + 状态编号）写入文件。

### 4.2. `LoadDAWGFromFile(fileName string)`

- 先读 `nbNodes`；
- 然后创建一个 `states` 切片大小 `nbNodes`；
- 依次读取每行，构建对应的 `state` 对象，并把 `char -> linkedNodeNumber` 的关系恢复到各自的 `letter`；
- 构建完成后，返回新的 `DAWG`.

---

## 5. 其他函数

- `FindRandomWord(wordSize int) (string, error)`  
  随机在 DAWG 中走若干步（`wordSize` 步），尝试构造一个单词。若最后停在一个 `final` 状态，就返回这个随机单词。
  - 若某一步发现没有子节点（`state.lettersCount == 0`），就放弃当前路线，重新来（通过 `continue INFINITE`）。
  - 这在结构稀疏时，可能会很长时间找不到合法长度的单词。

---

## 6. 代码整体点评

1. **数据结构设计**

   - 使用了二叉搜索树 + 单链表的组合来管理子节点（`letter`），在查询和遍历上做了折衷。
   - `state` 中的 `next`、`letter` 指针是为了**辅助合并**，这在 Trie -> DAWG 的转换上是一个可行的思路。

2. **Trie 构建**

   - `addWord` 的逻辑与常规 Trie 插入相似，但把孩子节点存在 BST 中，而非哈希/数组。
   - 统计新节点数量、最长单词长度，以便后续压缩和分层。

3. **压缩逻辑（`compressTrie`）**

   - `analyseSubTrie`: 递归收集各层节点；
   - 在同一层上，用 `equals()` 判断两个状态是否可合并；
   - `equals()` 基于**子节点逐一比对**和 `final` 标志的比较。
   - 合并时，将“被合并”的 `state` 的指向改到“保留的”那个 `state` 上，并删除多余节点。
   - 并发设计略显简单（`channels` 等），但核心思路仍是 BFS/DFS + 比对 + 合并。

4. **模糊搜索**

   - 采用递归回溯，支持“替换、删除、插入”来进行 Levenshtein 距离计算；
   - 复杂度可能较大，但对小距离或有限 `maxResults` 尚可。

5. **序列化 & 反序列化**

   - 通过给每个 `state` 动态分配一个 `number`，写出每个 state 的内容（`final` 标志、子节点信息），再按相同编号读回；
   - 比较简单、直接，也能方便后续部署与持久化。

6. **可能的改进**
   - 对 BST + 链表的结构做一些更高效的封装；
   - 并发处理的正确性和性能可能需要更多测试或更严谨的锁机制；
   - 模糊匹配的回溯算法在大数据规模/大距离时可能效率欠佳，可以考虑改进为更典型的“编辑距离 + Trie”DP 解法；
   - `FindRandomWord` 在可能没有固定长度单词的情况下会死循环。

---

## 7. 总体总结

这段代码展示了一个**从 Trie 构建并压缩成 DAWG**的完整流程，同时还实现了“模糊搜索（Levenshtein 距离）”、序列化/反序列化等功能。核心亮点包括：

- **压缩成 DAWG**的实现：

  - 先用 Trie 存储所有单词；
  - 再通过“层次遍历 + 等价判断”来合并重复的后缀子树。

- **查找方式**：

  - 精确查找：在 BST 中按字符查找；
  - 模糊查找：通过递归回溯，对可选的编辑操作进行分支搜索。

- **使用 BST + 链表**来管理子节点，不是传统的 `map[char]*state` 或数组结构，独具特色。

总体而言，这份代码对 **DAWG** 的压缩思路是相对“朴素”的：比较两个节点的所有子树是否相同，然后合并。它演示了 DAWG 的一个可行实现，也为读者提供了如何做**“最小化 Trie”**的思路。
