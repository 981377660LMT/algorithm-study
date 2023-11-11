# 给定字符串S及其转换定义,问转换t次后第k个字符是什么? 0<=t,1<=k<=1e18
# 转换定义 a=>bc b=>ca c=>ab
# !类似完全二叉树递归 第t,k个字母由第t-1,k>>1个字母转换而来
# 当k变为0时继续计算 三个一循环
# 分治

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

MAPPING = {"A": "BC", "B": "CA", "C": "AB"}


def main() -> None:
    def dfs(t: int, k: int) -> str:
        if t == 0:
            return s[k]
        if k == 0:
            mod_ = t % 3
            res = s[0]
            for _ in range(mod_):
                res = MAPPING[res][0]
            return res
        else:
            pre = dfs(t - 1, k >> 1)
            return MAPPING[pre][k & 1]

    s = input()
    q = int(input())
    for _ in range(q):
        t, k = map(int, input().split())
        # print(["a", "b", "c"][dfs(t, k)])
        print(dfs(t, k - 1))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
