# D - Strange Mirroring
# https://atcoder.jp/contests/abc380/tasks/abc380_d
# 给定一个字符串 s ，重复下述操作 无数次：
# 将 s 的字母大小写反转成 t ，加到 s 后面 给定 q 个询问，每个询问问第 k 个字符是什么。
#
# 1<=k<=1e18
# https://atcoder.jp/contests/abc380/submissions/59834055 (迭代写法)

if __name__ == "__main__":
    S = input()
    Q = int(input())
    K = list(map(int, input().split()))

    def dfs(depth: int, cur: int) -> str:
        """复制depth次后的第cur位是多少"""
        if depth == 0:
            return S[cur - 1]
        length = len(S) * 2**depth
        mid = length // 2
        if cur <= mid:
            return dfs(depth - 1, cur)
        return dfs(depth - 1, cur - mid).swapcase()

    for k in K:
        print(dfs(60, k), end=" ")
