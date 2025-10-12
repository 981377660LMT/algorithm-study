### fastDiff.js 代码详解

这是一个用于**文本比较（Diff）**的 JavaScript 库。它源自 Google 著名的 `diff-match-patch` 库，但经过了精简，只保留了核心的“比较（Diff）”功能，移除了“补丁（Patch）”和“匹配（Match）”部分，使其更轻量、更专注。

#### 核心目标

该库的核心目标是接收两个文本字符串（旧文本和新文本），并生成一个描述两者差异的“差异序列”。

#### 核心数据结构

它用一个数组来表示差异，数组中的每个元素是一个元组（一个包含两个元素的数组），格式为 `[操作, 文本]`：

- `[DIFF_EQUAL, '一些文字']` 或 `[0, '一些文字']`：表示这部分文字在两个版本中都存在，是**相同**的。
- `[DIFF_DELETE, '一些文字']` 或 `[-1, '一些文字']`：表示这部分文字在旧文本中有，但在新文本中被**删除**了。
- `[DIFF_INSERT, '一些文字']` 或 `[1, '一些文字']`：表示这部分文字在旧文本中没有，但在新文本中被**插入**了。

例如，将 "Hello world" 改为 "Goodbye world"，结果会是：
`[[-1, "Hello"], [1, "Goodbye"], [0, " world"]]`

#### 算法工作流程

代码的执行可以分为几个主要步骤，由 `diff_main` 函数 orchestrate：

1.  **入口与快速检查 (`diff_main`)**

    - 首先，检查两个文本是否完全相等。如果是，直接返回一个 `DIFF_EQUAL` 块，提高效率。
    - 它还会尝试一个针对光标位置优化的快速路径 `find_cursor_edit_diff`，用于处理典型的文本编辑器输入场景（例如，在光标处插入/删除字符）。

2.  **剥离公共部分（优化）**

    - `diff_commonPrefix`: 检查并剥离两个文本开头相同的**公共前缀**。
    - `diff_commonSuffix`: 检查并剥离两个文本末尾相同的**公共后缀**。
    - **目的**：这样做极大地缩小了需要进行复杂比较的文本范围，是重要的性能优化。例如，比较 "The quick brown fox" 和 "The slow brown fox"，可以立即剥离前缀 "The " 和后缀 " brown fox"，只需要比较 "quick" 和 "slow"。

3.  **核心差异计算 (`diff_compute_`)**

    - 这是处理被剥离前后缀后剩下的中间部分的函数。
    - 它包含一些快速路径：
      - 如果一个文本为空，那么差异就是单纯的插入或删除。
      - 如果一个文本完全包含在另一个文本中，差异就是首尾的插入或删除。
    - **半匹配启发式算法 (`diff_halfMatch_`)**：这是一种更高级的优化。它尝试寻找一个较长的公共子串（至少是长文本的一半），如果找到，就将问题一分为二，递归地解决公共子串前后的差异。这对于包含大块移动代码的场景非常有效。
    - **Myers 差分算法 (`diff_bisect_`)**：如果以上所有优化都不适用，代码将回退到这个核心算法。这是 Eugene Myers 在 1986 年提出的经典算法，它能在 O(ND) 时间复杂度内找到最短的编辑路径（N 是文本总长度，D 是差异量）。它通过在“编辑图”上从前向和后向同时搜索，直到路径相遇，从而高效地找到差异点。

4.  **结果合并与清理**
    - **合并**：将第 2 步中剥离的公共前缀和后缀作为 `DIFF_EQUAL` 块重新加到差异序列的首尾。
    - **合并相邻块 (`diff_cleanupMerge`)**：将连续的插入、删除或相等块合并成一个。例如 `[[-1, "a"], [-1, "b"]]` 会被合并为 `[[-1, "ab"]]`。
    - **语义化清理 (`diff_cleanupSemantic`)**：这是为了让差异结果更符合人类的阅读习惯。例如，将 `"The c<ins>at c</ins>ame."` 调整为 `"The <ins>cat </ins>came."`。它会移动编辑边界，使其尽量落在单词边界或换行符上，而不是单词中间。

### 使用案例

由于这是一个 `module.exports` 的 Node.js 模块，你可以在你的项目中使用它。

**1. 安装/准备**

假设你已经将 fastDiff.js 文件放在你的项目目录中。

**2. 创建一个使用它的文件 (e.g., `test.js`)**

```javascript
const diff = require('./fastDiff.js')

// 定义两个要比较的字符串
const text1 = 'The quick brown fox jumps over the lazy dog.'
const text2 = 'The quick red fox leaps over the lazy cat.'

// 调用 diff 函数
const differences = diff(text1, text2)

// 打印结果
console.log('原始差异序列:')
console.log(JSON.stringify(differences))

// 为了更直观地展示，我们可以写一个辅助函数来格式化输出
function prettyPrint(diffs) {
  let html = ''
  for (const [op, data] of diffs) {
    const text = data
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/\n/g, '&para;<br>')
    switch (op) {
      case diff.INSERT: // 1
        html += `<ins style="background:#e6ffe6;">${text}</ins>`
        break
      case diff.DELETE: // -1
        html += `<del style="background:#ffe6e6;">${text}</del>`
        break
      case diff.EQUAL: // 0
        html += `<span>${text}</span>`
        break
    }
  }
  return html
}

console.log('\n可视化结果 (HTML):')
const htmlResult = prettyPrint(differences)
console.log(htmlResult)

// 另一个例子
const oldCode = `function hello() {
  console.log("world");
}`
const newCode = `function hello() {
  console.log("JavaScript");
  return true;
}`

const codeDiffs = diff(oldCode, newCode)
console.log('\n代码比较的可视化结果:')
console.log(prettyPrint(codeDiffs))
```

