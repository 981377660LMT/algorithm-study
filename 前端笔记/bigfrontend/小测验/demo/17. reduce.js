1, 2
undefined, 3 // note, undefined as we don't return from our reduce
;[1, 2, 3].reduce((a, b) => {
  console.log(a, b)
})
//////////////////////////////////////////////////////////////////////
0, 1
undefined, 2 // undefined as we don't return from our reduce
undefined, 3
;[1, 2, 3].reduce((a, b) => {
  console.log(a, b)
}, 0)
