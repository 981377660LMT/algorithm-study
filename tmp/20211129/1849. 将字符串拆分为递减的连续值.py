# 请你判断能否将 s 拆分成两个或者多个 非空子字符串 ，使子字符串的 数值 按 降序 排列，且每两个 相邻子字符串 的数值之 差 等于 1 。
# 1 <= s.length <= 20


# 总结:
# 1. 枚举字符串分割点
# 2. cur index 来回溯


class Solution:
    def splitString(self, s: str) -> bool:
        def bt(cur: int, preNum: int) -> bool:
            if cur >= len(s):
                return True
            for i in range(cur + 1, len(s) + 1):
                if int(s[cur:i]) == preNum - 1:
                    if bt(i, preNum - 1):
                        return True
            return False

        s = s.lstrip('0')
        for i in range(1, len(s)):
            if bt(i, int(s[:i])):
                return True
        return False


# print(Solution().splitString("050043"))
print(Solution().splitString("1098765432"))
# 输出：true
# 解释：s 可以拆分为 ["05", "004", "3"] ，对应数值为 [5,4,3] 。
# 满足按降序排列，且相邻值相差 1 。