**3. 运行**

在终端中，进入该目录并运行：

```sh
node test.js
```

**4. 预期输出**

```
原始差异序列:
[[0,"The quick "],[-1,"brown"],[-1," fox "],[-1,"jumps"],[1,"red"],[1," fox "],[1,"leaps"],[0," over the lazy "],[-1,"dog"],[1,"cat"],[0,"."]]

可视化结果 (HTML):
<span>The quick </span><del style="background:#ffe6e6;">brown</del><del style="background:#ffe6e6;"> fox </del><del style="background:#ffe6e6;">jumps</del><ins style="background:#e6ffe6;">red</ins><ins style="background:#e6ffe6;"> fox </ins><ins style="background:#e6ffe6;">leaps</ins><span> over the lazy </span><del style="background:#ffe6e6;">dog</del><ins style="background:#e6ffe6;">cat</ins><span>.</span>

代码比较的可视化结果:
<span>function hello() {&para;<br>  console.log("</span><del style="background:#ffe6e6;">world</del><ins style="background:#e6ffe6;">JavaScript</ins><span>");&para;<br></span><ins style="background:#e6ffe6;">  return true;&para;<br></ins><span>}</span>
```

这个库非常适合用于需要文本比较功能的场景，例如：

- 代码版本控制系统（类似 Git 的 diff 可视化）。
- 在线文档的修订历史。
- 自动评分系统，比较学生答案与标准答案的差异。

---

好的，我们来详细讲解 fastDiff.js 中 `diff_main` 函数的两个可选参数：`cursor_pos` 和 `cleanup`。

---

### 1. `cursor_pos`

#### 作用与目的

`cursor_pos` 参数是一个**性能优化**，主要用于**文本编辑器**这类场景。

在文本编辑器中，用户的编辑行为通常是**局部**的，例如：

- 在光标处输入几个字符（插入）。
- 按退格键或 Delete 键删除几个字符（删除）。
- 选中一段文本，然后输入新内容（替换）。

这些操作的共同点是，文本的大部分内容保持不变，只有光标或选区周围发生了微小的变化。在这种情况下，如果还去执行完整的 Myers 差分算法，就显得“杀鸡用牛刀”，效率很低。

`cursor_pos` 的目的就是**提供关于编辑发生位置的“提示”**，让 `diff` 算法可以尝试一条**快速路径（shortcut）**，从而避免复杂的全局比较。

#### 参数格式

`cursor_pos` 可以是两种格式：

1.  **`number`**: 一个简单的数字，表示旧文本 (`text1`) 中光标的位置索引。这通常对应于没有选区时的插入或删除操作。

2.  **`Object`**: 一个包含更详细信息的对象，格式如下：
    ```javascript
    {
      oldRange: { index: number, length: number },
      newRange: { index: number, length: number }
    }
    ```
    - `oldRange`: 描述了在**旧文本**中被替换的选区。`index` 是选区开始的位置，`length` 是选区的长度。
    - `newRange`: 描述了在**新文本**中插入内容的选区。
    - 这种格式主要用于处理“选中并替换”的场景。

#### 内部实现 (`find_cursor_edit_diff` 函数)

当 `diff_main` 接收到 `cursor_pos` 参数时，它会调用 `find_cursor_edit_diff` 函数。这个函数会尝试几种常见的编辑模式：

1.  **在光标前插入/删除**:

    - 假设新旧文本的后半部分（从光标位置到末尾）是完全相同的。
    - 然后只比较光标前的部分，快速找出差异。
    - 如果假设成立，就构建差异序列并立即返回，跳过后续所有复杂计算。

2.  **在光标后插入/删除**:

    - 与上面相反，假设新旧文本的前半部分（从开头到光标位置）是完全相同的。
    - 然后只比较光标后的部分。
    - 如果假设成立，就快速返回结果。

3.  **范围替换**:
    - 如果 `cursor_pos` 是一个对象，它会检查旧文本中选区前后的内容是否与新文本的相应部分匹配。
    - 如果匹配，那么差异就只发生在 `oldRange` 和 `newRange` 之间。
    - 它会把 `oldRange` 对应的文本标记为 `DIFF_DELETE`，把 `newRange` 对应的文本标记为 `DIFF_INSERT`，然后返回结果。

如果 `find_cursor_edit_diff` 成功地匹配了以上任何一种模式，它就会返回一个差异数组。如果所有模式都匹配失败（意味着编辑行为比较复杂，不是简单的局部增删改），它会返回 `null`，`diff_main` 就会继续执行后续的常规比较流程（剥离公共前后缀、半匹配、Myers 算法等）。

---

### 2. `cleanup`

#### 作用与目的

`cleanup` 是一个布尔值参数 (`true` 或 `false`)，用于控制是否对差异结果进行**语义化清理（Semantic Cleanup）**。

默认情况下，`diff` 算法（尤其是 Myers 算法）只关心如何用**最少**的编辑步骤（插入和删除）来完成从 `text1`到 `text2` 的转换。这种“数学上最优”的结果有时在人类看来并不自然或不易读。

**`cleanup` 的目的就是让差异结果更符合人类的阅读习惯。**

#### 参数格式

- `cleanup: true`: 启用语义化清理。
- `cleanup: false` (或不提供该参数): 禁用语义化清理。

