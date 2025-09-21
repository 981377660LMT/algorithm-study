// The code below helps compile a regular-expression-like language
// into a deterministic finite automaton. For a good introduction to
// these concepts, see https://swtch.com/~rsc/regexp/regexp1.html
// https://swtch.com/~rsc/regexp/
//
// !其实可以看一下golang标准库的正则实现.

/**
 * @file content_engine.ts
 * @description A generic, reusable engine for matching sequences of objects against a regular-expression-like pattern.
 *
 * This code is a refactored and generalized version of prosemirror-model's content.ts.
 * It replaces the ProseMirror-specific `NodeType` with a generic `MatchableType` interface,
 * allowing it to validate and manipulate any sequence of objects that conform to this interface.
 *
 * @version 1.0.0
 * @license MIT
 */

// ===================================================================================
// 1. PUBLIC INTERFACE & TYPES
// ===================================================================================

/**
 * Describes the objects that can be part of a sequence.
 * Your own object types must implement this interface to be used with the engine.
 * @template T The concrete type implementing this interface.
 */
export interface MatchableType<T extends MatchableType<T>> {
  /** A unique name for this type. */
  readonly name: string
  /**
   * Whether an object of this type can be created automatically.
   * Used by `fillBefore` and `findWrapping` to insert missing content.
   */
  readonly isGeneratable: boolean
  /**
   * If this type is a container, this property holds the `ContentMatch`
   * instance that describes its own valid content. For leaf types, this
   * should point to `ContentMatch.empty`.
   */
  readonly contentMatcher: ContentMatch<T>
  /**
   * Checks if this type belongs to a named group.
   * @param groupName The name of the group.
   */
  isInGroup(groupName: string): boolean
  /**
   * Creates a default instance of this type.
   * Called by `fillBefore` to generate filling content.
   */
  createDefault(): { type: T; [key: string]: any }
  /**
   * Whether this type is a "leaf" and cannot contain other content.
   * `findWrapping` will not try to wrap content with a leaf type.
   */
  readonly isLeaf: boolean
}

/**
 * A simple interface to represent a sequence of matchable objects.
 * ProseMirror's `Fragment` is compatible with this.
 * @template T The concrete `MatchableType`.
 */
export interface MatchableFragment<T extends MatchableType<T>> {
  readonly childCount: number
  child(index: number): { readonly type: T }
  // Helper to create a fragment from an array of objects.
  // You might need to implement this part based on your project's structure.
  // static from<T extends MatchableType<T>>(nodes: { type: T }[]): MatchableFragment<T>;
}

type MatchEdge<T extends MatchableType<T>> = { type: T; next: ContentMatch<T> }

/**
 * Instances of this class represent a match state of a content expression.
 * They can be used to find out whether further content matches here,
 * and whether a given position is a valid end of the expression.
 * @template T The concrete `MatchableType`.
 */
export class ContentMatch<T extends MatchableType<T>> {
  readonly next: MatchEdge<T>[] = []
  readonly wrapCache: (T | readonly T[] | null)[] = []

  constructor(
    /** True when this match state represents a valid end of the expression. */
    readonly validEnd: boolean
  ) {}

  /**
   * Parses a content expression string into a `ContentMatch` entry state.
   * @param string The content expression.
   * @param types A map of type names to `MatchableType` instances.
   */
  static parse<T extends MatchableType<T>>(
    string: string,
    types: { readonly [name: string]: T }
  ): ContentMatch<T> {
    const stream = new TokenStream(string, types)
    if (stream.next == null) return ContentMatch.empty as ContentMatch<T>
    const expr = parseExpr(stream)
    if (stream.next) stream.err('Unexpected trailing text')
    const match = dfa(nfa(expr))
    checkForDeadEnds(match, stream)
    return match
  }

  /** Match a type, returning a match after that type if successful. */
  matchType(type: T): ContentMatch<T> | null {
    for (let i = 0; i < this.next.length; i++)
      if (this.next[i].type == type) return this.next[i].next
    return null
  }

  /** Try to match a fragment. Returns the resulting match when successful. */
  matchFragment(
    frag: MatchableFragment<T>,
    start = 0,
    end = frag.childCount
  ): ContentMatch<T> | null {
    let cur: ContentMatch<T> | null = this
    for (let i = start; cur && i < end; i++) cur = cur.matchType(frag.child(i).type)
    return cur
  }

