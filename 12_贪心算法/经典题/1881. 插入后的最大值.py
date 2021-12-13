# 返回插入操作后，用字符串表示的 n 的最大值。
class Solution:
    def maxValue(self, s: str, x: int) -> str:
        if s[0] == '-':
            idx = 1
            while idx < len(s) and x >= int(s[idx]):
                idx += 1
        else:
            idx = 0
            while idx < len(s) and x <= int(s[idx]):
                idx += 1
        res = s[:idx] + str(x) + s[idx:]
        return res


print(Solution().maxValue(s="-13", x=2))
# 输出："-123"
# 解释：向 n 中插入 x 可以得到 -213、-123 或者 -132 ，三者中最大的是 -123