#### 内部实现 (`diff_cleanupSemantic` 和 `diff_cleanupSemanticLossless` 函数)

当 `cleanup` 为 `true` 时，`diff_main` 会在生成初步的差异结果后，调用 `diff_cleanupSemantic` 函数。这个函数主要做两件事：

1.  **消除琐碎的相等块**:

    - 它会寻找那些被大量编辑包围的、非常短的 `DIFF_EQUAL` 块。
    - 例如，`"abc"` 变成 `"ab123c"`，结果可能是 `[[0, "ab"], [1, "123"], [0, "c"]]`。
    - 如果 `lastequality`（这里是 `"c"`）的长度相对于两边的编辑量来说非常小，算法会认为保留这个微小的相等块没有意义，不如将其视为一次更大的替换。它可能会将结果变成 `[[-1, "abc"], [1, "ab123c"]]`，然后再进行合并清理。这有助于简化差异。

2.  **移动编辑边界以对齐单词（核心功能）**:
    - 这是由 `diff_cleanupSemanticLossless` 函数完成的。
    - 它会寻找 `EQUAL - EDIT - EQUAL` 这样的模式（例如 `[ [0, "The "], [1, "cat "], [0, "came."]]`）。
    - 然后，它会尝试向左或向右“滑动”中间的编辑块（`EDIT`），并为每个可能的位置计算一个“分数”。
    - **评分标准 (`diff_cleanupSemanticScore_`)**:
      - 边界落在**换行符**或**空白行**上，得分最高。
      - 边界落在**空格**或**标点符号**上，得分较高。
      - 边界落在**字母或数字**中间，得分最低。
    - 最后，它会选择使总分最高的方案，将编辑边界移动到最“自然”的位置。

**示例**:

- **没有 `cleanup`**:
  - 比较 `"The cat came."` 和 `"The came."`
  - 结果可能为: `[[0, "The "], [-1, "cat "], [0, "came."]]` (这已经很好了)
- **没有 `cleanup` 的坏情况**:
  - 比较 `"The quick."` 和 `"The fast."`
  - 数学最优解可能是: `[[0, "The q"], [-1, "uic"], [1, "as"], [0, "k."]]`
- **使用 `cleanup: true`**:
  - 算法会发现把编辑边界放在单词中间（`q` 和 `k` 之间）的分数很低。
  - 它会尝试移动边界，发现 `[[0, "The "], [-1, "quick"], [1, "fast"], [0, "."]]` 这个方案的边界都在空格和标点上，得分最高。
  - 最终返回这个更易读的结果。

### 总结

- `cursor_pos`: 一个**性能优化**参数，通过提供编辑位置的提示，让算法尝试走“捷径”，避免对整个文本进行昂贵的比较，特别适用于文本编辑器场景。
- `cleanup`: 一个**可读性优化**参数，通过对差异结果进行后处理，将编辑边界对齐到单词或句子边界，使结果更符合人类的阅读习惯，但会稍微增加计算成本。

---

## fix_unicode

### `_fix_unicode` 参数详解

#### 核心目标：正确处理 Unicode 代理对（Surrogate Pairs）

要理解 `_fix_unicode`，首先需要了解 JavaScript 字符串如何处理复杂的 Unicode 字符。

- **背景**: JavaScript 内部使用 UTF-16 编码来表示字符串。对于大部分常用字符（如英文字母、中文汉字），一个字符就对应一个 16 位的编码单元。
- **问题**: 但是，对于超出这个范围的字符，比如很多 Emoji 表情（如 "😊"）、一些不常用的汉字或特殊符号，需要用**两个**16 位的编码单元来表示。这两个单元被称为**代理对（Surrogate Pair）**。

  - 例如，Emoji "😊" 在 JavaScript 字符串中实际上是两个字符：`"\uD83D"` (高代理) 和 `"\uDE0A"` (低代理)。

- **风险**: `diff` 算法在计算差异时，是逐个“字符”（即 16 位编码单元）进行比较的。这就存在一个风险：算法可能会在代理对的中间将它们拆分。
  - **错误示例**: 比较 `"A😊B"` 和 `"AC"`。一个不正确的 `diff` 结果可能是 `[[0, "A\uD83D"], [-1, "\uDE0AB"], [1, "C"]]`。这里，Emoji "😊" 被拆开了，高代理 `\uD83D` 留在了前面的相等块，低代理 `\uDE0A` 被错误地和 "B" 一起删除了。这破坏了字符的完整性，导致乱码或错误。

**`_fix_unicode` 的核心目标就是检测并修复这种代理对被错误拆分的情况，确保 Unicode 字符的完整性。**

#### 参数与实现

`_fix_unicode` 是一个布尔值，它主要在两个地方体现其作用：

1.  **在顶层函数 `diff` 中被设置为 `true`**

    ```javascript
    function diff(text1, text2, cursor_pos, cleanup) {
      // only pass fix_unicode=true at the top level, not when diff_main is
      // recursively invoked
      return diff_main(text1, text2, cursor_pos, cleanup, true)
    }
    ```

    - 这确保了当你从外部调用库的 `diff` 函数时，Unicode 修复功能总是开启的。
    - 注释也明确说明了，这个参数只在顶层调用时传递 `true`，在递归调用 `diff_main` 时（例如在 `diff_halfMatch_` 或 `diff_bisectSplit_` 内部）不会传递，此时 `_fix_unicode` 的值为 `undefined`（即 `false`）。这样做是为了效率，只在最后的结果上进行一次总的修复，而不是在中间过程中反复修复。

