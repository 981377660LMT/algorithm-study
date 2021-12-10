console.log(Array.from('3*4-2*5+'.matchAll(/(\()|(\d+)|([-+*/])|(\))/g)))
const matchArray = Array.from('"2-3/(5*2)+1+'.matchAll(/(\()|(\d+)|([-+*/])|(\))/g))
for (const [_, leftBracket, num, opt, rightBracket] of matchArray) {
  console.log(leftBracket, num, opt, rightBracket)
}
