# 972. 相等的有理数
# 整数、小数、循环节都不大于4位
# 也就是说，两个小数的循环部分的最大的 公共周期 为 12。
class Solution:
    def isRationalEqual(self, s: str, t: str) -> bool:
        def f(s):
            i = s.find("(")
            if i >= 0:
                s = s[:i] + s[i + 1 : -1] * 20
            return float(s[:20])

        return f(s) == f(t)


print(Solution().isRationalEqual(s="0.(52)", t="0.5(25)"))
# 输出：true
# 解释：因为 "0.(52)" 代表 0.52525252...，
# 而
# "0.5(25)" 代表 0.52525252525.....，则这两个字符串表示相同的数字。
