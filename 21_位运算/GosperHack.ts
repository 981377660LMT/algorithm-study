// Gosper's Hack是一种生成 [公式] 元集合所有 [公式] 元子集的算法，它巧妙地利用了位运算

function GosperHack(n: number, k: number) {
  let x = (1 << k) - 1
  const limit = 1 << n
  while (x < limit) {
    console.log(x.toString(2).padStart(5, '0'))
    const lowbit = x & -x
    const r = x + lowbit
    // xor
    x = r | (((x ^ r) >> 2) / lowbit)
  }
}

GosperHack(5, 3)
// const lowbit = x & -x 标识出 x 最低位的1 e.g. 0b10110 –> 0b00010
// const r = x + lowbit将 x 右端的连续一段1清零 e.g. 0b10110 –> 0b11000
// x = r | (((x ^ r) >> 2) / lowbit)  e.g. 0b11000 | 0b00001 = 0b11001
