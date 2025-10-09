// !codemirror 自动补全的匹配策略
// 自动补全（Autocomplete）中的建议项匹配：当用户在补全列表中输入时，如何筛选和高亮建议项。
// 模糊匹配与计分 (Fuzzy Matching & Scoring)

// Scores are counted from 0 (great match) down to negative numbers,
// assigning specific penalty values for specific shortcomings.
const enum Penalty {
  Gap = -1100, // Added for each gap in the match (not counted for by-word matches)
  NotStart = -700, // The match doesn't start at the start of the word
  CaseFold = -200, // At least one character needed to be case-folded to match
  ByWord = -100, // The match is by-word, meaning each char in the pattern matches the start of a word in the string
  NotFull = -100 // Used to push down matches that don't match the pattern fully relative to those that do
}

const enum Tp {
  NonWord,
  Upper,
  Lower
}

// A pattern matcher for fuzzy completion matching. Create an instance
// once for a pattern, and then use that to match any number of
// completions.
class FuzzyMatcher {
  chars: number[] = []
  folded: number[] = []
  astral: boolean

  // Buffers reused by calls to `match` to track matched character
  // positions.
  any: number[] = []
  precise: number[] = []
  byWord: number[] = []

  score = 0
  matched: readonly number[] = []

  constructor(readonly pattern: string) {
    for (let p = 0; p < pattern.length; ) {
      let char = codePointAt(pattern, p),
        size = codePointSize(char)
      this.chars.push(char)
      let part = pattern.slice(p, p + size),
        upper = part.toUpperCase()
      this.folded.push(codePointAt(upper == part ? part.toLowerCase() : upper, 0))
      p += size
    }
    this.astral = pattern.length != this.chars.length
  }

  ret(score: number, matched: readonly number[]) {
    this.score = score
    this.matched = matched
    return this
  }

