// 请你以字符串形式，返回一个与它相似度最大且可以简写的颜色(最后的结果都是 00 11 22 33 44 55 66 77 88 99 aa bb cc dd ee ff)
// color.length == 7
// color[0] == '#'
// `x11, 0x22, 0x33, 0x44... 都是 0x11 的倍数;00 到 ff 中找到一个相似度最大的;除以17，>8四舍五入`
const OFFSET = 17
function similarRGB(color: string): string {
  const helper = (target: string): string => {
    const hex = parseInt(target, 16)
    let [div, mod] = [~~(hex / OFFSET), hex % OFFSET]
    if (mod > 8) div++
    const decimal = div * OFFSET
    return decimal.toString(16).padStart(2, '0')
  }

  return `#${helper(color.slice(1, 3))}${helper(color.slice(3, 5))}${helper(color.slice(5, 7))}`
}

console.log(similarRGB('#09f166'))
// 输出："#11ee66"
// 解释：
// 因为相似度计算得出 -(0x09 - 0x11)^2 -(0xf1 - 0xee)^2 - (0x66 - 0x66)^2 = -64 -9 -0 = -73
// 这已经是所有可以简写的颜色中最相似的了

// 思路:将offset归一化，然后四舍五入判断最近的值
