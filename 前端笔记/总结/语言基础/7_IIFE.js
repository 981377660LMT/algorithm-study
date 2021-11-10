// IIFE内的var穿透了块作用域，name被提升至if()之前，且此时name为undefined。
var name = 'Tom'
;(function () {
  if (typeof name == 'undefined') {
    var name = 'Jack'
    console.log('Goodbye ' + name)
  } else {
    console.log('Hello ' + name)
  }
})()
