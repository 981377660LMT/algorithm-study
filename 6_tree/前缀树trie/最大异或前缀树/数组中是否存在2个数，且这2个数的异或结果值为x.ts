// 数组中是否存在2个数，且这2个数的异或结果值为x
// a^k=b等价于a^b=k:即对每一个数a,只需在数组中寻找是否存在另一个数等于a^k即可
const xorEqualToK = (nums: number[], k: number) => {
  const visited = new Set<number>()
  for (const num of nums) {
    if (visited.has(num)) return [num ^ k, num]
    visited.add(num ^ k)
  }
  return []
}

console.log(xorEqualToK([1, 2, 3, 4], 7))
console.log(3 ^ 4)
console.log(3 ^ 7)

// 暴力法：

// for(int i = 1; i < nums.length; ++i)
//     for(int j = 0; j < i; ++j)
//         if(nums[i] ^ nums[j] == x)
//             return true;

// 不用暴力：
// 假设数组中存在a和b，满足a ^ b == x，则a ^ x == b与b ^ x == a也必然满足。
// 反过来讲，你可以遍历数组nums，将当前遍历的数记为a，计算a ^ x，
// 若数组中存在1个数(记为b)，满足b == a ^ x，则说明数组中存在2个数a和b，
// 满足a ^ b == x。 理解了这一点的话，
// 就可利用HashSet，将nums数组中所有数值都存入HashSet，
