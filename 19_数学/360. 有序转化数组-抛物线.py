from typing import List

# 要注意，返回的这个数组必须按照 升序排列，并且我们所期望的解法时间复杂度为 O(n)。
# 计算函数值 f(x) = ax2 + bx + c

# 1、抛物线中心轴 b / -(2a)
# 2、分别求出中心轴两侧数组，然后合并就行了


class Solution:
    def sortTransformedArray(self, nums: List[int], a: int, b: int, c: int) -> List[int]:
        if a == 0:
            res = [b * x + c for x in nums]
            return res if b >= 0 else res[::-1]

        mid = b / -(2 * a)
        right = [a * x * x + b * x + c for x in nums if x > mid]
        left = [a * x * x + b * x + c for x in nums if x <= mid]
        if a > 0:
            left = left[::-1]
        if a < 0:
            right = right[::-1]

        # 合并两个有序数组(也可deque实现)
        res = []
        i, j = 0, 0
        while i < len(right) and j < len(left):
            if right[i] < left[j]:
                res.append(right[i])
                i += 1
            else:
                res.append(left[j])
                j += 1
        if i < len(right):
            res += right[i:]
        if j < len(left):
            res += left[j:]
        return res


print(Solution().sortTransformedArray(nums=[-4, -2, 2, 4], a=1, b=3, c=5))
# 输出: [3,9,15,33]
