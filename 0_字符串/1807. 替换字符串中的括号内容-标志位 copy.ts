/**
 * @param {string} s
 * @param {string[][]} knowledge
 * @return {string}
 */
function evaluate(s: string, knowledge: string[][]) {
  const mapping = Object.fromEntries(knowledge)

  s = s.replace(/\(([a-z]+)\)/g, function (_, group) {
    return mapping[group] ?? '?'
  })

  return s
}

console.log(
  evaluate('(name)is(age)yearsold', [
    ['name', 'bob'],
    ['age', 'two'],
  ])
)

export {}