  /** Get the first matching generatable type at this match position. */
  get defaultType(): T | null {
    for (let i = 0; i < this.next.length; i++) {
      const { type } = this.next[i]
      if (type.isGeneratable) return type
    }
    return null
  }

  /**
   * Try to match the given fragment, and if that fails, see if it can
   * be made to match by inserting nodes in front of it.
   */
  fillBefore(
    after: MatchableFragment<T>,
    toEnd = false,
    startIndex = 0
  ): {
    fragment: { type: T }[]
    match: ContentMatch<T>
  } | null {
    const seen: ContentMatch<T>[] = [this]
    function search(
      match: ContentMatch<T>,
      types: readonly T[]
    ): {
      fragment: { type: T }[]
      match: ContentMatch<T>
    } | null {
      const finished = match.matchFragment(after, startIndex)
      if (finished && (!toEnd || finished.validEnd))
        return {
          fragment: types.map(tp => tp.createDefault()),
          match: finished
        }

      for (let i = 0; i < match.next.length; i++) {
        const { type, next } = match.next[i]
        if (type.isGeneratable && seen.indexOf(next) == -1) {
          seen.push(next)
          const found = search(next, types.concat(type))
          if (found) return found
        }
      }
      return null
    }

    return search(this, [])
  }

  /** Find a set of wrapping node types that would allow a node of the given type to appear at this position. */
  findWrapping(target: T): readonly T[] | null {
    for (let i = 0; i < this.wrapCache.length; i += 2)
      if (this.wrapCache[i] == target) return this.wrapCache[i + 1] as readonly T[] | null
    const computed = this.computeWrapping(target)
    this.wrapCache.push(target, computed)
    return computed
  }

  private computeWrapping(target: T): readonly T[] | null {
    type Active = { match: ContentMatch<T>; type: T | null; via: Active | null }
    const seen = Object.create(null)
    const active: Active[] = [{ match: this, type: null, via: null }]
    while (active.length) {
      const current = active.shift()!
      const match = current.match
      if (match.matchType(target)) {
        const result: T[] = []
        for (let obj: Active | null = current; obj && obj.type; obj = obj.via) result.push(obj.type)
        return result.reverse()
      }
      for (let i = 0; i < match.next.length; i++) {
        const { type, next } = match.next[i]
        if (
          !type.isLeaf &&
          type.isGeneratable &&
          !(type.name in seen) &&
          (!current.type || next.validEnd)
        ) {
          active.push({ match: type.contentMatcher, type, via: current })
          seen[type.name] = true
        }
      }
    }
    return null
  }

  /** The number of outgoing edges this node has. */
  get edgeCount() {
    return this.next.length
  }

  /** Get the nth outgoing edge from this node. */
  edge(n: number): MatchEdge<T> {
    if (n >= this.next.length) throw new RangeError(`There's no ${n}th edge in this content match`)
    return this.next[n]
  }

  /** A representation of the automaton for debugging. */
  toString() {
    const seen: ContentMatch<T>[] = []
    function scan(m: ContentMatch<T>) {
      seen.push(m)
      for (let i = 0; i < m.next.length; i++)
        if (seen.indexOf(m.next[i].next) == -1) scan(m.next[i].next)
    }
    scan(this)
    return seen
      .map((m, i) => {
        let out = i + (m.validEnd ? '*' : ' ') + ' '
        for (let i = 0; i < m.next.length; i++)
          out += (i ? ', ' : '') + m.next[i].type.name + '->' + seen.indexOf(m.next[i].next)
        return out
      })
      .join('\n')
  }

  /** An empty match that is always valid. */
  static empty = new ContentMatch<any>(true)
}

// ===================================================================================
// 2. INTERNAL IMPLEMENTATION (PARSER & COMPILER)
// ===================================================================================

class TokenStream<T extends MatchableType<T>> {
  pos = 0
  tokens: string[]

