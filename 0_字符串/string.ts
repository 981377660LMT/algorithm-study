function isalpha(str: string): boolean {
  if (typeof str !== 'string') return false
  for (let index = 0; index < str.length; index++) {
    const code = str[index].codePointAt(0)
    if (code == undefined) return false
    if (!((code >= 97 && code <= 122) || (code >= 65 && code <= 90))) return false
  }
  return true
}

function isnumeric(str: string) {
  const num = parseFloat(str)
  return !Number.isNaN(num) && Number.isFinite(num)
}

function isdigit(str: string) {
  return /^\d+$/.test(str)
}

if (require.main === module) {
  console.log(isalpha('a'))
  console.log(isalpha('Aaa'))
  console.log(isnumeric('1e5'))
  console.log(isnumeric('-1'))
  console.log(isnumeric('-'))
  console.log(isdigit('1e5'))
  console.log(isdigit('123'))
}

export { isalpha, isnumeric, isdigit }
