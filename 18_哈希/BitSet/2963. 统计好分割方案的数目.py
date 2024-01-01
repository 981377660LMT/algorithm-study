# 2963. 统计好分割方案的数目
# https://leetcode.cn/problems/count-the-number-of-good-partitions/
# 给你一个下标从 0 开始、由 正整数 组成的数组 nums。
# 将数组分割成一个或多个 连续 子数组，如果不存在包含了相同数字的两个子数组，则认为是一种 好分割方案 。
# 返回 nums 的 好分割方案 的 数目。
# 由于答案可能很大，请返回答案对 109 + 7 取余 的结果。


from typing import List
from collections import defaultdict, Counter

MOD = int(1e9 + 7)


class Solution:
    def numberOfGoodPartitions(self, nums: List[int]) -> int:
        # 判断每个分割线左右是否包含相同数字,哈希表维护交集(注意每次只会变动一个元素).
        left, right = Counter(), Counter(nums)
        intersect = set()
        n = len(nums)
        res = 0
        for i in range(n - 1):
            v = nums[i]
            left[v] += 1
            right[v] -= 1
            if right[v] == 0:
                del right[v]
            if v in left and v in right:
                intersect.add(v)
            else:
                intersect.discard(v)
            if len(intersect) == 0:
                res += 1
        return pow(2, res, MOD)

    def numberOfGoodPartitions2(self, nums: List[int]) -> int:
        # 判断每个分割线左右是否包含相同数字,bitset求交集
        id_ = defaultdict(lambda: len(id_))
        nums = [id_[v] for v in nums]
        leftCounter, rightCounter = Counter(), Counter(nums)
        leftMask, rightMask = 0, sum(1 << v for v in rightCounter)
        res = 0
        for v in nums:
            leftCounter[v] += 1
            if leftCounter[v] == 1:
                leftMask |= 1 << v
            rightCounter[v] -= 1
            if rightCounter[v] == 0:
                rightMask ^= 1 << v
            res += leftMask & rightMask == 0
        return pow(2, res - 1, MOD)


# nums = [1,2,3,4]
print(Solution().numberOfGoodPartitions([1, 2, 3, 4]))
print(Solution().numberOfGoodPartitions([1, 2, 1, 3]))
print(Solution().numberOfGoodPartitions([1]))
print(Solution().numberOfGoodPartitions([1, 1, 1, 1]))
print(Solution().numberOfGoodPartitions([1, 2, 1]))
