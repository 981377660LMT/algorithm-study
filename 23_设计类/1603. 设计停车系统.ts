class ParkingSystem {
  private state: number
  /**
   *
   * @param big
   * @param medium
   * @param small
   * 每种停车位的数目0 <= big, medium, small <= 1000
   */
  constructor(big: number, medium: number, small: number) {
    this.state = big + (medium << 10) + (small << 20)
  }

  /**
   *
   * @param carType
   * 一辆车只能停在  carType 对应尺寸的停车位中。
   * 如果没有空车位，请返回 false ，否则将该车停入车位并返回 true 。
   */
  addCar(carType: number): boolean {
    const mask = 10 * (carType - 1)
    if ((this.state >> mask) & 0x3ff) {
      this.state -= 1 << mask
      return true
    }
    return false
  }
}
// 使用哈希表来进行记录。
// 这样做的好处是，当增加车类型，只需要重载一个构造方法即可。

// 由于 1000 的二进制表示只有 10 位，而 intint 有 32 位
// 我们可以使用一个 int 配合「位运算」来分段做。
// 这样 intint 分段的做法，在工程源码上也有体现：
// JDK 中的 ThreadPoolExecutor 使用了一个 ctl 变量 (int 类型)
// 的前 3 位记录线程池的状态，
// 后 29 位记录程池中线程个数。
// 这样的「二进制分段压缩存储」的主要目的，不是为了减少使用一个 int，
// 而是为了让「非原子性操作」变为「原子性操作」。
// 如果我们将「线程数量」和「线程池的状态」合二为一之后，
// 我们只需要修改一个 int，就可改变线程数量和线程池的状态
// 这时候只需要使用 CAS 做法（用户态）即可保证线程安全与原子性。
