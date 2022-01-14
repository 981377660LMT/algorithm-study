# 我们称一个字符串是 交替 的，需要满足`任意相邻字符都不同`。
# 你可以按任意顺序执行以下两种操作任意次：

# 类型 1 ：删除 字符串 s 的第一个字符并将它 添加 到字符串结尾。
# 类型 2 ：选择 字符串 s 中任意一个字符并将该字符 反转 ，也就是如果值为 '0' ，则反转得到 '1' ，反之亦然。

# 请你返回使 s 变成 交替 字符串的前提下， `类型 2 的 最少 操作次数` 。

INF = 0x3F3F3F3F
pattern1 = '01'


class Solution:
    def minFlips(self, s: str) -> int:
        n = len(s)
        res = INF
        opt1 = 0
        s += s
        for i, char in enumerate(s):
            if char != pattern1[i & 1]:
                opt1 += 1
            if i >= n:
                if s[i - n] != pattern1[(i - n) & 1]:
                    opt1 -= 1
            if i >= n - 1:
                res = min(res, opt1, n - opt1)
        return res


print(Solution().minFlips(s="111000"))
# 输出：2
# 解释：执行第一种操作两次，得到 s = "100011" 。
# 然后对第三个和第六个字符执行第二种操作，得到 s = "101010" 。
