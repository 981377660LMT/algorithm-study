from typing import List

# 1 <= nums.length <= 100
# 1 <= n <= 1e5


class Solution:
    def elementInNums(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        """
        每过一分钟，数组的 最左边 的元素将被移除，直到数组为空。
        然后，每过一分钟，数组的 尾部 将添加一个元素
        """
        n = len(nums)
        res = [-1] * len(queries)
        for i, (time, pos) in enumerate(queries):
            time %= 2 * n
            if time < n and pos < n - time:
                # 数组长度为 n-time
                res[i] = nums[time + pos]
            elif time > n and pos < time - n:
                # 数组长度为 time-n
                res[i] = nums[pos]
        return res


# 如果在时刻 timej，indexj < nums.length，那么答案是此时的 nums[indexj]；
# 如果在时刻 timej，indexj >= nums.length，那么答案是 -1。

print(Solution().elementInNums(nums=[0, 1, 2], queries=[[0, 2], [2, 0], [3, 2], [5, 0]]))
print(Solution().elementInNums(nums=[2], queries=[[0, 0], [1, 0], [2, 0], [3, 0]]))

