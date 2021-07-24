console.log(Symbol.for('a') === Symbol.for('a'))

const s = new Set([[1], [1], [1]].map(v => JSON.stringify(v)))
console.log(
  new Set(
    ['John', 'johnsmith@mail.com', 'john00@mail.com'].concat([
      'John',
      'johnsmith@mail.com',
      'john_newyork@mail.com',
    ])
  )
)
