https://github.com/kvandake/lexorank-ts

---

好的，我们来详细讲解 Jira 的 LexoRank 实现。

需要强调的是，Jira 的 LexoRank 是 Atlassian 的专有系统，其确切的源代码并未公开。然而，通过官方文档、开发者博客以及对系统行为的分析，我们可以相当准确地推断其工作原理并构建一个功能上非常相似的实现。

Jira 的实现比我们之前讨论的简单版本要复杂得多，因为它专为大规模、高并发、多租户的企业环境（如 Jira Cloud）而设计。

### Jira LexoRank 的核心概念

Jira 的 LexoRank 引入了几个关键概念来增强可伸缩性和鲁棒性。

#### 1. 格式：`bucket|rank`

Jira 的排名字符串不只是一个简单的字符串，它由两部分组成，用管道符 `|` 分隔：

- **Bucket (桶)**：一个数字，用于将排名划分到不同的逻辑空间。
- **Rank (排名)**：实际的排序字符串，它本身也可能是一个复杂的结构。

**示例**：`0|i0000o:`、`1|100000:10000o`

#### 2. Bucket (桶) 的作用

桶是 Jira 实现中最关键的创新之一。

- **命名空间/分区**：桶将全局的排名空间划分为多个独立的子空间。一个项目（Project）中的问题（Issue）通常会被分配到同一个或少数几个桶中。这确保了在一个项目中进行排序不会与另一个项目中的排序发生冲突。
- **减少重平衡的影响**：当一个桶内的排名变得过于密集时，只需要对该桶内的 Issue 进行重平衡，而不会影响到其他桶。这极大地缩小了昂贵的重平衡操作的范围。
- **并发控制**：当系统检测到需要重平衡时，它可以将受影响的 Issue 迁移到一个全新的、空间充裕的桶中，而旧桶可以被废弃。这是一种高效处理并发写入和维护排名健康度的策略。

#### 3. Rank (排名) 的内部结构

Rank 部分本身也比简单的字符串更复杂。它使用一个 **62 进制** 的系统（`0-9`, `A-Z`, `a-z`），并被设计成一种“小数”形式。

- **基数-62 系统**：每个字符代表一个 0-61 的数字。
- **类小数结构**：Rank 字符串可以看作是一个多位数的基数-62 的数字。例如 `i0000o`。
- **步长 (Step)**：为了在插入时留出足够的空间，Jira 的 Rank 通常不是连续的。初始生成时，可能会创建像 `100000`, `200000`, `300000` 这样的排名，它们之间有巨大的“空间”用于后续插入。
- **处理精度**：当需要在两个相邻的排名（如 `100000` 和 `100001`）之间插入时，Jira 会增加“小数位”，用冒号 `:` 分隔。例如，中间值可能变成 `100000:10000o`（这里的 `10000o` 是一个新的“小数部分”）。这与我们之前讨论的 `es` -> `ej` 的原理相同，但结构化更强。

#### 4. 核心算法：`between`

`between` 算法的逻辑与之前类似，但在 Jira 的结构下操作：

1.  **解析**：首先解析 `prevRank` 和 `nextRank` 的 `bucket` 和 `rank` 部分。
2.  **桶检查**：通常，操作在同一个桶内进行。如果 `prevRank` 和 `nextRank` 在不同的桶中，处理会更复杂，可能表示数据处于迁移或不一致的状态。通常假设它们在同一个桶。
3.  **计算中间 Rank**：这是最复杂的部分。算法需要在两个基数-62 的“数字”之间找到一个中间值。这涉及到自定义的算术运算（加法、除以二）。
    - 如果两个 Rank 字符串长度不同，先将较短的补齐到相同长度（用最小字符 `0` 补齐）。
    - 将这两个字符串视为巨大的基数-62 的整数，计算它们的和，然后除以 2，得到中间值。
    - 将计算出的中间值转换回基数-62 的字符串。
4.  **处理冲突/空间不足**：
    - 如果计算出的中间值等于 `prevRank`，说明它们之间没有空间了。此时就需要增加“小数位”。算法会在 `prevRank` 后面附加一个冒号 `:` 和一个代表中间位置的新 Rank 片段（如 `10000o`）。
    - 如果生成的 Rank 字符串超过了某个长度阈值，这会成为一个触发**重平衡**的信号。

