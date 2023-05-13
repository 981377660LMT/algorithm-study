# 如果一个字符串满足以下条件，则称其为 美丽字符串 ：

# 它由英语小写字母表的前 k 个字母组成。
# 不等于前一个字符
# 给你一个长度为 n 的美丽字符串 s 和一个正整数 k 。
# 请你找出并返回一个长度为 n 的美丽字符串，该字符串还满足：在字典序大于 s 的所有美丽字符串中字典序最小。如果不存在这样的字符串，则返回一个空字符串。
# 对于长度相同的两个字符串 a 和 b ，如果字符串 a 在与字符串 b 不同的第一个位置上的字符字典序更大，则字符串 a 的字典序大于字符串 b 。
# 例如，"abcd" 的字典序比 "abcc" 更大，因为在不同的第一个位置（第四个字符）上 d 的字典序大于 c 。

# 2663. 字典序最小的美丽字符串
# !不存在长度>=2的回文子串等价于不存在长2和3的回文串(ARC考过一次)
# !记录前两个字符，然后dfs


class Solution:
    def smallestBeautifulString(self, s: str, k: int) -> str:
        def dfs(pos: int, isLimit: bool, pre1: int, pre2: int):
            if pos == n:
                if not isLimit:
                    yield res
                return
            lower = ords[pos] if isLimit else 97
            for cur in range(lower, 97 + k):
                if cur == pre1 or cur == pre2:
                    continue
                res.append(chr(cur))
                yield from dfs(pos + 1, (isLimit and cur <= ords[pos]), cur, pre1)
                res.pop()

        res = []
        n = len(s)
        ords = list(map(ord, s))
        res = next(dfs(0, True, 0, -1), [])
        return "".join(res)


print(Solution().smallestBeautifulString(s="dcb" * 100, k=4))
