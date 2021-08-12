/**
 * @param {number} turnedOn
 * @return {string[]}
 * @description 二进制手表顶部有 4 个 LED 代表 小时（0-11），底部的 6 个 LED 代表 分钟（0-59）。每个 LED 代表一个 0 或 1，最低位在右侧。
 */
const readBinaryWatch = function (turnedOn: number): string[] {
  const res: string[] = []
  for (let h = 0; h < 12; h++) {
    for (let m = 0; m < 60; m++) {
      if (
        h.toString(2).replace(/0/g, '').length + m.toString(2).replace(/0/g, '').length ===
        turnedOn
      ) {
        res.push(`${h}:${m.toString().padStart(2, '0')}`)
      }
    }
  }

  return res
}

console.log(readBinaryWatch(1))
// ["0:01","0:02","0:04","0:08","0:16","0:32","1:00","2:00","4:00","8:00"]
// 小时不会以零开头：分钟必须由两位数组成
