function typeOf(obj) {
  return Object.prototype.toString.call(obj).split(' ')[1].slice(0, -1).toLowerCase()
}

const values = [
  true,
  10,
  'foo',
  Symbol('bar'),
  null,
  undefined,
  1n,
  NaN,
  {},
  [],
  function () {},
  new Error(),
  new Date(),
  /a/,
  new Map(),
  new Set(),
  Promise.resolve(),
  function* () {},
  class {}
]

values.forEach(v => console.log(typeOf(v)))
