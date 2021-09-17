type Hex = `#${string}`
type RGB = `rgb(${number},${number},${number})`

const hexToRGB = (hex: Hex): RGB => {
  const tmp = parseInt(hex.slice(1), 16)
  const r = tmp >> 16
  const g = (tmp >> 8) & 0xff
  const b = tmp & 0xff
  return `rgb(${r},${g},${b})`
}

const rgbToHex = (rgb: RGB): Hex => {
  const rgbArr = rgb.match(/\d+/g)!.map(Number)
  const color = (rgbArr[0] << 16) | (rgbArr[1] << 8) | rgbArr[2]
  return `#${color.toString(16)}`
}

console.log(hexToRGB(`#ffffff`))

console.log(rgbToHex(`rgb(255,255,255)`))
