# 删除子字符串 "ab" 并得到 x 分。
# 删除子字符串"ba" 并得到 y 分
# 请返回对 s 字符串执行上面操作若干次能得到的最大得分。
class Solution:
    def maximumGain(self, s: str, x: int, y: int) -> int:

        if y > x:
            s = s[::-1]
            x, y = y, x

        res = 0

        # 讨论x>=y的情况 多匹配'ab'
        stack1 = []
        for char in s:
            if char == 'b' and stack1 and stack1[-1] == 'a':
                stack1.pop()
                res += x
            else:
                stack1.append(char)

        stack2 = []
        for char in stack1:
            if char == 'a' and stack2 and stack2[-1] == 'b':
                stack2.pop()
                res += y
            else:
                stack2.append(char)

        return res


print(Solution().maximumGain(s="cdbcbbaaabab", x=4, y=5))
# 输出：19
# 解释：
# - 删除 "cdbcbbaaabab" 中加粗的 "ba" ，得到 s = "cdbcbbaaab" ，加 5 分。
# - 删除 "cdbcbbaaab" 中加粗的 "ab" ，得到 s = "cdbcbbaa" ，加 4 分。
# - 删除 "cdbcbbaa" 中加粗的 "ba" ，得到 s = "cdbcba" ，加 5 分。
# - 删除 "cdbcba" 中加粗的 "ba" ，得到 s = "cdbc" ，加 5 分。
# 总得分为 5 + 4 + 5 + 5 = 19 。

