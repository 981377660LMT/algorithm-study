// 使得每个组包含正好的 K 个字符，但第一个组可能比 K 短，但仍必须包含至少一个字符
// 必须在两个组之间插入短划线，并且所有小写字母都应转换为大写

const licenseKeyFormatting = (key: string, k: number) =>
  key
    .split('')
    .filter(v => v !== '-')
    .map(v => v.toUpperCase())
    .reverse()
    .reduce((pre, cur, index) => pre + (index % k === 0 ? '-' + cur : cur), '')
    .slice(1)
    .split('')
    .reverse()
    .join('')

console.log(licenseKeyFormatting('5F3Z-2e-9-w', 4))
// ('5F3Z-2E9W')
console.log(licenseKeyFormatting('2-5g-3-J', 2))
// ('2-5G-3J')

export {}
