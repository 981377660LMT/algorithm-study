好的，我们来详细讲解 Google 的 `diff-match-patch` 这个经典库的未压缩版本，并提供具体的使用案例。

### `diff-match-patch.js` 详细讲解

#### 一、 总体概述

`diff-match-patch` 是由 Google 的 Neil Fraser 开发的一个非常著名的高性能库，用于处理纯文本的比较和同步。它主要包含三个核心功能，分别对应其名称：

1.  **Diff (差异比较)**：比较两个文本，并生成一个描述两者差异的“编辑脚本”（diffs）。
2.  **Match (模糊匹配)**：在一个长文本中，以给定的位置为中心，模糊地查找一个短文本模式（pattern）的最佳匹配位置。
3.  **Patch (补丁应用)**：根据 `diff` 的结果生成“补丁”（patches），并将这些补丁应用到另一个文本上，从而实现文本的更新或同步。

这个库被广泛应用于各种需要文本比较的场景，如代码版本控制（类似 Git 的 diff 功能）、在线协作编辑器（如 Google Docs）、富文本编辑器的历史记录（undo/redo）等。

#### 二、 核心功能详解

##### 1. Diff (差异比较)

这是库最核心、最复杂的部分。它的目标是找出从 `text1` 转换到 `text2` 所需的最少编辑步骤。

**主要入口函数**: `diff_main(text1, text2, opt_checklines)`

**工作流程**:
`diff_main` 的执行过程是一个精心设计的多阶段优化流程，旨在尽可能快地处理不同类型的文本差异。

1.  **预处理（性能优化）**:

    - **检查完全相等**：如果两个文本完全相同，直接返回一个 `EQUAL` 块。
    - **剥离公共前后缀**：调用 `diff_commonPrefix` 和 `diff_commonSuffix` 快速找到并分离出两个文本开头和结尾完全相同的部分。这是一个巨大的性能提升，因为它将一个大问题缩小为只比较中间不同部分的小问题。

2.  **核心计算 (`diff_compute_`)**:
    对剥离掉公共部分后的文本，`diff_compute_` 会按顺序尝试一系列策略：

    - **简单情况**：如果一个文本为空，则另一个文本就是纯粹的插入或删除。
    - **子串检查**：检查短文本是否是长文本的子串。如果是，可以快速生成 diff。
    - **半匹配启发式 (`diff_halfMatch_`)**：尝试寻找一个足够长的公共子串（“锚点”），将问题一分为二，递归解决。这是一种不保证最优解但速度极快的启发式优化。
    - **行模式 (`diff_lineMode_`)**：对于非常长的文本（默认超过 100 个字符），它会先将文本按行切分，用单个字符代表每一行，进行一次快速的“行 diff”。然后再对有差异的行进行精确的“字符 diff”。这大大降低了处理大段代码或文章时的计算复杂度。
    - **二分法 (`diff_bisect_`)**：如果以上所有快速方法都失败了，就启动最终的、最精确的 **Myers 差分算法**。该算法能保证找到数学上的最优解（最短编辑距离），但计算成本也最高。

3.  **后处理（结果美化）**:
    - **合并结果**：将预处理阶段剥离的公共前后缀重新加回到 diff 结果的两端。
    - **清理与合并 (`diff_cleanupMerge`)**：合并连续的相同操作（如两个 `DELETE` 合并成一个），并提取 `INSERT` 和 `DELETE` 块之间的公共部分，使结果更紧凑。
    - **语义化清理 (`diff_cleanupSemantic`)**：让 diff 结果更符合人类阅读习惯。例如，将编辑边界从单词中间移动到单词之间的空格上。
      - **原始**: `The c<ins>at c</ins>ame.`
      - **清理后**: `The <ins>cat </ins>came.`

**输出格式**:
Diff 的结果是一个元组数组，每个元组包含 `[操作, 文本]`。

- `[DIFF_DELETE, "Hello"]` 或 `[-1, "Hello"]`: 删除 "Hello"
- `[DIFF_INSERT, "Goodbye"]` 或 `[1, "Goodbye"]`: 插入 "Goodbye"
- `[DIFF_EQUAL, " world"]` 或 `[0, " world"]`: 保留 " world"

---

##### 2. Match (模糊匹配)

这个功能用于在一个大文本中寻找一个模式串的最佳位置，即使没有完美匹配。

**主要入口函数**: `match_main(text, pattern, loc)`

- `text`: 被搜索的长文本。
- `pattern`: 要查找的模式串。
- `loc`: 一个“期望位置”，算法会优先在该位置附近查找。

**工作流程**:
它主要基于 **Bitap 算法** (`match_bitap_`)，这是一种高效的模糊字符串搜索算法。它会计算一个综合“分数”，该分数同时考虑了**匹配的精确度**（有多少字符错误）和**位置的接近度**（匹配位置离 `loc` 有多远）。最终返回一个分数最低（即最佳）的匹配起始索引。如果找不到任何可接受的匹配，则返回 `-1`。

---

##### 3. Patch (补丁)

Patch 功能建立在 Diff 的基础上，用于生成和应用补丁，非常适合数据同步场景。

