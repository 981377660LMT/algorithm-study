function badSensor(sensor1: number[], sensor2: number[]): number {
  const n = sensor1.length

  let errorIndex = 0
  while (errorIndex < n) {
    if (sensor1[errorIndex] !== sensor2[errorIndex]) break
    errorIndex++
  }

  if (errorIndex >= n - 1) return -1

  const withoutLast1 = sensor1.slice(errorIndex, n - 1).join('')
  const withoutLast2 = sensor2.slice(errorIndex, n - 1).join('')
  const withoutError1 = sensor1.slice(errorIndex + 1).join('')
  const withoutError2 = sensor2.slice(errorIndex + 1).join('')

  if (withoutLast1 === withoutError2 && withoutError1 === withoutLast2) return -1
  else if (withoutLast1 === withoutError2) return 1
  else return 2
}

console.log(badSensor([2, 3, 4, 5], [2, 1, 3, 4]))
// 输入：sensor1 = [2,3,4,5], sensor2 = [2,1,3,4]
// 输出：1
// 解释：传感器 2 返回了所有正确的数据.
// 传感器2对第二个数据点采集的数据，被传感器1丢弃了，传感器1返回的最后一个数据被替换为 5 。

// 输入：sensor1 = [2,2,2,2,2], sensor2 = [2,2,2,2,5]
// 输出：-1
// 解释：无法判定拿个传感器是有缺陷的。
// 假设任一传感器丢弃的数据是最后一位，那么，另一个传感器就能给出与之对应的输出。

console.log('as'.slice(0, 2))

export {}
