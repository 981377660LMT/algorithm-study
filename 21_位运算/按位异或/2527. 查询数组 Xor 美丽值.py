# 2527. 查询数组 Xor 美丽值

# 给你一个下标从 0 开始的整数数组 nums 。
# !三个下标 i ，j 和 k 的 有效值 定义为 ((nums[i] | nums[j]) & nums[k]) 。
# !一个数组的 xor 美丽值 是数组中所有满足 0 <= i, j, k < n  的三元组 (i, j, k) 的 有效值 的异或结果。
# 请你返回 nums 的 xor 美丽值。

# 按位统计后就是布尔代数
# !布尔代数中 异或等于加法 与等于乘法

# 脑筋急转弯:
# 假设这一位有k个1
# 那么异或取1时, nums[k]为1,nums[i]与nums[j]不全为0
# !有 k*(n^2-(n-k)*(n-k))种取法 也就是 k*(n^2-n^2+2nk-nk) = (2nk-k^2)*k 种取1的方法
# 所以当k为奇数时,这一位取1,否则取0，这实际上就是看每个比特位的异或值是否为 1。


from functools import reduce
from typing import List


class Solution:
    def xorBeauty(self, nums: List[int]) -> int:
        return reduce(lambda x, y: x ^ y, nums, 0)
