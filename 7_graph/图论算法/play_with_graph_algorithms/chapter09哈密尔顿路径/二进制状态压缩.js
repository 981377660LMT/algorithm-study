// visited数组可以用一个属来表示
// 例如 0 1 2 3四个点
// 0 2访问过 1 3没访问过
// 即为[1,0,1,0]
// 二进制数1010
console.log(0b0101)

// 看第0位是否为1
console.log(0b0101 & 0b0001)
// 等于0b0001
// 看第1位是否为1
console.log(0b0101 & 0b0010)
// 等于0b0000

// 看visited 第i位是否为1 2 ^ i即1左移i位
// visited & (2 ^ i) === 1
// visited & 1<<i === 1

// 如果visited第i位为0，设为1：
// visited+(1<<i)

// 如果visited第i位为1，设为0：
// visited-(1<<i)
