// 'rgb(255, 255, 255)' -> '#FFFFFF' 的多种思路
function rgb2Hex(rgb: string) {
  const toHex = (str: string) => Number(str).toString(16).toUpperCase().padStart(2, '0')
  const match = rgb.match(/\d+/g)
  if (!match) return ''
  const r = match[0]
  const g = match[1]
  const b = match[2]
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

console.log(rgb2Hex('rgb(255, 255, 255)'))
console.log(rgb2Hex('rgb(255, 1, 255)'))

function hexToRgba(hex: string): string {
  const isValidChars = /^#[a-fA-F\d]+$/.test(hex)
  const isValidLength = [4, 5, 7, 9].includes(hex.length)
  if (!isValidChars || !isValidLength) throw new Error('Invalid HEX')

  const arr = hex
    .toLowerCase()
    .split('')
    .slice(1)
    .reduce((pre, cur) => `${pre}${hex.length < 7 ? cur.repeat(2) : cur}`, '')
    .match(/(\w{2})/g)!
    .map(num => parseInt(num, 16))

  const [r, g, b, a = 255] = arr

  return `rgba(${r},${g},${b},${Math.round((a * 100) / 255) / 100})`
}

console.log(hexToRgba('#fff'))
// 'rgba(255,255,255,1)'
export {}
