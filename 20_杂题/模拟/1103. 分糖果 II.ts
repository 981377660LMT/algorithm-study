function distributeCandies(candies: number, num_people: number): number[] {
  const res = Array<number>(num_people).fill(0)
  let index = 0

  while (candies > 0) {
    const cur = Math.min(index + 1, candies)
    res[index % num_people] += cur
    candies -= cur
    index++
  }

  return res
}

// 给第一个小朋友 1 颗糖果，第二个小朋友 2 颗，依此类推，直到给最后一个小朋友 n 颗糖果。
// 然后，我们再回到队伍的起点，给第一个小朋友 n + 1 颗糖果，第二个小朋友 n + 2 颗，依此类推，直到给最后一个小朋友 2 * n 颗糖果。
