function isalpha(s: string): boolean {
  const len = s.length
  if (len === 0) return false
  for (let i = 0; i < len; i++) {
    const cc = s.charCodeAt(i)
    // 'A'..'Z' 是 65..90， 'a'..'z' 是 97..122
    // 如果在 65..122 之外，或者在 90..97 之间，则不是字母
    if (cc < 65 || cc > 122 || (cc > 90 && cc < 97)) {
      return false
    }
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