2.  **在 `diff_cleanupMerge` 函数中被使用**
    ```javascript
    function diff_cleanupMerge(diffs, fix_unicode) {
      // ...
      switch (diffs[pointer][0]) {
        case DIFF_EQUAL:
          var previous_equality = ...;
          if (fix_unicode) {
            // ... 核心修复逻辑 ...
          }
          // ...
      }
      // ...
    }
    ```
    - 这是 `_fix_unicode` 发挥作用的核心地带。
    - 当 `diff_cleanupMerge` 函数在合并和清理差异块时，如果 `fix_unicode` 为 `true`，它会执行一段特殊的检查逻辑。

#### 修复逻辑详解

`diff_cleanupMerge` 在遇到一个 `DIFF_EQUAL` 块时，会检查这个块与它前后相邻的编辑块（`INSERT`/`DELETE`）的边界：

1.  **检查前边界**:

    - `if (previous_equality >= 0 && ends_with_pair_start(diffs[previous_equality][1]))`
    - 这行代码检查**前一个** `EQUAL` 块的**最后一个字符**是不是一个**高代理**（即代理对的开始部分）。
    - 如果是，说明代理对在这里被拆分了。它的后半部分（低代理）肯定被错误地归入了接下来的 `INSERT` 或 `DELETE` 块的开头。
    - **修复操作**: 算法会将这个高代理从前一个 `EQUAL` 块中“切掉”，然后加到当前正在处理的 `INSERT` 和 `DELETE` 文本的开头，从而在编辑块内部重新“团聚”这个 Emoji。

2.  **检查后边界**:
    - `if (starts_with_pair_end(diffs[pointer][1]))`
    - 这行代码检查**当前这个** `EQUAL` 块的**第一个字符**是不是一个**低代理**（即代理对的结束部分）。
    - 如果是，说明代理对在这里也被拆分了。它的前半部分（高代理）肯定被错误地归入了前一个 `INSERT` 或 `DELETE` 块的末尾。
    - **修复操作**: 算法会将这个低代理从当前 `EQUAL` 块中“切掉”，然后加到当前正在处理的 `INSERT` 和 `DELETE` 文本的末尾。

通过在两个边界上进行这种“挪动”，算法确保了没有一个完整的 Unicode 字符（如 Emoji）会被差异边界无情地拆散。

### 总结

- `_fix_unicode` 是一个内部参数，用于**开启 Unicode 代理对的修复逻辑**。
- 它解决了 `diff` 算法可能将多字节字符（如 Emoji）从中间拆分，导致结果错误的问题。
- 这个修复操作被设计为在**最后清理合并阶段 (`diff_cleanupMerge`)** 执行一次，而不是在递归的每一步都执行，以保证效率。
- 作为库的使用者，你不需要关心这个参数，因为顶层的 `diff` 函数已经为你默认开启了它，保证了结果的正确性。

---

## find_cursor_edit_diff

### `find_cursor_edit_diff(oldText, newText, cursor_pos)` 函数详解

#### 核心目标与定位

这个函数是 fastDiff.js 库中的一个关键**性能优化**。它的存在是为了**绕过（shortcut）** 完整且耗时的 Myers 差分算法，专门用于快速处理在**文本编辑器**中非常常见的、简单的编辑场景。

当用户在编辑器中进行操作时，绝大多数情况是：

1.  在光标处插入或删除字符。
2.  选中一段文本，然后用新文本替换它。

这些操作的共同点是**局部性**：文本的大部分内容保持不变，只有光标或选区周围发生了变化。`find_cursor_edit_diff` 就是用来识别并快速处理这些局部变化的。如果它能成功处理，就直接返回差异结果，从而避免了对整个文本进行复杂的比较。

#### 参数

- `oldText` (string): 编辑前的原始文本。
- `newText` (string): 编辑后的新文本。
- `cursor_pos` (number | object): 描述编辑位置的“提示信息”。
  - 如果是 `number`，表示旧文本中光标的索引位置（没有选区）。
  - 如果是 `object`，通常包含 `oldRange` 和 `newRange`，描述了编辑前后的选区范围。

#### 算法工作流程

函数内部通过一系列的**假设和验证**来工作，可以分为三个主要场景：

**场景 1 & 2: 在光标处插入或删除 (没有选区)**

这是由第一个 `if (oldRange.length === 0 ...)` 块处理的。它假设用户没有选中任何文本，只是在光标处进行了输入或删除。

它会依次测试两种可能性：

**A. `editBefore:` 块 - 变化发生在光标之前**

- **假设**: 文本的变化完全发生在光标位置**之前**，而光标位置**之后**的所有文本都保持不变。
  - 例如：`oldText`="ab|cde", `newText`="ax|cde" ( `|` 代表光标)
- **验证步骤**:
  1.  根据文本长度的变化，计算出新文本中对应的光标位置 `newCursor`。
  2.  将 `newText` 在 `newCursor` 处分割成 `newBefore` 和 `newAfter`。
  3.  **核心验证**: 检查 `newAfter` 是否严格等于 `oldAfter`。如果不是，说明光标后的文本也变了，假设不成立，立即用 `break editBefore` 跳出这个代码块。
  4.  如果验证通过，说明差异确实只存在于 `oldBefore` 和 `newBefore` 之间。
  5.  函数会进一步找出 `oldBefore` 和 `newBefore` 的公共前缀（这部分是 `DIFF_EQUAL`），剩下的部分就是实际被删除和插入的内容。
  6.  最后调用 `make_edit_splice` 组装成 `[EQUAL, DELETE, INSERT, EQUAL]` 格式的差异数组并返回。

