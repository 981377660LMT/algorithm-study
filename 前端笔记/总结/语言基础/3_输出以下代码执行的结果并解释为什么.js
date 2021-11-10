const obj = {
  2: 3,
  3: 4,
  length: 2,
  splice: Array.prototype.splice,
  push: Array.prototype.push,
}

obj.push(1)
obj.push(2)
console.log(obj)
// {
//   '2': 1,
//   '3': 2,
//   length: 4,
//   splice: [Function: splice],
//   push: [Function: push]
// }
// 执行obj.push(1)时，当前length为2，正好替换了obj['2']的值，然后length变为3，obj.push(2)时就是替换obj['3']的值。就出来了浏览器的运行结果
