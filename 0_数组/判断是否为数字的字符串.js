const isNumeric = s => !isNaN(parseFloat(s)) && isFinite(s)

console.log(isNumeric('12.3'))
console.log(isNumeric('12.3a'))
console.log(isNumeric('12px'))