**B. `editAfter:` 块 - 变化发生在光标之后**

- **假设**: 文本的变化完全发生在光标位置**之后**，而光标位置**之前**的所有文本都保持不变。
  - 例如：`oldText`="abc|de", `newText`="abc|xyz"
- **验证步骤**:
  1.  将 `newText` 在 `oldCursor` 处分割成 `newBefore` 和 `newAfter`。
  2.  **核心验证**: 检查 `newBefore` 是否严格等于 `oldBefore`。如果不是，说明光标前的文本也变了，假设不成立，用 `break editAfter` 跳出。
  3.  如果验证通过，说明差异只存在于 `oldAfter` 和 `newAfter` 之间。
  4.  函数会找出 `oldAfter` 和 `newAfter` 的公共后缀，剩下的部分就是实际的差异。
  5.  调用 `make_edit_splice` 组装差异数组并返回。

---

**场景 3: 替换选区内容**

这是由第二个 `if (oldRange.length > 0 ...)` 块处理的。它假设用户**选中了一段文本** (`oldRange.length > 0`)，然后输入了新内容将其替换。

- **假设**: 文本的变化被严格限制在原始选区的范围内，而选区**之前**和**之后**的文本都保持不变。
  - 例如：`oldText`="a[bcd]e", `newText`="a[xyz]e" ( `[]` 代表选区)
- **验证步骤**:
  1.  从 `oldText` 中提取出选区前的前缀 `oldPrefix` 和选区后的后缀 `oldSuffix`。
  2.  从 `newText` 中也提取出相应长度的前缀 `newPrefix` 和后缀 `newSuffix`。
  3.  **核心验证**: 检查 `oldPrefix` 是否等于 `newPrefix` **并且** `oldSuffix` 是否等于 `newSuffix`。如果任何一个不相等，说明变化超出了选区范围，假设不成立，用 `break replaceRange` 跳出。
  4.  如果验证通过，那么差异就非常明确了：
      - 被删除的内容 (`oldMiddle`) 就是 `oldText` 中位于 `oldPrefix` 和 `oldSuffix` 之间的部分。
      - 被插入的内容 (`newMiddle`) 就是 `newText` 中位于 `newPrefix` 和 `newSuffix` 之间的部分。
  5.  调用 `make_edit_splice` 组装差异数组并返回。

---

#### 返回值

- 如果以上任何一个场景的假设被成功验证，函数会返回一个计算好的、格式正确的差异数组。
- 如果所有场景的假设都失败了（意味着这是一个更复杂的、非局部的编辑），函数会返回 `null`。`diff_main` 函数在接收到 `null` 后，就会放弃这个“捷径”，转而执行标准的、更耗时的完整比较流程。

---

## `diff_halfMatch_`

好的，我们来详细讲解 `diff_halfMatch_` 函数。

### `diff_halfMatch_(text1, text2)` 函数详解

#### 核心目标：分而治之的性能优化

`diff_halfMatch_` 是 fastDiff.js 库中一个非常重要的**启发式性能优化**。它的核心思想是“分而治之”（Divide and Conquer）。

在直接使用计算成本高昂的 Myers 算法（`diff_bisect_`）之前，`diff_halfMatch_` 会先尝试一个“捷径”：它会检查两个字符串 `text1` 和 `text2` 是否共享一个**足够长**的公共子串。

- **如果找到了**这样一个公共子串，它就会将原始的大问题（比较 `text1` 和 `text2`）分解为两个小得多的问题：

  1.  比较公共子串**之前**的文本。
  2.  比较公共子串**之后**的文本。
      然后将这两部分的比较结果与中间的公共子串拼接起来，得到最终结果。这通常比直接对整个长文本运行 Myers 算法要快得多。

- **如果没有找到**，函数会返回 `null`，然后 `diff_compute_` 会继续执行后续的 Myers 算法。

**重要提示**：如注释中所说，这个优化可能会产生**非最小**的差异结果。因为它优先考虑了性能，通过找到一个大的公共块来分割问题，但这不一定能保证最终的编辑步骤（增/删）是最少的。但在处理大型文本时，这种性能提升往往是值得的。

#### 算法工作流程

1.  **预检 (Pre-checks)**

    ```javascript
    var longtext = text1.length > text2.length ? text1 : text2
    var shorttext = text1.length > text2.length ? text2 : text1
    if (longtext.length < 4 || shorttext.length * 2 < longtext.length) {
      return null // Pointless.
    }
    ```

    - 首先，确定哪个是长字符串 (`longtext`)，哪个是短字符串 (`shorttext`)。
    - 然后进行快速判断：如果长字符串本身就很短（小于 4 个字符），或者短字符串的长度还不到长字符串的一半，那么找到一个“足够长”（超过一半）的公共子串是不可能的。在这种情况下，继续搜索是“无意义的”（Pointless），直接返回 `null`。

2.  **两次“播种-扩展”尝试**
    算法的核心逻辑在内部函数 `diff_halfMatchI_` 中。`diff_halfMatch_` 会调用它两次，从两个不同的“战略位置”开始尝试寻找公共子串：

    ```javascript
    // 尝试1: 从长字符串的 1/4 处开始
    var hm1 = diff_halfMatchI_(longtext, shorttext, Math.ceil(longtext.length / 4))
    // 尝试2: 从长字符串的 1/2 处（中点）开始
    var hm2 = diff_halfMatchI_(longtext, shorttext, Math.ceil(longtext.length / 2))
    ```

    - 这两次尝试都使用了“播种-扩展”策略：从指定位置取一个“种子”（seed），在短字符串中找到这个种子，然后向前后扩展，看能找到多长的公共部分。
    - 选择 1/4 和 1/2 位置是基于启发式：如果存在一个大的公共块，它很可能与字符串的中心区域重叠。从这两个点开始搜索，有很大概率能快速找到它。

