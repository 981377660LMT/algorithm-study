# yield与返回bool的dfs
# !搜索一条路径的问题,dfs可以采用yield的方式,也可以采用返回bool的方式


# https://leetcode.cn/problems/lexicographically-smallest-beautiful-string/
# 2663. 字典序最小的美丽字符串


class Solution:
    def smallestBeautifulString1(self, s: str, k: int) -> str:
        """生成器dfs返回路径."""

        def dfs(pos: int, isLimit: bool, pre1: int, pre2: int):
            if pos == n:
                if not isLimit:
                    yield path
                return
            lower = ords[pos] if isLimit else 97
            for cur in range(lower, 97 + k):
                if cur == pre1 or cur == pre2:
                    continue
                path.append(cur)
                yield from dfs(pos + 1, (isLimit and cur == lower), cur, pre1)
                path.pop()

        path = []
        n = len(s)
        ords = list(map(ord, s))
        res = next(dfs(0, True, 0, -1), [])
        return "".join([chr(v) for v in res])

    def smallestBeautifulString2(self, s: str, k: int) -> str:
        """
        !返回bool的dfs返回路径.
        """

        def dfs(pos: int, isLimit: bool, pre1: int, pre2: int) -> bool:
            if pos == n:
                if not isLimit:
                    return True
                return False
            lower = ords[pos] if isLimit else 97
            for cur in range(lower, 97 + k):
                if cur == pre1 or cur == pre2:
                    continue
                path.append(cur)
                if dfs(pos + 1, (isLimit and cur == lower), cur, pre1):
                    return True  # !找到一条路径就返回
                path.pop()
            return False

        path = []
        n = len(s)
        ords = list(map(ord, s))
        dfs(0, True, 0, -1)
        return "".join([chr(v) for v in path])
