function multi7(n: number) {
  const str7 = n.toString(7)
  return parseInt(`${str7}0`, 7)
}

// 轉乘 7 進位向左位移， 再轉回 10 進位
console.log(multi7(7))
