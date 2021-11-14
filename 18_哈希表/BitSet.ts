class Bitset {
  private words: Uint8Array

  constructor(nBits: number) {
    const bytes = Math.ceil(nBits / 8)
    this.words = new Uint8Array(bytes)
  }

  add(bitIndex: number): void {
    const row = Math.floor(bitIndex / 8)
    const col = bitIndex % 8
    this.words[row] |= 1 << col
  }

  has(bitIndex: number): boolean {
    const row = Math.floor(bitIndex / 8)
    const col = bitIndex % 8
    return ((this.words[row] >> col) & 1) === 1
  }
}

if (require.main === module) {
  const bitSet = new Bitset(10000)
  bitSet.add(1000)
  console.log(bitSet.has(1000))
  console.log(bitSet.has(1001))
}

export { Bitset }