  constructor(readonly string: string, readonly types: { readonly [name: string]: T }) {
    this.tokens = string.split(/\s*(?=\b|\W|$)/)
    if (this.tokens[this.tokens.length - 1] == '') this.tokens.pop()
    if (this.tokens[0] == '') this.tokens.shift()
  }

  get next() {
    return this.tokens[this.pos]
  }

  eat(tok: string): boolean {
    return this.next == tok && this.pos++ > -1
  }

  err(str: string): never {
    throw new SyntaxError(str + " (in content expression '" + this.string + "')")
  }
}

type Expr<T extends MatchableType<T>> =
  | { type: 'choice'; exprs: Expr<T>[] }
  | { type: 'seq'; exprs: Expr<T>[] }
  | { type: 'plus'; expr: Expr<T> }
  | { type: 'star'; expr: Expr<T> }
  | { type: 'opt'; expr: Expr<T> }
  | { type: 'range'; min: number; max: number; expr: Expr<T> }
  | { type: 'name'; value: T }

function parseExpr<T extends MatchableType<T>>(stream: TokenStream<T>): Expr<T> {
  const exprs: Expr<T>[] = []
  do {
    exprs.push(parseExprSeq(stream))
  } while (stream.eat('|'))
  return exprs.length == 1 ? exprs[0] : { type: 'choice', exprs }
}

function parseExprSeq<T extends MatchableType<T>>(stream: TokenStream<T>): Expr<T> {
  const exprs: Expr<T>[] = []
  do {
    exprs.push(parseExprSubscript(stream))
  } while (stream.next && stream.next != ')' && stream.next != '|')
  return exprs.length == 1 ? exprs[0] : { type: 'seq', exprs }
}

function parseExprSubscript<T extends MatchableType<T>>(stream: TokenStream<T>): Expr<T> {
  let expr = parseExprAtom(stream)
  for (;;) {
    if (stream.eat('+')) expr = { type: 'plus', expr }
    else if (stream.eat('*')) expr = { type: 'star', expr }
    else if (stream.eat('?')) expr = { type: 'opt', expr }
    else if (stream.eat('{')) expr = parseExprRange(stream, expr)
    else break
  }
  return expr
}

function parseNum<T extends MatchableType<T>>(stream: TokenStream<T>): number {
  if (/\D/.test(stream.next)) stream.err("Expected number, got '" + stream.next + "'")
  const result = Number(stream.next)
  stream.pos++
  return result
}

function parseExprRange<T extends MatchableType<T>>(
  stream: TokenStream<T>,
  expr: Expr<T>
): Expr<T> {
  let min = parseNum(stream),
    max = min
  if (stream.eat(',')) {
    if (stream.next != '}') max = parseNum(stream)
    else max = -1
  }
  if (!stream.eat('}')) stream.err('Unclosed braced range')
  return { type: 'range', min, max, expr }
}

function resolveName<T extends MatchableType<T>>(
  stream: TokenStream<T>,
  name: string
): readonly T[] {
  const types = stream.types
  const type = types[name]
  if (type) return [type]
  const result: T[] = []
  for (const typeName in types) {
    const type = types[typeName]
    if (type.isInGroup(name)) result.push(type)
  }
  if (result.length == 0) stream.err("No type or group '" + name + "' found")
  return result
}

function parseExprAtom<T extends MatchableType<T>>(stream: TokenStream<T>): Expr<T> {
  if (stream.eat('(')) {
    const expr = parseExpr(stream)
    if (!stream.eat(')')) stream.err('Missing closing paren')
    return expr
  } else if (!/\W/.test(stream.next)) {
    const exprs = resolveName(stream, stream.next).map(type => {
      return { type: 'name', value: type } as Expr<T>
    })
    stream.pos++
    return exprs.length == 1 ? exprs[0] : { type: 'choice', exprs }
  } else {
    stream.err("Unexpected token '" + stream.next + "'")
  }
}

type NFAEdge<T extends MatchableType<T>> = { term: T | undefined; to: number | undefined }