3.  **选择最佳匹配**

    ```javascript
    var hm
    if (!hm1 && !hm2) {
      return null
    } else if (!hm2) {
      hm = hm1
    } else if (!hm1) {
      hm = hm2
    } else {
      // Both matched.  Select the longest.
      hm = hm1[4].length > hm2[4].length ? hm1 : hm2
    }
    ```

    - 如果两次尝试都失败了（都返回 `null`），说明没有找到满足条件的半匹配，函数返回 `null`。
    - 如果只有一次成功，就用那次的结果。
    - 如果两次都成功了，就比较它们找到的公共子串的长度（存储在返回数组的第 5 个元素 `hm[4]` 中），选择**更长**的那个作为最终的匹配结果 `hm`。

4.  **整理并返回结果**
    ```javascript
    // A half-match was found, sort out the return data.
    var text1_a, text1_b, text2_a, text2_b
    if (text1.length > text2.length) {
      // ...
    } else {
      // ...
    }
    var mid_common = hm[4]
    return [text1_a, text1_b, text2_a, text2_b, mid_common]
    ```
    - 一旦确定了最佳匹配 `hm`，就需要将 `hm` 中相对于 `longtext` 和 `shorttext` 的分割结果，重新映射回原始的 `text1` 和 `text2`。
    - 最终，函数返回一个包含 5 个元素的数组：
      - `[0]`: `text1` 在公共子串之前的部分 (`text1_a`)
      - `[1]`: `text1` 在公共子串之后的部分 (`text1_b`)
      - `[2]`: `text2` 在公共子串之前的部分 (`text2_a`)
      - `[3]`: `text2` 在公共子串之后的部分 (`text2_b`)
      - `[4]`: 找到的公共子串本身 (`mid_common`)

调用者 `diff_compute_` 接收到这个数组后，就会递归地调用 `diff_main` 来分别比较 `(text1_a, text2_a)` 和 `(text1_b, text2_b)`，从而完成分治过程。

---

## `diff_bisect_`

好的，我们来详细讲解 `diff_bisect_` 函数。这是整个 `diff` 算法库中最为核心和复杂的函数。

### `diff_bisect_(text1, text2)` 函数详解

#### 核心目标：Myers 差分算法的实现

`diff_bisect_` 是著名的 **Myers 差分算法**（1986 年）的一种高效实现。当所有更快的启发式方法（如 `find_cursor_edit_diff`, `diff_halfMatch_`）都失败后，它作为**最终的、最可靠的手段**被调用。

它的目标是找到两个字符串之间的**最优差异**，即生成一个**最短的编辑脚本**（Shortest Edit Script），意味着最少的插入和删除操作。

#### 核心思想：寻找“中间蛇” (Middle Snake)

要理解这个算法，首先需要想象一个“编辑图”（Edit Graph）：一个二维网格，`text1` 的长度是宽，`text2` 的长度是高。从左上角走到右下角的一条路径就代表一种编辑方案。

- 向右走一步：删除 `text1` 的一个字符。
- 向下走一步：插入 `text2` 的一个字符。
- 沿对角线走一步：`text1` 和 `text2` 的一个字符匹配（相等）。

Myers 算法的目标就是找到一条经过最多对角线（匹配）的路径。

直接从头找到尾的计算量很大。`diff_bisect_` 采用了 Myers 算法的一个关键优化：**双向搜索并寻找“中间蛇”**。

1.  **双向搜索**：算法不再只从左上角（起点）向右下角（终点）搜索。而是**同时**：

    - 从**起点**开始，向前搜索。
    - 从**终点**开始，向后搜索。

2.  **寻找交点**：当向前搜索的路径和向后搜索的路径在图的中间区域**相遇或重叠**时，算法就找到了一个分割点（“中间蛇”）。

3.  **分而治之**：一旦找到这个分割点，原始的大问题就被一分为二，变成了两个更小的 `diff` 问题。然后算法递归地解决这两个小问题。

**为什么这样更快？** 搜索的复杂度与编辑距离 `D` 相关（大约是 O(ND)）。通过双向搜索，我们只需要分别搜索 `D/2` 的距离，而不是完整的 `D`。这会使搜索空间大大减小，从而显著提升性能。

#### 算法工作流程

让我们将上述思想映射到代码上：

1.  **初始化**

    ```javascript
    var text1_length = text1.length
    var text2_length = text2.length
    var max_d = Math.ceil((text1_length + text2_length) / 2)
    var v_offset = max_d
    var v1 = new Array(2 * max_d) // 向前搜索的前沿
    var v2 = new Array(2 * max_d) //向后搜索的前沿
    // ... 初始化 v1, v2 为 -1 ...
    v1[v_offset + 1] = 0
    v2[v_offset + 1] = 0
    ```

    - `v1` 和 `v2` 是两个核心数组，分别记录了向前和向后搜索的“前沿阵地”。`v[k]` 存储的是在对角线 `k` 上能达到的最远 `x` 坐标。
    - `v_offset` 用于处理 `k` 可能是负数的情况，通过偏移量将数组索引映射到正数范围。
    - `max_d` 是理论上可能的最大编辑距离，作为循环的上限。

