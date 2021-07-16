"use strict";
console.log(~~(2 / 3));
console.log(~~(44 / 3));
console.log(~(42 / 3));
console.log(~~(42 / 3));
// 对于浮点数，~~value可以代替parseInt(value)，而且前者效率更高些
console.log(~~'0', !!'0');