function nfa<T extends MatchableType<T>>(expr: Expr<T>): NFAEdge<T>[][] {
  const nfa: NFAEdge<T>[][] = [[]]
  connect(compile(expr, 0), node())
  return nfa

  function node(): number {
    return nfa.push([]) - 1
  }
  function edge(from: number, to?: number, term?: T): NFAEdge<T> {
    const edge = { term, to }
    nfa[from].push(edge)
    return edge
  }
  function connect(edges: NFAEdge<T>[], to: number) {
    edges.forEach(edge => (edge.to = to))
  }

  function compile(expr: Expr<T>, from: number): NFAEdge<T>[] {
    if (expr.type == 'choice') {
      return expr.exprs.reduce((out, expr) => out.concat(compile(expr, from)), [] as NFAEdge<T>[])
    } else if (expr.type == 'seq') {
      for (let i = 0; ; i++) {
        const next = compile(expr.exprs[i], from)
        if (i == expr.exprs.length - 1) return next
        connect(next, (from = node()))
      }
    } else if (expr.type == 'star') {
      const loop = node()
      edge(from, loop)
      connect(compile(expr.expr, loop), loop)
      return [edge(loop)]
    } else if (expr.type == 'plus') {
      const loop = node()
      connect(compile(expr.expr, from), loop)
      connect(compile(expr.expr, loop), loop)
      return [edge(loop)]
    } else if (expr.type == 'opt') {
      return [edge(from)].concat(compile(expr.expr, from))
    } else if (expr.type == 'range') {
      let cur = from
      for (let i = 0; i < expr.min; i++) {
        const next = node()
        connect(compile(expr.expr, cur), next)
        cur = next
      }
      if (expr.max == -1) {
        connect(compile(expr.expr, cur), cur)
      } else {
        for (let i = expr.min; i < expr.max; i++) {
          const next = node()
          edge(cur, next)
          connect(compile(expr.expr, cur), next)
          cur = next
        }
      }
      return [edge(cur)]
    } else if (expr.type == 'name') {
      return [edge(from, undefined, expr.value)]
    }
    throw new Error('Unknown expr type')
  }
}

function cmp(a: number, b: number): number {
  return b - a
}

function nullFrom<T extends MatchableType<T>>(
  nfa: NFAEdge<T>[][],
  node: number
): readonly number[] {
  const result: number[] = []
  scan(node)
  return result.sort(cmp)

  function scan(node: number): void {
    const edges = nfa[node]
    if (edges.length == 1 && !edges[0].term) return scan(edges[0].to!)
    result.push(node)
    for (let i = 0; i < edges.length; i++) {
      const { term, to } = edges[i]
      if (!term && result.indexOf(to!) == -1) scan(to!)
    }
  }
}

function dfa<T extends MatchableType<T>>(nfa: NFAEdge<T>[][]): ContentMatch<T> {
  const labeled: { [key: string]: ContentMatch<T> } = Object.create(null)
  return explore(nullFrom(nfa, 0))

  function explore(states: readonly number[]): ContentMatch<T> {
    const out: [T, number[]][] = []
    states.forEach(node => {
      nfa[node].forEach(({ term, to }) => {
        if (!term) return
        let set: number[] | undefined
        for (let i = 0; i < out.length; i++) if (out[i][0] == term) set = out[i][1]
        nullFrom(nfa, to!).forEach(node => {
          if (!set) out.push([term, (set = [])])
          if (set.indexOf(node) == -1) set.push(node)
        })
      })
    })
    const state = (labeled[states.join(',')] = new ContentMatch<T>(
      states.indexOf(nfa.length - 1) > -1
    ))
    for (let i = 0; i < out.length; i++) {
      const states = out[i][1].sort(cmp)
      state.next.push({ type: out[i][0], next: labeled[states.join(',')] || explore(states) })
    }
    return state
  }
}

function checkForDeadEnds<T extends MatchableType<T>>(
  match: ContentMatch<T>,
  stream: TokenStream<T>
) {
  for (let i = 0, work = [match]; i < work.length; i++) {
    const state = work[i]
    let dead = !state.validEnd
    const nodes: string[] = []
    for (let j = 0; j < state.next.length; j++) {
      const { type, next } = state.next[j]
      nodes.push(type.name)
      if (dead && type.isGeneratable) dead = false
      if (work.indexOf(next) == -1) work.push(next)
    }
    if (dead)
      stream.err('Only non-generatable types (' + nodes.join(', ') + ') in a required position')
  }
}

export {}
