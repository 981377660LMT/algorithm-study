# 返回插入操作后，用字符串表示的 n 的最大值。
class Solution:
    def maxValue(self, s: str, x: int) -> str:
        if s[0] == '-':
            cand = 1
            while cand < len(s) and x >= int(s[cand]):
                cand += 1
        else:
            cand = 0
            while cand < len(s) and x <= int(s[cand]):
                cand += 1
        res = s[:cand] + str(x) + s[cand:]
        return res


print(Solution().maxValue(s="-13", x=2))
# 输出："-123"
# 解释：向 n 中插入 x 可以得到 -213、-123 或者 -132 ，三者中最大的是 -123

