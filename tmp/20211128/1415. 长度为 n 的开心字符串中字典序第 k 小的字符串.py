# 1 <= n <= 10
# 1 <= k <= 100
# 一个 「开心字符串」定义为：
# 仅包含小写字母 ['a', 'b', 'c'].
# 对所有在 1 到 s.length - 1 之间的 i ，满足 s[i] != s[i + 1] （字符串的下标从 1 开始）。


from typing import List


class Solution:
    def getHappyString(self, n: int, k: int) -> str:
        choose = ['a', 'b', 'c']
        res = []

        def bt(path: List[str]) -> None:
            if len(res) == k:
                return
            if len(path) == n:
                return res.append(''.join(path))
            for char in choose:
                if not path or char != path[-1]:
                    path.append(char)
                    bt(path)
                    path.pop()

        bt([])

        return res[k - 1] if k - 1 < len(res) else ''


print(Solution().getHappyString(n=3, k=9))
# 输出："cab"
# 解释：长度为 3 的开心字符串总共有 12 个 ["aba", "abc", "aca", "acb", "bab", "bac", "bca", "bcb", "cab", "cac", "cba", "cbc"] 。第 9 个字符串为 "cab"