  // Matches a given word (completion) against the pattern (input).
  // Will return a boolean indicating whether there was a match and,
  // on success, set `this.score` to the score, `this.matched` to an
  // array of `from, to` pairs indicating the matched parts of `word`.
  //
  // The score is a number that is more negative the worse the match
  // is. See `Penalty` above.
  match(word: string): { score: number; matched: readonly number[] } | null {
    if (this.pattern.length == 0) return this.ret(Penalty.NotFull, [])
    if (word.length < this.pattern.length) return null
    let { chars, folded, any, precise, byWord } = this
    // For single-character queries, only match when they occur right
    // at the start
    if (chars.length == 1) {
      let first = codePointAt(word, 0),
        firstSize = codePointSize(first)
      let score = firstSize == word.length ? 0 : Penalty.NotFull
      if (first == chars[0]) {
      } else if (first == folded[0]) score += Penalty.CaseFold
      else return null
      return this.ret(score, [0, firstSize])
    }
    let direct = word.indexOf(this.pattern)
    if (direct == 0)
      return this.ret(word.length == this.pattern.length ? 0 : Penalty.NotFull, [
        0,
        this.pattern.length
      ])

    let len = chars.length,
      anyTo = 0
    if (direct < 0) {
      for (let i = 0, e = Math.min(word.length, 200); i < e && anyTo < len; ) {
        let next = codePointAt(word, i)
        if (next == chars[anyTo] || next == folded[anyTo]) any[anyTo++] = i
        i += codePointSize(next)
      }
      // No match, exit immediately
      if (anyTo < len) return null
    }

    // This tracks the extent of the precise (non-folded, not
    // necessarily adjacent) match
    let preciseTo = 0
    // Tracks whether there is a match that hits only characters that
    // appear to be starting words. `byWordFolded` is set to true when
    // a case folded character is encountered in such a match
    let byWordTo = 0,
      byWordFolded = false
    // If we've found a partial adjacent match, these track its state
    let adjacentTo = 0,
      adjacentStart = -1,
      adjacentEnd = -1
    let hasLower = /[a-z]/.test(word),
      wordAdjacent = true
    // Go over the option's text, scanning for the various kinds of matches
    for (
      let i = 0, e = Math.min(word.length, 200), prevType = Tp.NonWord;
      i < e && byWordTo < len;

    ) {
      let next = codePointAt(word, i)
      if (direct < 0) {
        if (preciseTo < len && next == chars[preciseTo]) precise[preciseTo++] = i
        if (adjacentTo < len) {
          if (next == chars[adjacentTo] || next == folded[adjacentTo]) {
            if (adjacentTo == 0) adjacentStart = i
            adjacentEnd = i + 1
            adjacentTo++
          } else {
            adjacentTo = 0
          }
        }
      }
      let ch,
        type =
          next < 0xff
            ? (next >= 48 && next <= 57) || (next >= 97 && next <= 122)
              ? Tp.Lower
              : next >= 65 && next <= 90
              ? Tp.Upper
              : Tp.NonWord
            : (ch = fromCodePoint(next)) != ch.toLowerCase()
            ? Tp.Upper
            : ch != ch.toUpperCase()
            ? Tp.Lower
            : Tp.NonWord
      if (!i || (type == Tp.Upper && hasLower) || (prevType == Tp.NonWord && type != Tp.NonWord)) {
        if (chars[byWordTo] == next || (folded[byWordTo] == next && (byWordFolded = true)))
          byWord[byWordTo++] = i
        else if (byWord.length) wordAdjacent = false
      }
      prevType = type
      i += codePointSize(next)
    }

    if (byWordTo == len && byWord[0] == 0 && wordAdjacent)
      return this.result(Penalty.ByWord + (byWordFolded ? Penalty.CaseFold : 0), byWord, word)
    if (adjacentTo == len && adjacentStart == 0)
      return this.ret(
        Penalty.CaseFold - word.length + (adjacentEnd == word.length ? 0 : Penalty.NotFull),
        [0, adjacentEnd]
      )
    if (direct > -1)
      return this.ret(Penalty.NotStart - word.length, [direct, direct + this.pattern.length])
    if (adjacentTo == len)
      return this.ret(Penalty.CaseFold + Penalty.NotStart - word.length, [
        adjacentStart,
        adjacentEnd
      ])
    if (byWordTo == len)
      return this.result(
        Penalty.ByWord +
          (byWordFolded ? Penalty.CaseFold : 0) +
          Penalty.NotStart +
          (wordAdjacent ? 0 : Penalty.Gap),
        byWord,
        word
      )
    return chars.length == 2
      ? null
      : this.result((any[0] ? Penalty.NotStart : 0) + Penalty.CaseFold + Penalty.Gap, any, word)
  }

  result(score: number, positions: number[], word: string) {
    let result: number[] = [],
      i = 0
    for (let pos of positions) {
      let to = pos + (this.astral ? codePointSize(codePointAt(word, pos)) : 1)
      if (i && result[i - 1] == pos) result[i - 1] = to
      else {
        result[i++] = pos
        result[i++] = to
      }
    }
    return this.ret(score - word.length, result)
  }
}

class StrictMatcher {
  matched: readonly number[] = []
  score: number = 0
  folded: string

  constructor(readonly pattern: string) {
    this.folded = pattern.toLowerCase()
  }

  match(word: string): { score: number; matched: readonly number[] } | null {
    if (word.length < this.pattern.length) return null
    let start = word.slice(0, this.pattern.length)
    let match =
      start == this.pattern ? 0 : start.toLowerCase() == this.folded ? Penalty.CaseFold : null
    if (match == null) return null
    this.matched = [0, start.length]
    this.score = match + (word.length == this.pattern.length ? 0 : Penalty.NotFull)
    return this
  }
}

function surrogateLow(ch: number) {
  return ch >= 0xdc00 && ch < 0xe000
}
function surrogateHigh(ch: number) {
  return ch >= 0xd800 && ch < 0xdc00
}

/// Find the code point at the given position in a string (like the
/// [`codePointAt`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/codePointAt)
/// string method).
function codePointAt(str: string, pos: number) {
  let code0 = str.charCodeAt(pos)
  if (!surrogateHigh(code0) || pos + 1 == str.length) return code0
  let code1 = str.charCodeAt(pos + 1)
  if (!surrogateLow(code1)) return code0
  return ((code0 - 0xd800) << 10) + (code1 - 0xdc00) + 0x10000
}

