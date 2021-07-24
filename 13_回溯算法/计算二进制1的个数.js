/**
 *
 * @param {number} num
 * @returns  {number}
 */
const binaryOnes = num => num.toString(2).match(/1/g).length

console.log(binaryOnes(11))
