# 每 两个 连续竖线 '|' 为 一对 。换言之，第一个和第二个 '|' 为一对，第三个和第四个 '|' 为一对，以此类推。
# 请你返回 不在 竖线对之间，s 中 '*' 的数目。


class Solution:
    def countAsterisks(self, s: str) -> int:
        """标志位简化逻辑"""
        res = []
        isOk = True
        for char in s:
            if char == '|':
                isOk = not isOk
            elif isOk:
                res.append(char)
        return res.count('*')
        return sum(p.count('*') for p in s.split('|')[::2])

