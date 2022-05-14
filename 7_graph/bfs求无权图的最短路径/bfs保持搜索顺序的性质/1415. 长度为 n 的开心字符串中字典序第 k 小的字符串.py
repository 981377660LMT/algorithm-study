# 1 <= n <= 10
# 1 <= k <= 100
# 一个 「开心字符串」定义为：
# 仅包含小写字母 ['a', 'b', 'c'].
# 对所有在 1 到 s.length - 1 之间的 i ，满足 s[i] != s[i + 1] （字符串的下标从 1 开始）。


from itertools import islice
from typing import Generator, List


class Solution:
    def getHappyString(self, n: int, k: int) -> str:
        def bt(index: int, pre: str, path: List[str]) -> Generator[str, None, None]:
            if index == n:
                yield ''.join(path)
                return

            for next in ('a', 'b', 'c'):
                if next == pre:
                    continue
                path.append(next)
                yield from bt(index + 1, next, path)
                path.pop()

        iter = bt(0, '', [])
        return next(islice(iter, k - 1, None), '')


print(Solution().getHappyString(n=3, k=9))
# 输出："cab"
# 解释：长度为 3 的开心字符串总共有 12 个 ["aba", "abc", "aca", "acb", "bab", "bac", "bca", "bcb", "cab", "cac", "cba", "cbc"] 。第 9 个字符串为 "cab"

