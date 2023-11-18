/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-inner-declarations */

/**
 * @deprecated
 */
class ACAutoMatonArrayLegacy {
  /**
   * 模式串列表.
   */
  private readonly _pattern: string[] = []

  /**
   * Trie树.
   * `_children[i]`表示节点(状态)`i`的所有子节点,`0`表示虚拟根节点.
   */
  private readonly _children: Int32Array[] = []

  /**
   * Trie树结点附带的信息.
   * `_matching[i]`表示节点(状态)`i`对应的字符串唯一标识.
   */
  private readonly _matching: number[][] = []

  /**
   * Trie树结点附带的信息.
   * `_wordCount[i]`表示节点(状态)`i`的匹配个数.
   */
  private _wordCount!: Uint32Array

  /**
   * `_fail[i]`表示节点(状态)`i`的失配指针.
   */
  private _fail!: Uint32Array

  /**
   * 构建时是否在 {@link _matching} 中处理出每个结点匹配到的模式串id.
   */
  private _heavy = false

  private readonly _size: number
  private readonly _margin: number

  /**
   * @param size 字符集大小.默认为26,即所有小写字母.
   * @param margin 字符集的起始字符.默认为97,即`a`的ascii码.
   */
  constructor(size = 26, margin = 97) {
    this._size = size
    this._margin = margin
    this._children.push(new Int32Array(size).fill(-1))
    this._matching.push([])
  }

  /**
   * 将模式串`pattern`插入到Trie树中.模式串一般是`被禁用的单词`.
   * @param pid 模式串的唯一标识id.
   * @param pattern 模式串.
   * @param didInsert 模式串插入后的回调函数,入参为结束字符所在的结点(状态).
   */
  insert(pid: number, pattern: string, didInsert?: (endState: number) => void): this {
    if (!pattern) return this

    let root = 0
    for (let i = 0; i < pattern.length; i++) {
      const char = pattern[i].charCodeAt(0)! - this._margin
      const nexts = this._children[root]
      if (~nexts[char]) {
        root = nexts[char]
      } else {
        const nextState = this._children.length
        nexts[char] = nextState
        root = nextState
        this._children.push(new Int32Array(this._size).fill(-1))
        this._matching.push([])
      }
    }

    this._matching[root].push(pid)
    this._pattern.push(pattern)
    if (didInsert) didInsert(root)
    return this
  }

  /**
   * 构建失配指针.
   * bfs为字典树的每个结点添加失配指针,结点要跳转到哪里.
   * AC自动机的失配指针指向的节点所代表的字符串 是 当前节点所代表的字符串的 最长后缀.
   * @param heavy 是否处理出每个结点匹配到的模式串id. 默认为False.
   * @param dp AC自动机构建过程中的回调函数,入参为`(next结点的fail指针, next结点)`.
   */
  build(heavy?: boolean, dp?: (move: number, next: number) => void): void {
    const n = this._children.length
    this._wordCount = new Uint32Array(n)
    for (let i = 0; i < n; i++) this._wordCount[i] = this._matching[i].length
    this._fail = new Uint32Array(n)
    this._heavy = !!heavy

    let queue: number[] = []
    for (let i = 0; i < this._size; i++) {
      const next = this._children[0][i]
      if (~next) queue.push(next)
    }
    while (queue.length) {
      const nextQueue: number[] = []
      const len = queue.length
      for (let i = 0; i < len; i++) {
        const cur = queue[i]
        const curFail = this._fail[cur]
        const curChildren = this._children[cur]
        for (let j = 0; j < this._size; j++) {
          const next = curChildren[j]
          if (~next) {
            const move = this.move(curFail, j)
            this._fail[next] = move // !更新子节点的fail指针
            this._wordCount[next] += this._wordCount[move] // !更新子节点的匹配个数
            heavy && this._matching[next].push(...this._matching[move]) // !更新move状态匹配的模式串下标
            dp && dp(move, next)
            nextQueue.push(next)
          }
        }
      }
      queue = nextQueue
    }
  }

  /**
   * 从当前状态`state`沿着字符`input`转移到的下一个状态.
   * 沿着失配链上跳,找到第一个可由input转移的节点.
   * @param state 当前状态.
   * @param input 输入字符或字符的ascii码.
   * @returns 下一个状态.
   */
  move(state: number, input: string | number): number {
    const char = typeof input === 'string' ? input.charCodeAt(0)! - this._margin : input
    while (true) {
      const nexts = this._children[state]
      if (~nexts[char]) return nexts[char]
      if (!state) return 0
      state = this._fail[state]
    }
  }

  /**
   * 从状态`state`开始匹配字符串`searchString`.
   * @param state ac自动机的状态.根节点状态为0.
   * @param searchString 待匹配的字符串.
   * @returns 每个模式串在`searchString`中出现的下标.
   */
  match(state: number, searchString: string): Map<number, number[]> {
    if (!this._heavy) throw new Error('需要调用build(true)构建AC自动机')
    const res = new Map<number, number[]>()
    let root = state
    for (let i = 0; i < searchString.length; i++) {
      root = this.move(root, searchString[i])
      this._matching[root].forEach(pid => {
        const start = i - this._pattern[pid].length + 1
        !res.has(pid) && res.set(pid, [])
        res.get(pid)!.push(start)
      })
    }
    return res
  }

  /**
   * 当前状态`state`匹配到的模式串个数.
   */
  count(state: number): number {
    return this._wordCount[state]
  }

  /**
   * 当前状态`state`是否为匹配状态.
   */
  accept(state: number): boolean {
    return !!this._wordCount[state]
  }
}

export { ACAutoMatonArrayLegacy }

if (require.main === module) {
  // 1032. 字符流
  // https://leetcode.cn/problems/stream-of-characters/
  class StreamChecker {
    private readonly _acm: ACAutoMatonArrayLegacy
    private _state = 0
    constructor(words: string[]) {
      this._acm = new ACAutoMatonArrayLegacy()
      words.forEach((word, i) => this._acm.insert(i, word))
      this._acm.build()
    }

    query(letter: string): boolean {
      this._state = this._acm.move(this._state, letter)
      return this._acm.accept(this._state)
    }
  }

  /**
   * Your StreamChecker object will be instantiated and called as such:
   * var obj = new StreamChecker(words)
   * var param_1 = obj.query(letter)
   */

  // 面试题 17.17. 多次搜索
  // https://leetcode.cn/problems/multi-search-lcci/
  function multiSearch(big: string, smalls: string[]): number[][] {
    const acm = new ACAutoMatonArrayLegacy()
    smalls.forEach((small, i) => acm.insert(i, small))
    acm.build(true)

    const res: number[][] = Array(smalls.length)
    for (let i = 0; i < smalls.length; i++) res[i] = []
    const matching = acm.match(0, big)
    matching.forEach((starts, pid) => {
      starts.forEach(start => res[pid].push(start))
    })
    return res
  }
}
