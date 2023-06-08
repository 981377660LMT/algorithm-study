/**
 * 长这样的都叫状态机.
 */
interface Automaton<V> {
  next(state: number, newValue: V): number
  accept?(state: number): boolean
  readonly size?: number
}

// !DFA 和 NFA 有没有接口表示???
