// https://oi-wiki.org/string/automaton/
// 自动机是一个对信号序列进行判定的数学模型

/**
 * 确定有限状态自动机(DFA).转移方式确定，状态数量有限.
 * Moore自动机: 状态机的输出仅与当前状态有关.
 * Mealy自动机: 状态机的输出与当前状态和输入有关.
 *
 * - 状态集State: 自动机的若干状态.包含起始状态、接受状态集合(于此结束则表示可以接受输入).
 * - 字符集Input: 可接受的输入.
 * - 转移函数: 自动机从一个状态转移到另一个状态的转移方式.
 */
interface DFA<State, Input> {
  move(state: State, input: Input): State
  accept(state: State): boolean
}

export type { DFA }
