from typing import List


class Solution:
    def bestRotation(self, nums: List[int]) -> int:
        # diff[k]:表示移动K步后，可以产生贡献 其中k<=len(nums)-1 diff数组要多开1个位置
        n = len(nums)
        diff = [0] * (n + 10)
        for i, num in enumerate(nums):
            if num > i:
                # print(i + 1, n + i - num + 1)
                diff[i + 1] += 1
                diff[i + 1 + n - num] -= 1
            else:
                # print(i - num + 1, i + 1)
                diff[0] += 1
                diff[n] -= 1
                diff[i - num + 1] -= 1
                diff[i + 1] += 1
        for i in range(1, len(nums)):
            diff[i] += diff[i - 1]

        return diff.index(max(diff))


print(Solution().bestRotation([2, 3, 1, 4, 0]))
# 输出：3
# 解释：
# 下面列出了每个 K 的得分：
# K = 0,  A = [2,3,1,4,0],    score 2
# K = 1,  A = [3,1,4,0,2],    score 3
# K = 2,  A = [1,4,0,2,3],    score 3
# K = 3,  A = [4,0,2,3,1],    score 4
# K = 4,  A = [0,2,3,1,4],    score 3
# 所以我们应当选择 K = 3，得分最高。
