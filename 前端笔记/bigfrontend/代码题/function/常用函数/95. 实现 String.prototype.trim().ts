/**
 * @param {string} str
 * @return {string}
 */
function trim(str: string): string {
  return str.replace(/^\s+|\s+$/g, '')
}

console.log(trim('   df  '))
