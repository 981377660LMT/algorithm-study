const obj = {
  foo: 'bar',
}

console.log('foo' in obj) // true
console.log(['foo'] in obj) // true - ['foo'] converts to string
