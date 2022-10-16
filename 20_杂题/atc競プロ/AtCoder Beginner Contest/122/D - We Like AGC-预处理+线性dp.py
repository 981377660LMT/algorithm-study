# https://blog.hamayanhamayan.com/entry/2019/03/25/232610
# https://blog.csdn.net/flyawayl/article/details/88837889

# 求长度为n的字符串个数
# 并且这些字符串要满足如下两个条件:
# !1. 只包括ACGT四种字符
# !2. 不存在子串"AGC"
# !3. 邻位交换字符也不出现子串"AGC"

# 3<=n<=100
# 如何避免子串出现"AGC"
# !dfs(index,pre1,pre2,pre3)

from collections import defaultdict
from functools import lru_cache
from itertools import product
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


# 预处理转移
mp = defaultdict(list)
for (pre1, pre2, pre3) in product("ACGT#", repeat=3):
    if pre3 == "#" and (pre1 != "#" or pre2 != "#"):
        continue
    if pre2 == "#" and pre1 != "#":
        continue

    for cur in "ACGT":
        s = pre1 + pre2 + pre3 + cur
        if "AGC" in s:
            continue

        # !交换相邻两位字符也不出现子串"AGC"
        ok = True
        sb = [pre1, pre2, pre3, cur]
        for i in range(3):
            sb[i], sb[i + 1] = sb[i + 1], sb[i]
            if "AGC" in "".join(sb):
                ok = False
                break
            sb[i], sb[i + 1] = sb[i + 1], sb[i]
        if ok:
            mp[((pre1, pre2, pre3))].append(cur)

if __name__ == "__main__":
    n = int(input())

    @lru_cache(None)
    def dfs(index: int, pre1: str, pre2: str, pre3: str) -> int:
        if index == n:
            return 1
        res = 0
        for cur in mp[(pre1, pre2, pre3)]:
            res += dfs(index + 1, pre2, pre3, cur)
            res %= MOD
        return res

    res = dfs(0, "#", "#", "#")
    dfs.cache_clear()
    print(res)
