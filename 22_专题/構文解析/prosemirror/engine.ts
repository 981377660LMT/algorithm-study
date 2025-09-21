// 词法分析、语法分析、NFA 和 DFA 编译
//
// ===================================================================================
//  Standalone Regular Expression Engine for Object Sequences
//  Extracted and adapted from prosemirror-model/src/content.ts
// ===================================================================================

/**
 * The basic "alphabet" of our regular expression.
 * Any object you want to match must at least have a name.
 */
export interface Term {
  readonly name: string
  // Optional: for supporting group matching in expressions
  isInGroup?(groupName: string): boolean
}

// ===================================================================================
// 1. The DFA State (The Final Product)
// ===================================================================================

/**
 * Represents a state in the compiled Deterministic Finite Automaton (DFA).
 */
export class MatchState<T extends Term> {
  // An array of transitions to the next state.
  readonly next: { term: T; next: MatchState<T> }[] = []

  /**
   * @param validEnd True if this state represents a valid end of a sequence.
   */
  constructor(readonly validEnd: boolean) {}

  /**
   * Tries to match a term and returns the next state if successful.
   * This is the core function for running the DFA.
   */
  matchTerm(term: T): MatchState<T> | null {
    for (let i = 0; i < this.next.length; i++) {
      if (this.next[i].term === term) return this.next[i].next
    }
    return null
  }

  /**
   * Tries to match a sequence of terms.
   */
  matchSequence(sequence: readonly T[]): MatchState<T> | null {
    let current: MatchState<T> | null = this
    for (const term of sequence) {
      if (!current) break
      current = current.matchTerm(term)
    }
    return current
  }
}

// ===================================================================================
// 2. The Compiler Pipeline
// ===================================================================================

/**
 * Compiles a content expression string into a DFA.
 * @param expression The regular expression string (e.g., "a (b|c)*").
 * @param terms A map from term names to the actual Term objects.
 * @returns The starting state of the compiled DFA.
 */
export function compile<T extends Term>(
  expression: string,
  terms: { readonly [name: string]: T }
): MatchState<T> {
  const stream = new TokenStream(expression, terms)
  if (stream.next == null) return new MatchState<T>(true) // Empty expression matches empty sequence
  const ast = parseExpr(stream)
  if (stream.next) stream.err('Unexpected trailing text')
  const nfaGraph = nfa(ast)
  return dfa(nfaGraph)
}

// ===================================================================================
// 3. Internal Implementation (The "Engine Room")
// ===================================================================================

// -- 3.1 Lexer (Token Stream) --

class TokenStream<T extends Term> {
  pos = 0
  tokens: string[]

  constructor(readonly string: string, readonly terms: { readonly [name: string]: T }) {
    this.tokens = string.split(/\s*(?=\b|\W|$)/)
    if (this.tokens[this.tokens.length - 1] === '') this.tokens.pop()
    if (this.tokens[0] === '') this.tokens.shift()
  }

  get next() {
    return this.tokens[this.pos]
  }
  eat(tok: string): boolean {
    return this.next === tok && this.pos++ > -1
  }
  err(str: string): never {
    throw new SyntaxError(`${str} (in expression '${this.string}')`)
  }
}

// -- 3.2 Parser (String -> Abstract Syntax Tree) --

type Expr<T extends Term> =
  | { type: 'choice'; exprs: Expr<T>[] }
  | { type: 'seq'; exprs: Expr<T>[] }
  | { type: 'plus'; expr: Expr<T> }
  | { type: 'star'; expr: Expr<T> }
  | { type: 'opt'; expr: Expr<T> }
  | { type: 'range'; min: number; max: number; expr: Expr<T> }
  | { type: 'name'; value: T }

function parseExpr<T extends Term>(stream: TokenStream<T>): Expr<T> {
  const exprs: Expr<T>[] = []
  do {
    exprs.push(parseExprSeq(stream))
  } while (stream.eat('|'))
  return exprs.length == 1 ? exprs[0] : { type: 'choice', exprs }
}

function parseExprSeq<T extends Term>(stream: TokenStream<T>): Expr<T> {
  const exprs: Expr<T>[] = []
  do {
    exprs.push(parseExprSubscript(stream))
  } while (stream.next && stream.next !== ')' && stream.next !== '|')
  return exprs.length == 1 ? exprs[0] : { type: 'seq', exprs }
}

