function isalpha(str: string): boolean {
  if (typeof str !== 'string') return false
  const code = str.codePointAt(0)
  if (code == undefined) return false
  return (code >= 65 && code <= 90) || (code >= 97 && code <= 122)
}

function isnumeric(str: string) {
  if (typeof str !== 'string') return false
  return !Number.isNaN(parseFloat(str)) && Number.isFinite(str)
}

export { isalpha, isnumeric }
