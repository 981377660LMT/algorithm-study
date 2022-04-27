# 原来按位与运算有递减的性质，及a & b <=a且a & b <=b

# 求子数组的与，使得这个值最接近target，并返回最近距离
# 1 <= arr.length <= 10^5
# 1 <= arr[i] <= 10^6

from typing import List


class Solution:
    def closestToTarget(self, arr: List[int], target: int) -> int:
        """一个数组的所有不同的前缀与和的个数，不会超过第一个数中的 1 的个数，因为每次与上一个新的数，要么值不变，要么消掉当前前缀和中的至少一个 1
        
        set滚动更新解法
        时间复杂度O(nlog(max(arr))) 
        """
        res = abs(arr[0] - target)
        dp = set([arr[0]])
        for num in arr[1:]:
            ndp = {num & x for x in dp} | {num}
            for x in ndp:
                res = min(res, abs(x - target))
            dp = ndp
        return res


print(Solution().closestToTarget(arr=[9, 12, 3, 7, 15], target=5))
