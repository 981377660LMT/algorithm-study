export {}

const INF = 2e15

// # 给你一个下标从 0 开始、长度为 n 的数组 usageLimits 。

// # 你的任务是使用从 0 到 n - 1 的数字创建若干组，并确保每个数字 i 在 所有组 中使用的次数总共不超过 usageLimits[i] 次。此外，还必须满足以下条件：

// # 每个组必须由 不同 的数字组成，也就是说，单个组内不能存在重复的数字。
// # 每个组（除了第一个）的长度必须 严格大于 前一个组。
// # 在满足所有条件的情况下，以整数形式返回可以创建的最大组数。

// # 最多的数字要在开头使用

function maxIncreasingGroups(usageLimits: number[]): number {}

// class Solution:
//     def maxIncreasingGroups(self, usageLimits: List[int]) -> int:
//         usageLimits.sort(reverse=True)
//         bit = BIT2(len(usageLimits))
//         for i, limit in enumerate(usageLimits):
//             bit.add(i, i + 1, limit)
//         res = 0
//         count = 1
//         while True:
//             if bit.query(0, count) < count:
//                 break
//             res += 1
//             bit.add(count - 1, -count)
//             count += 1
