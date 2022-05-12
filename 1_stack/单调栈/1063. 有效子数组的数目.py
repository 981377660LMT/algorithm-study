from typing import List

# 子数组中，最左侧的元素不大于其他元素。
# 返回满足下面条件的 非空、连续 子数组的数目：
# 1 <= A.length <= 50000


# 要求输出的子数最左边的元素不大于其他元素，所以我们需要找到每一个元素的从其右边开始第一个小于这个元素的数，
# 那么这两个数的索引差就是，当前元素对象（即上文的每一个元素）作为子数组的最左边值，最大能够延伸的长度

# 比如数组，[2, 3, 4, 1]，可以看到1小于2，所以2作为子数组的开头最大能够延伸到1的前一位，就是到4，所以满足条件的子数组就是[2], [2, 3], [2, 3, 4]。


class Solution:
    def validSubarrays(self, nums: List[int]) -> int:
        nums.append(-int(1e20))
        n = len(nums)
        stack = []
        res = 0

        for i in range(n):
            while stack and nums[stack[-1]] > nums[i]:
                res += i - stack.pop()
            stack.append(i)

        return res


print(Solution().validSubarrays([1, 4, 2, 5, 3]))
