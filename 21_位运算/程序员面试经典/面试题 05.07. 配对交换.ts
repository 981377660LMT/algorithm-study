// 交换某个整数的奇数位和偶数位，尽量使用较少的指令
// （也就是说，位0与位1交换，位2与位3交换，以此类推）。

// int -> 32位
// 奇数位全1 -> 0101.... 表示为 0x55555555
// 偶数位全1 -> 1010.... 表示为 0xaaaaaaaa
// ans = (提取奇数位 << 1) + (提取偶数位 >> 1)
function exchangeBits(num: number): number {
  const evenMask = 0xaaaaaaaa
  const oddMask = 0x55555555
  return ((num & evenMask) >> 1) | ((num & oddMask) << 1)
}

console.log(exchangeBits(0b101011))