2.  **主循环 (逐层扩展搜索)**

    ```javascript
    for (var d = 0; d < max_d; d++) {
      // ...
    }
    ```

    - `d` 代表当前的编辑距离（或搜索的“步数”）。循环的每一次都代表向前和向后各“走”一步。

3.  **向前走一步 (Walk the front path)**

    ```javascript
    for (var k1 = -d + k1start; k1 <= d - k1end; k1 += 2) {
        // ... 计算 x1, y1 ...
        while (x1 < text1_length && y1 < text2_length && text1.charAt(x1) === text2.charAt(y1)) {
            x1++;
            y1++;
        }
        v1[k1_offset] = x1;
        // ... 检查是否与 v2 重叠 ...
        if (front && ... v2[k2_offset] !== -1) {
            // ... Overlap detected.
            return diff_bisectSplit_(text1, text2, x1, y1);
        }
    }
    ```

    - 内层循环遍历所有在 `d` 步内可以到达的对角线 `k1`。
    - 在每个 `k1` 上，它先通过一次非对角线移动（插入或删除）到达一个新的点，然后尽可能地沿着对角线“滑行”（`while` 循环），处理所有连续的匹配字符。这个“滑行”就是所谓的“蛇”。
    - 更新 `v1` 数组，记录在 `k1` 这条线上最远到达的位置。
    - **关键**：每次更新后，立即检查这个新位置是否已经被 `v2`（反向路径）访问过。如果访问过，说明两条路径相遇了！立即调用 `diff_bisectSplit_` 并返回。

4.  **向后走一步 (Walk the reverse path)**

    ```javascript
    for (var k2 = -d + k2start; k2 <= d - k2end; k2 += 2) {
        // ... 逻辑与向前走一步类似，但比较的是反向的字符 ...
        // ... 检查是否与 v1 重叠 ...
        if (!front && ... v1[k1_offset] !== -1) {
            // ... Overlap detected.
            return diff_bisectSplit_(text1, text2, x1, y1);
        }
    }
    ```

    - 这部分逻辑与向前走一步完全对称，但是它是从字符串的末尾开始，反向进行比较。
    - 同样，在每一步之后，它都会检查是否与 `v1`（正向路径）相遇。

5.  **分治与递归 (`diff_bisectSplit_`)**
    当 `diff_bisect_` 找到重叠点 `(x, y)` 后，它会调用 `diff_bisectSplit_`。

    ```javascript
    function diff_bisectSplit_(text1, text2, x, y) {
      var text1a = text1.substring(0, x)
      var text2a = text2.substring(0, y)
      var text1b = text1.substring(x)
      var text2b = text2.substring(y)

      var diffs = diff_main(text1a, text2a)
      var diffsb = diff_main(text1b, text2b)

      return diffs.concat(diffsb)
    }
    ```

    - 这个函数的作用非常清晰：它根据找到的分割点 `(x, y)` 将两个字符串都切成两半。
    - 然后递归地调用 `diff_main` 来分别计算前半部分的差异和后半部分的差异。
    - 最后将两个结果合并，完成整个 `diff` 计算。

6.  **最坏情况**
    如果循环结束都没有找到重叠点（这在理论上可能发生于超时或文本完全无关联），则返回一个表示完全替换的结果。

### 总结

`diff_bisect_` 是一个优雅而强大的算法。它不是直接解决整个问题，而是通过**双向搜索**高效地找到问题的“腰部”（中间蛇），然后将问题**一分为二**，递归解决。这种分治策略是它能够在合理时间内处理复杂差异的关键。它是整个库的基石，保证了在所有优化都失效时，依然能得到一个数学上最优的解。

---

## `diff_cleanupMerge`

好的，我们来详细讲解 `diff_cleanupMerge` 函数。

### `diff_cleanupMerge(diffs, fix_unicode)` 函数详解

#### 核心目标：规范化和优化差异列表

`diff_cleanupMerge` 是一个后处理函数，它的任务不是计算差异，而是对已经生成的原始差异列表（`diffs` 数组）进行**清理、合并和优化**。

一个原始的 `diff` 结果可能是零碎和冗余的，例如：`[[DELETE, "a"], [DELETE, "b"], [INSERT, "c"], [INSERT, "d"]]`。`diff_cleanupMerge` 会将这样的结果变得更紧凑、更高效，例如合并为 `[[DELETE, "ab"], [INSERT, "cd"]]`。

它主要执行以下几个关键任务：

1.  **合并同类项**：将连续的 `INSERT`、`DELETE` 或 `EQUAL` 块合并成一个。
2.  **提取公共部分**：在一组删除和插入操作之间，提取出公共的前后缀，将它们转换成 `EQUAL` 块。
3.  **修复 Unicode 代理对**：确保多字节字符（如 Emoji）不会被差异边界拆分。
4.  **移动编辑块**：通过移动编辑块来消除一些不必要的 `EQUAL` 块，使结果更简洁。

#### 算法工作流程

这个函数分两个主要阶段（Pass）进行处理。

---

### 第一阶段：主循环合并与提取

这是函数的主体部分，通过一次遍历来完成大部分合并工作。

1.  **初始化与遍历**

    ```javascript
    diffs.push([DIFF_EQUAL, '']) // 在末尾添加一个哨兵
    var pointer = 0
    var count_delete = 0,
      count_insert = 0
    var text_delete = '',
      text_insert = ''
    while (pointer < diffs.length) {
      // ...
    }
    ```

    - 算法首先在 `diffs` 数组末尾添加一个空的 `EQUAL` 块作为“哨兵”。这确保了循环的最后一次也能正确处理之前累积的 `INSERT`/`DELETE` 块。
    - 它使用一个 `pointer` 遍历 `diffs` 数组。
    - 当遇到 `INSERT` 或 `DELETE` 时，它不会立即处理，而是将它们的内容累加到 `text_insert` 和 `text_delete` 变量中，并用 `count_insert` 和 `count_delete` 记录块的数量。

