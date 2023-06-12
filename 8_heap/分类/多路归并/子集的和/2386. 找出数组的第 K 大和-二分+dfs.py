# 2386. 找出数组的第 K 大和-二分+dfs
# n == nums.length
# 1 <= n <= 1e5
# -1e9 <= nums[i] <= 1e9
# 1 <= k <= min(2000, 2^n)

# https://leetcode.cn/problems/find-the-k-sum-of-an-array/solution/zhuan-huan-dui-by-endlesscheng-8yiq/

# 0.将数组中所有非负数的和即为posSum.
# !1.对nums所有数取绝对值。任意一个子集的和可以统一成从posSum中减去某些数(此时所有的数是非负数)
# 2.对绝对值数组排序,二分求出从posSum中减去的`第k-1小`的子序列和，就是原数组的第k大子序列和。


from typing import List


class Solution:
    def kSum(self, nums: List[int], k: int) -> int:
        """
        二分+dfs求出数组的第k大子序列和.
        时间复杂度O(nlogn+klogU).
        """
        n = len(nums)
        posSum = sum([v for v in nums if v > 0])
        absNums = sorted([abs(v) for v in nums])

        def check(mid: int) -> bool:
            """从absNums里选若干个数使得和不超过mid,选法是否不小于k-1种."""

            def dfs(index: int, curSum: int) -> None:
                if index == n:
                    return
                nonlocal count
                if count >= k - 1:  # !限制dfs次数
                    return
                cand = curSum + absNums[index]
                if cand <= mid:  # !剪枝,后面的都无法选择了
                    count += 1
                    dfs(index + 1, cand)  # !先进入选择的分支，再进入不选择的分支
                    dfs(index + 1, curSum)

            count = 0
            dfs(0, 0)
            return count >= k - 1

        left, right = 0, sum(absNums)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return posSum - left


if __name__ == "__main__":
    # nums = [2,4,-2], k = 5
    print(Solution().kSum(nums=[2, 4, -2], k=5))
