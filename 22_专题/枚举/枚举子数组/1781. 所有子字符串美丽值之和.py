# 一个字符串的 美丽值 定义为：出现频率最高字符与出现频率最低字符的出现次数之差。


# 比方说，"abaacc" 的美丽值为 3 - 1 = 2 。
# 给你一个字符串 s ，请你返回它所有子字符串的 美丽值 之和。
# 1 <= s.length <= 500
# s 只包含小写英文字母。


class Solution:
    def beautySum(self, s: str) -> int:
        res = 0
        n = len(s)
        for i in range(n):
            counter = [0] * 26
            for j in range(i, n):
                counter[ord(s[j]) - ord("a")] += 1
                res += max(counter) - min(x for x in counter if x)
        return res