2.  **遇到 `EQUAL` 块时进行处理**
    当 `switch` 语句遇到 `case DIFF_EQUAL:` 时，意味着一连串的 `INSERT`/`DELETE` 序列结束了，此时开始进行集中处理：

    - **a. Unicode 修复 (`if (fix_unicode)`)**:

      - 这是非常重要的一步。它检查累积的编辑块与前后 `EQUAL` 块的边界，看是否有 Unicode 代理对（如 Emoji）被拆散。
      - 如果前一个 `EQUAL` 块以高代理结尾，或当前 `EQUAL` 块以低代理开头，它会把被拆散的部分“挪”回 `text_insert` 和 `text_delete` 中，以保证字符的完整性。

    - **b. 提取公共前缀/后缀**:

      ```javascript
      // Factor out any common prefixes.
      commonlength = diff_commonPrefix(text_insert, text_delete)
      // ...
      // Factor out any common suffixes.
      commonlength = diff_commonSuffix(text_insert, text_delete)
      // ...
      ```

      - 这是优化的核心。算法比较累积的 `text_delete` 和 `text_insert`。
      - **示例**: 如果 `text_delete` 是 `"apple pie"`，`text_insert` 是 `"apple tart"`。
        - `diff_commonPrefix` 会发现公共前缀 `"apple "`。这个前缀会被从 `text_delete` 和 `text_insert` 中移除，并追加到**前一个** `EQUAL` 块的末尾。
        - `diff_commonSuffix` 会发现公共后缀（在这个例子中没有）。如果有，它会被从 `text_delete` 和 `text_insert` 中移除，并添加到**当前这个** `EQUAL` 块的开头。
      - 经过这一步，`text_delete` 变为 `"pie"`，`text_insert` 变为 `"tart"`，而 `"apple "` 成为了公共部分。

    - **c. 替换旧的 `diff` 块**:

      - 在提取完公共部分后，算法用处理后剩下的 `text_delete` 和 `text_insert` 生成新的 `[DELETE, ...]` 和 `[INSERT, ...]` 块。
      - 然后，它使用 `splice` 方法，将原始的、零碎的 `INSERT`/`DELETE` 序列（共 `n` 个块）替换为这 1 个或 2 个新的、紧凑的块。

    - **d. 合并 `EQUAL` 块**:
      - 在完成上述替换后，如果新的 `diff` 序列中出现了连续的 `EQUAL` 块，它们会被合并成一个。

---

### 第二阶段：移动编辑块以消除 `EQUAL`

在第一轮大合并之后，可能会出现一些可以进一步优化的情况。

```javascript
// Second pass: look for single edits surrounded on both sides by equalities
// which can be shifted sideways to eliminate an equality.
// e.g: A<ins>BA</ins>C -> <ins>AB</ins>AC
var changes = false
pointer = 1
while (pointer < diffs.length - 1) {
  if (diffs[pointer - 1][0] === DIFF_EQUAL && diffs[pointer + 1][0] === DIFF_EQUAL) {
    // ...
  }
}
```

- 这个阶段专门寻找 `[EQUAL, EDIT, EQUAL]` 这样的模式。
- **示例**: `[EQUAL, "A"], [INSERT, "BA"], [EQUAL, "C"]`
  - 它检查中间的 `INSERT` 块 `"BA"` 是否以**前一个** `EQUAL` 块的内容 `"A"` 结尾。
  - 这里不匹配。然后检查 `INSERT` 块 `"BA"` 是否以**后一个** `EQUAL` 块的内容 `"C"` 开头。
  - 也不匹配。
- **另一个示例**: `[EQUAL, "A"], [INSERT, "BC"], [EQUAL, "B"]`
  - 它检查 `INSERT` 块 `"BC"` 的开头是否与后一个 `EQUAL` 块 `"B"` 匹配。
  - 匹配！算法会进行“移位”：
    - 将后一个 `EQUAL` 块的内容 `"B"` 追加到前一个 `EQUAL` 块，使其变为 `[EQUAL, "AB"]`。
    - 将 `INSERT` 块的内容也进行相应移动，变为 `[INSERT, "CB"]`。
    - 删除后一个 `EQUAL` 块。
  - 结果就从 `[EQUAL, "A"], [INSERT, "BC"], [EQUAL, "B"]` 变成了 `[EQUAL, "AB"], [INSERT, "CB"]`。虽然这个例子看起来没优化，但对于 `A<ins>BA</ins>C` -> `<ins>AB</ins>AC` 这样的情况，它可以消除一个 `EQUAL` 块，使 `diff` 更简洁。

#### 递归调用

如果在第二阶段发生了任何“移位”，`changes` 标志会变为 `true`。函数最后会检查这个标志，如果为 `true`，它会**再次调用自己** (`diff_cleanupMerge(diffs, fix_unicode)`)，因为移位操作可能创造了新的可合并的机会。

### 总结

`diff_cleanupMerge` 是一个强大的规范化引擎。它通过两轮细致的清理工作，将 `diff` 算法的原始输出转化为一种更紧凑、更高效、更具逻辑性的最终形式。它是保证 `diff` 结果质量的关键步骤。

---

## `diff_cleanupSemantic`
