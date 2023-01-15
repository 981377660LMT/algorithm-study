/**
 * @param {string} s
 * @param {string[][]} knowledge
 * @return {string}
 */
function evaluate(s: string, knowledge: string[][]): string {
  const mapping = Object.fromEntries(knowledge)
  return s.replace(/\(([a-z]+)\)/g, (_, group) => mapping[group] ?? '?')
}

console.log(
  evaluate('(name)is(age)yearsold', [
    ['name', 'bob'],
    ['age', 'two']
  ])
)

export {}
