from itertools import islice
from typing import Generator, List


# 2^n
class Solution:
    def solve(self, n: int, k: int) -> str:
        """
        返回'0''1''2'组成的长为n的字典序的第k个字符串 相邻字符不能相同
        用dfs搜，搜出来直接就是字典序，并且用生成器可以节省空间，加速
        如果用bfs搜，搜出来是实际大小排序
        """

        def bt(index: int, pre: int, path: List[int]) -> Generator[str, None, None]:
            if index == n:
                yield ''.join(map(str, path))
                return

            for next in range(3):
                if next == pre:
                    continue
                path.append(next)
                yield from bt(index + 1, next, path)
                path.pop()

        iter = bt(0, -1, [])
        return next(islice(iter, k, None), '')


print(Solution().solve(n=2, k=0))
print(Solution().solve(n=2, k=1))
print(Solution().solve(n=2, k=2))
print(Solution().solve(n=2, k=3))
print(Solution().solve(n=2, k=4))
print(Solution().solve(n=2, k=5))
print(Solution().solve(n=2, k=6))
print(Solution().solve(n=2, k=7))

