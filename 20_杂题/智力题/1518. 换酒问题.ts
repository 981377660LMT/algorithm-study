// 小区便利店正在促销，用 numExchange 个空酒瓶可以兑换一瓶新酒。
// 你购入了 numBottles 瓶酒。
// 如果喝掉了酒瓶中的酒，那么酒瓶就会变成空的。
// 请你计算 最多 能喝到多少瓶酒。

// !不可以借空瓶

function numWaterBottles1(numBottles: number, numExchange: number): number {
  return numBottles + Math.floor((numBottles - 1) / (numExchange - 1))
}

function numWaterBottles2(numBottles: number, numExchange: number): number {
  let res = numBottles
  let empty = numBottles
  while (empty >= numExchange) {
    const count = Math.floor(empty / numExchange)
    res += count
    empty -= count * numExchange
    empty += count
  }

  return res
}

console.log(numWaterBottles1(9, 3))
console.log(numWaterBottles2(9, 3))
// 3个空瓶=1个啤酒=1个空瓶+1单位酒 --> 2个空瓶=1单位酒 --> 1个空瓶=0.5单位酒
// 那么该人总共可以喝 公式A:n+n/(m-1)=9+9/(3-1)=13瓶，但是这个答案是不完美的。
// 在实际问题中我们最后不可能恰好喝完手上又没有空瓶的，如果恰好喝完请问最后喝的一杯酒是空手接白刃吗？
// 所以我们要对n减去1