**主要函数**:

- `patch_make(text1, diffs)`: 接收原始文本和 `diffs` 数组，生成一个“补丁对象”数组。每个补丁对象不仅包含 diff 信息，还包含了它应该被应用在原始文本的哪个位置，以及一些上下文信息（即变化区域前后的不变文本），用于校验位置的正确性。

- `patch_toText(patches)`: 将补丁对象数组转换成人类可读的文本格式（类似 `git format-patch`）。

- `patch_fromText(text)`: 从文本格式解析出补丁对象数组。

- `patch_apply(patches, text)`: 将补丁应用到给定的文本上。它会利用补丁中的上下文信息去查找应用位置。即使文本已经发生了一些轻微变化，只要上下文还能匹配上，它也能智能地、正确地应用补丁。该函数返回一个数组 `[新文本, 应用成功与否的布尔值数组]`。

#### 三、 使用案例

```javascript
// 引入库 (假设已在 HTML 中通过 <script> 标签引入)
// 或者在 Node.js 中: const diff_match_patch = require('./diff_match_patch_uncompressed.js');

// 1. 创建实例
var dmp = new diff_match_patch()

// ==================================================
// 案例 1: Diff (差异比较)
// ==================================================
console.log('--- 案例 1: Diff ---')
var text1 = 'The quick brown fox jumps over the lazy dog.'
var text2 = 'That quick brown fox jumped over a lazy dog.'

var diffs = dmp.diff_main(text1, text2)
console.log('原始 Diff 结果:', JSON.stringify(diffs))

// 对结果进行美化，使其更易读
dmp.diff_cleanupSemantic(diffs)
console.log('语义优化后 Diff 结果:', JSON.stringify(diffs))

// 将 diff 结果转换为漂亮的 HTML
var html = dmp.diff_prettyHtml(diffs)
console.log('HTML 可视化结果:')
console.log(html)
// 你可以将这个 html 字符串插入到网页的 <div> 中来显示差异
// 结果: <del style="background:#ffe6e6;">The</del><ins style="background:#e6ffe6;">That</ins><span> quick brown fox jump</span><del style="background:#ffe6e6;">s</del><ins style="background:#e6ffe6;">ed</ins><span> over </span><del style="background:#ffe6e6;">the</del><ins style="background:#e6ffe6;">a</ins><span> lazy dog.</span>

// ==================================================
// 案例 2: Match (模糊匹配)
// ==================================================
console.log('\n--- 案例 2: Match ---')
var longText = 'The quick brown fox jumps over the lazy dog.'
var pattern = 'jumps'
var expectedLoc = 20 // "jumps" 实际在 20

// 精确匹配
var matchLoc1 = dmp.match_main(longText, pattern, expectedLoc)
console.log(`在位置 ${expectedLoc} 附近查找 "${pattern}": 找到于 ${matchLoc1}`) // 输出: 20

// 模糊匹配
var fuzzyPattern = 'jumped'
var matchLoc2 = dmp.match_main(longText, fuzzyPattern, expectedLoc)
console.log(`在位置 ${expectedLoc} 附近模糊查找 "${fuzzyPattern}": 找到于 ${matchLoc2}`) // 输出: 20 (因为 "jumps" 是最佳匹配)

// ==================================================
// 案例 3: Patch (制作与应用补丁)
// ==================================================
console.log('\n--- 案例 3: Patch ---')
var originalText = 'The rain in Spain stays mainly in the plain.'
var modifiedText = 'The rain in Spain falls mainly on the plain.'

// 步骤 1: 生成 diff
var diffsForPatch = dmp.diff_main(originalText, modifiedText)

// 步骤 2: 基于 diff 和原始文本制作补丁
var patches = dmp.patch_make(originalText, diffsForPatch)
console.log('生成的补丁数量:', patches.length)

// (可选) 步骤 3: 将补丁序列化为文本，便于存储或传输
var patchText = dmp.patch_toText(patches)
console.log('文本格式的补丁:\n', patchText)
/*
输出类似:
@@ -15,15 +15,14 @@
 ain in Spain 
-stays
+falls
  mainly on
*/

// 步骤 4: 将补丁应用到原始文本上
var results = dmp.patch_apply(patches, originalText)
var newText = results[0]
var appliedStatus = results[1]

console.log('补丁应用是否成功:', appliedStatus[0]) // true
console.log('应用补丁后的新文本:', newText)
console.log('新文本是否与目标文本一致:', newText === modifiedText) // true

// 模拟一个稍微被修改过的原始文本
var slightlyModifiedOriginal = 'The gentle rain in Spain stays mainly in the plain.'
var results2 = dmp.patch_apply(patches, slightlyModifiedOriginal)
console.log('\n在轻微修改过的文本上应用补丁是否成功:', results2[1][0]) // 仍然是 true
console.log('应用后的结果:', results2[0]) // "The gentle rain in Spain falls mainly on the plain."
// 注意：即使原始文本加了 "gentle"，补丁依然能通过上下文 "ain in Spain " 找到正确位置并成功应用。
```
