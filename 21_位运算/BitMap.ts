// 假设要表示的数据集合的容量是n，那么申请的 bit 位是10*n。
const BYTE_SIZE = 8 // 1byte = 8bit

class Bitset {
  private readonly bytes: Buffer
  private readonly capacity: number

  constructor(bytes = 0) {
    this.bytes = Buffer.alloc(bytes) // 申请bytes个字节大小
    this.capacity = bytes * BYTE_SIZE
  }

  /**
   * 作用：将第byteIndex个数对应的字节位bitIndex，设置为1
   * 思路：最简单的8个bit可以表示[0, 7]共8个数字，分别是：
   *   0000 0001（表示0）
   *   0000 0010（表示1）
   *   ......
   */
  add(num: number) {
    if (num > this.capacity) return
    const row = ~~(num / BYTE_SIZE)
    const col = num % BYTE_SIZE
    // 对应的字节的bit位置位0
    this.bytes[row] = this.bytes[row] | (1 << col)
  }

  has(num: number) {
    if (num > this.capacity) return
    const row = ~~(num / BYTE_SIZE)
    const col = num % BYTE_SIZE
    return (this.bytes[row] & (1 << col)) !== 0
  }
}

if (require.main === module) {
  const bitmap = new Bitset(10)
  bitmap.add(76)
  console.log(bitmap.has(76))
  console.log(bitmap)
}

export { Bitset }
