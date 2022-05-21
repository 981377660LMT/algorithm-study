const buffer = new ArrayBuffer(4 * 1000)
const hashTable1 = new Int32Array(buffer)
// 两个效果一样
const hashTable2 = new Int32Array(1000)
