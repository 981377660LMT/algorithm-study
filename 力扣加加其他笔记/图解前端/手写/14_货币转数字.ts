function numFormat(str: string) {
  return str.replace(/[^\d]/g, '')
}

console.log(numFormat('$1,234'))
