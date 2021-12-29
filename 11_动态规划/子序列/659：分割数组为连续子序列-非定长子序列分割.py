# 注意与1296. 划分数组为连续数字的集合的区别
# 给你一个按升序排序的整数数组 num（可能包含重复数字）， 1 <= nums.length <= 10000
# 请你将它们分割成`一个或多个长度至少为 3` 的子序列

# 贪心
from typing import List
from collections import Counter


class Solution:
    def isPossible(self, nums: List[int]) -> bool:
        remain = Counter(nums)
        # 记录以key结尾的字序列个数
        end = Counter()

        for num in nums:
            if remain[num] == 0:
                continue

            remain[num] -= 1

            if end[num - 1] > 0:
                end[num - 1] -= 1
                end[num] += 1
            elif remain[num + 1] > 0 and remain[num + 2] > 0:
                remain[num + 1] -= 1
                remain[num + 2] -= 1
                end[num + 2] += 1
            else:
                return False
        return True


print(Solution().isPossible([1, 2, 3, 3, 4, 5]))
print(Solution().isPossible([1, 2, 3, 3, 4, 4, 5, 5]))