/// Given a Unicode codepoint, return the JavaScript string that
/// respresents it (like
/// [`String.fromCodePoint`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/fromCodePoint)).
function fromCodePoint(code: number) {
  if (code <= 0xffff) return String.fromCharCode(code)
  code -= 0x10000
  return String.fromCharCode((code >> 10) + 0xd800, (code & 1023) + 0xdc00)
}

/// The amount of positions a character takes up in a JavaScript string.
function codePointSize(code: number): 1 | 2 {
  return code < 0x10000 ? 1 : 2
}

export {}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  /**
   * 辅助函数，用于打印匹配结果
   * @param matcherName 匹配器名称
   * @param pattern 搜索模式
   * @param word 被匹配的单词
   * @param result 匹配结果
   */
  function printResult(
    matcherName: string,
    pattern: string,
    word: string,
    result: { score: number; matched: readonly number[] } | null
  ) {
    console.log(`--- [${matcherName}] 模式: "${pattern}", 单词: "${word}" ---`)
    if (result) {
      console.log(`  ✅ 匹配成功!`)
      console.log(`  得分: ${result.score}`)
      console.log(`  匹配位置: ${JSON.stringify(result.matched)}`)
      // 根据匹配位置高亮显示
      let highlighted = ''
      let last = 0
      for (let i = 0; i < result.matched.length; i += 2) {
        const from = result.matched[i]
        const to = result.matched[i + 1]
        highlighted += word.slice(last, from)
        highlighted += `\x1b[31m${word.slice(from, to)}\x1b[0m` // 使用 ANSI escape code 高亮
        last = to
      }
      highlighted += word.slice(last)
      console.log(`  高亮: ${highlighted}`)
    } else {
      console.log(`  ❌ 未匹配`)
    }
    console.log('')
  }

  // --- 1. FuzzyMatcher 使用案例 ---
  // 它的匹配策略更灵活，允许字符不相邻、大小写不敏感等。
  console.log('============== FuzzyMatcher 演示 ==============')

  // 创建一个针对模式 "cm" 的模糊匹配器
  const fuzzy = new FuzzyMatcher('cm')

  // 场景 1: 驼峰式匹配 (CodeMirror)
  // 'c' 匹配 'C', 'm' 匹配 'M'。这是非常好的匹配（按词匹配）。
  let result1 = fuzzy.match('CodeMirror')
  printResult('Fuzzy', 'cm', 'CodeMirror', result1)

  // 场景 2: 包含子串，但大小写不同 (codemirror)
  // 'c' 匹配 'c', 'm' 匹配 'm'。
  let result2 = fuzzy.match('codemirror')
  printResult('Fuzzy', 'cm', 'codemirror', result2)

  // 场景 3: 字符不相邻，有间隔 (custom)
  // 'c' 匹配 'c', 'm' 匹配 'm'，但中间有 'usto'。会有间隙(Gap)惩罚，分数更低。
  let result3 = fuzzy.match('custom')
  printResult('Fuzzy', 'cm', 'custom', result3)

  // 场景 4: 不在开头 (become)
  // 'c' 和 'm' 都能匹配，但 'c' 不在单词开头。会有非开头(NotStart)惩罚。
  let result4 = fuzzy.match('become')
  printResult('Fuzzy', 'cm', 'become', result4)

  // 场景 5: 完全不匹配
  let result5 = fuzzy.match('typescript')
  printResult('Fuzzy', 'cm', 'typescript', result5)

  // --- 2. StrictMatcher 使用案例 ---
  // 它的策略非常严格，只匹配单词的开头部分（前缀匹配），但会忽略大小写。
  console.log('\n============== StrictMatcher 演示 ==============')

  // 创建一个针对模式 "auto" 的严格匹配器
  const strict = new StrictMatcher('auto')

  // 场景 1: 精确前缀匹配
  let result6 = strict.match('autocomplete')
  printResult('Strict', 'auto', 'autocomplete', result6)

  // 场景 2: 大小写不同的前缀匹配
  let result7 = strict.match('AutoComplete')
  printResult('Strict', 'auto', 'AutoComplete', result7)

  // 场景 3: 不是前缀，不匹配
  let result8 = strict.match('my-autocomplete')
  printResult('Strict', 'auto', 'my-autocomplete', result8)
}
