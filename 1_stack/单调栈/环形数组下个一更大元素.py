# 环形数组下个一更大元素
class Solution:
    def solve(self, nums):
        n = len(nums)
        nums = nums + nums

        res = [-1] * n
        stack = []
        for i, num in enumerate(nums):
            while stack and num > stack[-1][0]:
                _, pos = stack.pop()
                if pos < n:
                    res[pos] = num

            stack.append((num, i))

        return res


print(Solution().solve(nums=[3, 4, 0, 2]))