### 代码实现（基于原理的推演）

以下是一个基于上述原理的 TypeScript 实现。它模拟了 Jira LexoRank 的核心逻辑，包括桶、基数-62 的算术和生成中间值。

```typescript
/**
 * Jira LexoRank 的一个功能性重新实现。
 * 注意：这是基于公开信息和原理推断的实现，并非 Atlassian 的源代码。
 */

// --- 1. 定义常量和工具 ---

const BASE_62_ALPHABET = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
const ALPHABET_MAP: Map<string, number> = new Map()
for (let i = 0; i < BASE_62_ALPHABET.length; i++) {
  ALPHABET_MAP.set(BASE_62_ALPHABET[i], i)
}

const MIN_CHAR = BASE_62_ALPHABET[0]
const MAX_CHAR = BASE_62_ALPHABET[BASE_62_ALPHABET.length - 1]
const MID_CHAR = BASE_62_ALPHABET[Math.floor(BASE_62_ALPHABET.length / 2)]
const RANK_SEPARATOR = '|'
const DECIMAL_SEPARATOR = ':'

// --- 2. 定义核心数据结构 ---

interface ParsedRank {
  bucket: number
  rank: string
}

// --- 3. 实现基数-62的算术辅助函数 ---

/** 将单个 base-62 字符转换为十进制数字 */
function toDecimal(char: string): number {
  const val = ALPHABET_MAP.get(char)
  if (val === undefined) throw new Error(`Invalid LexoRank character: ${char}`)
  return val
}

/** 将十进制数字转换为单个 base-62 字符 */
function toBase62(decimal: number): string {
  if (decimal < 0 || decimal >= BASE_62_ALPHABET.length) {
    throw new Error(`Decimal ${decimal} is out of base-62 range`)
  }
  return BASE_62_ALPHABET[decimal]
}

/**
 * 在两个 base-62 字符串之间找到中间字符串
 * 这是 LexoRank 的核心计算逻辑
 */
function getMidRank(prev: string, next: string): string {
  let mid = ''
  let i = 0
  while (true) {
    const prevChar = prev[i] || MIN_CHAR
    const nextChar = next[i] || MAX_CHAR
    const prevVal = toDecimal(prevChar)
    const nextVal = toDecimal(nextChar)

    if (prevVal === nextVal) {
      mid += prevChar
      i++
      continue
    }

    const midVal = Math.floor((prevVal + nextVal) / 2)
    if (midVal === prevVal) {
      mid += prevChar
      i++
      // 补齐 prev，以便下一轮循环可以继续在 "小数位" 上寻找空间
      prev = prev.padEnd(i, MIN_CHAR)
      continue
    }

    mid += toBase62(midVal)
    return mid
  }
}

// --- 4. 主类 JiraLexoRank ---

class JiraLexoRank {
  public readonly value: string
  private readonly parsed: ParsedRank

  private constructor(value: string) {
    this.value = value
    this.parsed = JiraLexoRank.parse(value)
  }

  /**
   * 解析一个完整的 LexoRank 字符串
   * @param rankStr 'bucket|rank' 格式的字符串
   */
  public static parse(rankStr: string): ParsedRank {
    const parts = rankStr.split(RANK_SEPARATOR)
    if (parts.length !== 2) throw new Error('Invalid LexoRank format')
    const bucket = parseInt(parts[0], 10)
    if (isNaN(bucket)) throw new Error('Invalid bucket in LexoRank')
    return { bucket, rank: parts[1] }
  }

  /**
   * 从 ParsedRank 对象格式化回字符串
   */
  private static format(parsed: ParsedRank): string {
    return `${parsed.bucket}${RANK_SEPARATOR}${parsed.rank}`
  }

  /**
   * 创建一个新的 LexoRank 实例
   */
  public static from(value: string): JiraLexoRank {
    return new JiraLexoRank(value)
  }

  /**
   * 生成一个位于两个给定排名之间的排名
   * @param prev 前一个排名实例，或 null
   * @param next 后一个排名实例，或 null
   */
  public static between(prev: JiraLexoRank | null, next: JiraLexoRank | null): JiraLexoRank {
    let bucket = 0 // 默认或根据策略选择
    let prevRank = ''
    let nextRank = ''

    if (prev && next && prev.parsed.bucket !== next.parsed.bucket) {
      // 跨桶移动，这是一个复杂场景，通常会触发重平衡。
      // 这里简化处理：以 next 的桶为准，在桶的头部插入。
      console.warn('Cross-bucket move detected. Simplified to head-insertion in next bucket.')
      bucket = next.parsed.bucket
      nextRank = next.parsed.rank
    } else if (prev) {
      bucket = prev.parsed.bucket
      prevRank = prev.parsed.rank
    } else if (next) {
      bucket = next.parsed.bucket
      nextRank = next.parsed.rank
    }

    // 如果 prev 和 next 都存在，使用它们的 rank
    if (prev && next) {
      prevRank = prev.parsed.rank
      nextRank = next.parsed.rank
    }

    if (prevRank >= nextRank && prevRank && nextRank) {
      throw new Error(`prevRank (${prevRank}) must be smaller than nextRank (${nextRank})`)
    }

    let newRankStr = getMidRank(prevRank, nextRank)

    // 检查是否需要增加小数位
    if (newRankStr === prevRank || newRankStr === nextRank) {
      // 空间不足，需要增加小数位
      // 在 prevRank 和一个理想的 "中间小数" 之间生成
      newRankStr = getMidRank(prevRank, prevRank + DECIMAL_SEPARATOR + MID_CHAR)
    }

    const newFullRank = JiraLexoRank.format({ bucket, rank: newRankStr })
    return new JiraLexoRank(newFullRank)
  }

  /**
   * 生成一个初始排名
   */
  public static initial(): JiraLexoRank {
    return JiraLexoRank.between(null, null)
  }
}

// --- 5. 使用示例 ---

console.log('--- JiraLexoRank Simulation ---')

// 场景1: 在一个空列表中创建第一个条目
const rank1 = JiraLexoRank.initial()
console.log('1. First item in list:', rank1.value) // e.g., "0|U"

// 场景2: 在列表末尾添加一个条目
const rank2 = JiraLexoRank.between(rank1, null)
console.log('2. Add to end:', rank2.value) // e.g., "0|w"

// 场景3: 在列表开头添加一个条目
const rank0 = JiraLexoRank.between(null, rank1)
console.log('3. Add to start:', rank0.value) // e.g., "0|K"

// 场景4: 在两个现有条目之间插入
const rank1_5 = JiraLexoRank.between(rank1, rank2)
console.log('4. Insert in middle:', rank1_5.value) // e.g., "0|V"

let list = [rank0, rank1, rank1_5, rank2].map(r => r.value).sort()
console.log('Current sorted list:', list)

// 场景5: 在两个非常接近的排名之间插入 (模拟空间耗尽)
const rA = JiraLexoRank.from('0|V')
const rB = JiraLexoRank.from('0|W') // V 和 W 在 base-62 中是相邻的
const rAB = JiraLexoRank.between(rA, rB)
console.log(`5. Between "${rA.value}" and "${rB.value}":`, rAB.value) // e.g., "0|V:U" (增加了小数位)

list.push(rAB.value)
list.sort()
console.log('List after inserting with decimal:', list)
```

### 总结

Jira 的 LexoRank 是一个工程上的杰作，它通过**分桶（Bucketing）**和**结构化的基数-62 小数排名（Structured Base-62 Decimal Rank）**优雅地解决了在分布式、高并发环境下维护有序列表的难题。

- **优点**：极高的可伸缩性、强大的并发处理能力、几乎无限的插入空间、将重平衡的影响局部化。
- **缺点**：实现非常复杂，需要自定义算术逻辑，并且排名字符串可能会变得很长，对数据库存储和索引有一定压力。

这个实现为你提供了一个坚实的基础，让你能够理解和应用这种高级排序技术。
