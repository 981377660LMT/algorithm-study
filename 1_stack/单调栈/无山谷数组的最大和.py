# !将数组变为`无山谷数组`后的最大和
# Consider a list of integers A such that A[i] ≤ nums[i]. Also, there are no j and k such that there exist j < i < k and A[j] > A[i] and A[i] < A[k].
# Return the maximum possible sum of A.


# 数组取最大和时 肯定是山脉数组

from typing import List


class Solution:
    def solve(self, nums):
        def getSumAsPeek(nums) -> List[int]:
            """对每个前缀，找到受限制的递增序列在每个位置处的最大值"""
            stack = []
            res = []
            curSum = 0  # 栈中 数*count

            for num in nums:
                count = 1
                while stack and stack[-1][0] > num:
                    curSum -= stack[-1][0] * stack[-1][1]
                    count += stack[-1][1]
                    stack.pop()

                # "run-length encoding"
                stack.append((num, count))
                curSum += num * count
                res.append(curSum)

            return res

        # 每个位置作为山脉顶峰
        pre = getSumAsPeek(nums)
        suf = getSumAsPeek(nums[::-1])[::-1]
        # print(left, right)
        return max(left + right - mid for left, right, mid in zip(pre, suf, nums))


print(Solution().solve(nums=[10, 6, 8]))
