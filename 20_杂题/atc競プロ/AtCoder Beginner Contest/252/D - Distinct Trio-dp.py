"""
n,ai<=2e5
求ai aj ak 不同的三元组数 (i,j,k) 其中 i<j<k

三种的情况还好说,如果是k元组,就需要dp了
https://atcoder.jp/contests/abc252/editorial/4011
"""


from collections import Counter
from functools import lru_cache
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    @lru_cache(None)
    def dfs(value: int, remain: int) -> int:
        """
        value まで確定していて、集合のサイズが size である通り数
        1,2,… の順に 「n を集合に入れるか？」を考えていきましょう
        """
        if remain < 0:
            return 0
        if value == max_ + 1:
            return int(remain == 0)

        res = dfs(value + 1, remain)
        if remain - 1 >= 0 and value in counter:
            res += dfs(value + 1, remain - 1) * counter[value]
        return res

    n = int(input())
    counter = Counter(map(int, input().split()))
    max_ = max(counter)
    res = dfs(1, 3)
    dfs.cache_clear()
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
