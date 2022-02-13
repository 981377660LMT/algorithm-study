class BitSet {
  private capacity: number
  private bit: bigint
  private size: number

  constructor(capacity: number) {
    this.capacity = capacity
    this.bit = 0n
    this.size = 0
  }

  add(index: number): void {
    if (!(this.bit & (1n << BigInt(index)))) {
      this.size++
      this.bit |= 1n << BigInt(index)
    }
  }

  has(index: number): boolean {
    return (this.bit & (1n << BigInt(index))) === 1n
  }

  delete(index: number): void {
    if (this.bit & (1n << BigInt(index))) {
      this.size--
      this.bit &= ~(1n << BigInt(index))
    }
  }

  flip(): void {
    this.bit ^= (1n << BigInt(this.capacity)) - 1n
    this.size = this.capacity - this.size
  }

  all(): boolean {
    return this.size === this.capacity
  }

  one(): boolean {
    return this.size !== 0
  }

  count(): number {
    return this.size
  }

  toString(): string {
    return this.bit.toString(2).padStart(this.capacity, '0')
  }
}

export { BitSet }
