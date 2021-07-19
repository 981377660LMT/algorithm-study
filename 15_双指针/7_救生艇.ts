// 每艘船最多可同时载两人，但条件是这些人的重量之和最多为 limit。
const numRescueBoats = (people: number[], limit: number) => {
  people.sort((a, b) => a - b)
  let l = 0
  let r = people.length - 1
  let res = 0

  while (l <= r) {
    if (people[l] + people[r] > limit) {
      r--
      res++
    } else {
      r--
      l++
      res++
    }
  }

  return res
}
// 6
console.log(numRescueBoats([44, 10, 29, 12, 49, 41, 23, 5, 17, 26], 50))
console.log(numRescueBoats([3, 5, 3, 4], 5))
console.log(numRescueBoats([3, 8, 7, 1, 4], 9))

export {}
