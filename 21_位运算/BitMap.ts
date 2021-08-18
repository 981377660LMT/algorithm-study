// 假设要表示的数据集合的容量是n，那么申请的 bit 位是10*n。
const BYTE_SIZE = 8 // 1byte = 8bit

class BitMap {
  private bytes: Buffer
  private size: number

  constructor(bytes = 0) {
    this.bytes = Buffer.alloc(bytes) // 申请bytes个字节大小
    this.size = bytes * BYTE_SIZE
  }

  /**
   * 作用：将第byteIndex个数对应的字节位bitIndex，设置为1
   * 思路：最简单的8个bit可以表示[0, 7]共8个数字，分别是：
   *   0000 0001（表示0）
   *   0000 0010（表示1）
   *   ......
   */
  set(num: number) {
    if (num > this.size) return
    const byteIndex = ~~(num / BYTE_SIZE)
    const bitIndex = num % BYTE_SIZE
    // 对应的字节的bit位置位0
    this.bytes[byteIndex] = this.bytes[byteIndex] | (1 << bitIndex)
  }

  has(num: number) {
    if (num > this.size) return
    const byteIndex = ~~(num / BYTE_SIZE)
    const bitIndex = num % BYTE_SIZE
    return (this.bytes[byteIndex] & (1 << bitIndex)) !== 0
  }
}

if (require.main === module) {
  const bitmap = new BitMap(10)
  bitmap.set(76)
  console.log(bitmap.has(76))
  console.log(bitmap)
}

export { BitMap }
