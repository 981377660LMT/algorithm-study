# 最小字典序拼接字符串
# !n个单词中选择k个进行排列 求能拼接出的字典序最小的字符串(拼接最小序)
# !k<=n<=50
# len(si)<=50
# 应该是一个50^4的dp
# !倒着dp 记忆化搜索

from functools import cmp_to_key, lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

BIG = chr(ord("z") + 1)


def mergeSort(s1: str, s2: str) -> int:
    """拼接两个字符串，字典序最小"""
    return -1 if s1 + s2 < s2 + s1 else 1


if __name__ == "__main__":
    n, k = map(int, input().split())
    words = [input() for _ in range(n)]
    words.sort(key=cmp_to_key(lambda s1, s2: mergeSort(s1, s2)))  # 拼接最小序

    @lru_cache(None)
    def dfs(index: int, count: int) -> str:
        """前index个单词中选count个单词拼接出的字典序最小的字符串"""
        if index == n:
            return "" if count == k else BIG

        res1 = dfs(index + 1, count)
        res2 = words[index] + dfs(index + 1, count + 1)
        return min(res1, res2)

    res = dfs(0, 0)
    dfs.cache_clear()
    print(res)
