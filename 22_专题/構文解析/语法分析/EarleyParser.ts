//
// Earley Parser
//
// Description:
//   We are given CFG, i.e.,
//     A -> B
//     A -> aAa|bAb
//     B -> aa|bb
//     B -> a|b
//  It determines that a given string is matched by the CFG.
//  !上下文無關文法（CFG）
//
// Algorithm:
//   Earley algorithm. It generates all states with memoisation.
//   Here, state is given by (rule, pos-in-rule, pos-in-text).
//
// Complexity:
//   O(|G|^2 n^3) in the worst case.
//   If a grammar is simple, it usually reduced to O(|G|^2 n^2).
//
// Remark:
//   Because of simplicity, This implementation does not allow the
//   !epsilon rule. Please expand epsilon rule by hand.
//   (TODO!)
//
//
// struct earley_parser {
//   vector<int> terminal;
//   vector<vector<vector<int>>> grammar;
//   int add_symbol(char c = 0) {
//     terminal.push_back(c);·
//     grammar.push_back({});
//     return grammar.size()-1;
//   }
//   void add_grammar(int A, vector<int> As) {
//     As.push_back(0);
//     grammar[A].push_back(As);
//   }
//   earley_parser() { add_symbol(); add_symbol(); }
//   bool parse(const char s[], int init) {
//     int n = strlen(s);
//     struct state { int a, k, p, i; };
//     vector<vector<vector<state>>> chart(n+1, vector<vector<state>>(grammar.size()));
//     auto enqueue = [&](vector<state> &curr, const state &S) {
//       for (auto &T: curr)
//         if (T.a == S.a && T.k == S.k && T.p == S.p && T.i == S.i) return;
//       curr.push_back(S);
//     };
//     auto symbol = [&](const state &S) { return grammar[S.a][S.k][S.p]; };
//     grammar[1] = { {init, 0} };
//     vector<state> curr = {{1, 0, 0, 0}}, next;
//     for (int k = 0; k <= n; ++k) {
//       for (int i = 0; i < curr.size(); ++i) {
//         state S = curr[i];
//         int B = symbol(S);
//         if (B) {
//           if (!terminal[B]) {
//             for (int j = 0; j < grammar[B].size(); ++j)
//               enqueue(curr, {B, j, 0, k});
//           } else if (terminal[B] == s[k]) {
//             enqueue(next, {S.a, S.k, S.p+1, S.i});
//           }
//         } else {
//           for (auto &T: chart[S.i][S.a])
//             enqueue(curr, {T.a, T.k, T.p+1, T.i});
//         }
//       }
//       for (auto &T: curr)
//         chart[k][symbol(T)].push_back(T);
//       curr.swap(next);
//       next.clear();
//     }
//     for (auto &T: chart[n][0])
//       if (T.a == 1) return true;
//     return false;
//   }
// };

type State = readonly [a: number, k: number, p: number, i: number]

/**
 * Earley解析器.
 *
 * Rust中宏调用的解析是基于Earley算法的.
 *
 * @link https://www.python123.io/index/topics/algorithm_100_days/100-days-of-algorithms-94
 */
class EarleyParser {
  private readonly _terminal: number[] = []
  private readonly _grammar: number[][][] = []

  constructor() {
    this.addSymbol()
    this.addSymbol()
  }

  addSymbol(c = 0): number {
    this._terminal.push(c)
    this._grammar.push([])
    return this._grammar.length - 1
  }

  addGrammar(from: number, to: number[]): void {
    to.push(0)
    this._grammar[from].push(to)
  }

  parse(s: ArrayLike<number>, init: number): boolean {
    const n = s.length
    const chart: State[][][] = Array(n + 1)
    for (let i = 0; i <= n; i++) {
      chart[i] = Array(this._grammar.length)
      for (let j = 0; j < this._grammar.length; j++) {
        chart[i][j] = []
      }
    }

    this._grammar[1] = [[init, 0]]
    let curr: State[] = [[1, 0, 0, 0]]
    let next: State[] = []
    for (let k = 0; k <= n; k++) {
      for (let j = 0; j < curr.length; j++) {
        const state = curr[j]
        const [sa, sk, sp, si] = state
        const b = this._getSymbol(state)

        if (b) {
          if (!this._terminal[b]) {
            for (let j = 0; j < this._grammar[b].length; j++) {
              this._enqueue(curr, [b, j, 0, k])
            }
          } else if (this._terminal[b] === s[k]) {
            this._enqueue(next, [sa, sk, sp + 1, si])
          }
        } else {
          // eslint-disable-next-line no-loop-func
          chart[si][sa].forEach(t => {
            this._enqueue(curr, [t[0], t[1], t[2] + 1, t[3]])
          })
        }
      }

      curr.forEach(t => {
        chart[k][this._getSymbol(t)].push(t)
      })
      // eslint-disable-next-line semi-style
      ;[curr, next] = [next, curr]
      next = []
    }

    return chart[n][0].some(t => t[0] === 1)
  }

  // eslint-disable-next-line class-methods-use-this
  private _enqueue(cur: State[], state: State): void {
    for (let j = 0; j < cur.length; j++) {
      const [a, k, p, i] = cur[j]
      if (a === state[0] && k === state[1] && p === state[2] && i === state[3]) return
    }
    cur.push(state)
  }

  private _getSymbol(state: State): number {
    return this._grammar[state[0]][state[1]][state[2]]
  }
}

export {}

if (require.main === module) {
  const parser = new EarleyParser()
  const A = parser.addSymbol()
  const B = parser.addSymbol()

  const a = parser.addSymbol(97)
  const b = parser.addSymbol(98)
  parser.addGrammar(A, [B])
  parser.addGrammar(A, [a, A, a])
  parser.addGrammar(A, [b, A, b])
  parser.addGrammar(B, [a])
  parser.addGrammar(B, [b])
  parser.addGrammar(B, [a, a])
  parser.addGrammar(B, [b, b])
  const ords = 'abba'.split('').map(c => c.charCodeAt(0))
  console.log(parser.parse(ords, A))
  console.log(ords)
}
