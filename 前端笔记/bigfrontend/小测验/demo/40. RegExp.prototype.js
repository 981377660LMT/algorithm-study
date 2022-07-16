console.log(/^4\d\d$/.test('404')) // true
console.log(/^4\d\d$/.test(404)) // true - converts to string
console.log(/^4\d\d$/.test(['404'])) // true - converts to string
console.log(/^4\d\d$/.test([404])) // true - converts to string
