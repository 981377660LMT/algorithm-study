// 'rgb(255, 255, 255)' -> '#FFFFFF' 的多种思路
const rgb2Hex = (rgb: string) => {
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