function parseExprSubscript<T extends Term>(stream: TokenStream<T>): Expr<T> {
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

function parseExprAtom<T extends Term>(stream: TokenStream<T>): Expr<T> {
  if (stream.eat('(')) {
    const expr = parseExpr(stream)
    if (!stream.eat(')')) stream.err('Missing closing paren')
    return expr
  } else if (!/\W/.test(stream.next)) {
    const exprs = resolveName(stream, stream.next).map(
      type => ({ type: 'name', value: type } as Expr<T>)
    )
    stream.pos++
    return exprs.length == 1 ? exprs[0] : { type: 'choice', exprs }
  } else {
    stream.err(`Unexpected token '${stream.next}'`)
  }
}

function parseNum<T extends Term>(stream: TokenStream<T>): number {
  if (/\D/.test(stream.next)) stream.err(`Expected number, got '${stream.next}'`)
  const result = Number(stream.next)
  stream.pos++
  return result
}

function parseExprRange<T extends Term>(stream: TokenStream<T>, expr: Expr<T>): Expr<T> {
  let min = parseNum(stream),
    max = min
  if (stream.eat(',')) {
    if (stream.next !== '}') max = parseNum(stream)
    else max = -1
  }
  if (!stream.eat('}')) stream.err('Unclosed braced range')
  return { type: 'range', min, max, expr }
}

function resolveName<T extends Term>(stream: TokenStream<T>, name: string): readonly T[] {
  const defined = stream.terms
  const found = defined[name]
  if (found) return [found]
  const result: T[] = []
  for (const termName in defined) {
    const term = defined[termName]
    if (term.isInGroup?.(name)) result.push(term)
  }
  if (result.length == 0) stream.err(`No term or group '${name}' found`)
  return result
}

// -- 3.3 NFA Compiler (AST -> Non-deterministic Finite Automaton) --

type NFAEdge<T extends Term> = { term: T | undefined; to: number | undefined }

function nfa<T extends Term>(expr: Expr<T>): NFAEdge<T>[][] {
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
    if (expr.type === 'choice') {
      return expr.exprs.reduce((out, expr) => out.concat(compile(expr, from)), [] as NFAEdge<T>[])
    } else if (expr.type === 'seq') {
      for (let i = 0; ; i++) {
        const next = compile(expr.exprs[i], from)
        if (i === expr.exprs.length - 1) return next
        connect(next, (from = node()))
      }
    } else if (expr.type === 'star') {
      const loop = node()
      edge(from, loop)
      connect(compile(expr.expr, loop), loop)
      return [edge(loop)]
    } else if (expr.type === 'plus') {
      const loop = node()
      connect(compile(expr.expr, from), loop)
      connect(compile(expr.expr, loop), loop)
      return [edge(loop)]
    } else if (expr.type === 'opt') {
      return [edge(from)].concat(compile(expr.expr, from))
    } else if (expr.type === 'range') {
      let cur = from
      for (let i = 0; i < expr.min; i++) {
        const next = node()
        connect(compile(expr.expr, cur), next)
        cur = next
      }
      if (expr.max === -1) {
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
    } else if (expr.type === 'name') {
      return [edge(from, undefined, expr.value)]
    }
    throw new Error('Unknown expression type')
  }
}

// -- 3.4 DFA Compiler (NFA -> Deterministic Finite Automaton) --

function dfa<T extends Term>(nfa: NFAEdge<T>[][]): MatchState<T> {
  const labeled: { [key: string]: MatchState<T> } = Object.create(null)
  return explore(nullFrom(nfa, 0))

  function explore(states: readonly number[]): MatchState<T> {
    const key = states.join(',')
    if (key in labeled) return labeled[key]

    const out: [T, number[]][] = []
    for (const node of states) {
      for (const { term, to } of nfa[node]) {
        if (!term) continue
        let set: number[] | undefined
        for (let i = 0; i < out.length; i++) if (out[i][0] === term) set = out[i][1]
        for (const node of nullFrom(nfa, to!)) {
          if (!set) out.push([term, (set = [])])
          if (set.indexOf(node) === -1) set.push(node)
        }
      }
    }

    const state = (labeled[key] = new MatchState<T>(states.indexOf(nfa.length - 1) > -1))
    for (const [term, newStates] of out) {
      state.next.push({ term, next: explore(newStates.sort((a, b) => b - a)) })
    }
    return state
  }
}

function nullFrom<T extends Term>(nfa: NFAEdge<T>[][], node: number): readonly number[] {
  const result: number[] = []
  const seen: boolean[] = []

  function scan(node: number): void {
    if (seen[node]) return
    seen[node] = true
    const edges = nfa[node]
    if (edges.length === 1 && !edges[0].term) return scan(edges[0].to!)
    result.push(node)
    for (const { term, to } of edges) {
      if (!term) scan(to!)
    }
  }

  scan(node)
  return result.sort((a, b) => b - a)
}

const A = { name: 'a' }
const B = { name: 'b' }
const C = { name: 'c' }
const ALL_TERMS = { a: A, b: B, c: C }

const expression = 'a (b | c)+'
const startState = compile(expression, ALL_TERMS)

function test(sequence: Term[]): boolean {
  const finalState = startState.matchSequence(sequence)
  return !!finalState && finalState.validEnd
}

console.log(`Testing against: ${expression}`)
console.log('[a, b] ->', test([A, B])) // true
console.log('[a, c] ->', test([A, C])) // true
console.log('[a, b, c] ->', test([A, B, C])) // true
console.log('[a] ->', test([A])) // false (b|c)+ requires at least one
console.log('[b] ->', test([B])) // false (must start with a)
console.log('[a, b, a] ->', test([A, B, A])) // false (a is not allowed after b)
