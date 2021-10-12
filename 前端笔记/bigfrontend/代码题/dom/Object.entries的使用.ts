console.log(Object.entries([1, 2, 3])) // [ [ '0', 1 ], [ '1', 2 ], [ '2', 3 ] ]
console.log(Object.entries(new Map([[1, 2]]))) // []  注意这里
console.log(Object.entries({ 1: 2 })) // [ [ '1', 2 ] ]
