from math import factorial

# 康托展开
# 给定 n 和 k，返回第k个排列。
class Solution:
    def getPermutation(self, n: int, k: int) -> str:
        nums = list(range(1, n + 1))
        sb = []
        k -= 1  # 第一个排列，索引为0

        # 如 n=5,x=62 时：
        # 用 61 / 4! = 2 余 13，说明 a[5]=2,说明比首位小的数有 2 个，所以首位为 3。
        for i in range(n, 0, -1):
            div, mod = divmod(k, factorial(i - 1))
            k = mod
            sb.append(str(nums[div]))
            nums.pop(div)

        return ''.join(sb)


print(Solution().getPermutation(n=3, k=3))
# n = 3, k = 3
