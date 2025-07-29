# G - Longest Chord Chain
# https://atcoder.jp/contests/abc410/tasks/abc410_g
# 在圆周上有 2N 个点标号 1,2,…,2N（顺时针）。给定 N 条不重合端点的弦（匹配），第 i 条弦连接 Ai 和 Bi。你需要执行以下两步操作各一次：
#
# 从这 N 条弦中保留任意个两两不相交的子集，删除其他。
# 再在圆上任意加一条新弦（端点可任意选，不必是整数点）。
# 此时所有弦（原保留的 + 新加的）两两相交的点的个数最大能达到多少？
# 等价地：保留一组无交叉弦 C 后，新弦能够交叉 C 中恰好那些“跨过某个割点”的弦。
# 对于任何一组无交叉弦，其最大可能被一条新弦交叉的数量，恰即这组弦的“最大嵌套深度”——也就是这组弦按含嵌关系能形成的最长链的长度。
# !因此题目简化为：在原 N 条弦（匹配）里，找出一条最大长度的“嵌套链”——即存在弦 (c1,…,ck)，可在某个起点处线性化后满足
# L1 < L2 < … < Lk < Rk < … < R2 < R1
# 的最长 k。


from bisect import bisect_left


if __name__ == "__main__":
    N = int(input())
    A, B = [0] * N, [0] * N
    for i in range(N):
        A[i], B[i] = map(int, input().split())

    intervals = []
    n2 = 2 * N
    for a, b in zip(A, B):
        if a > b:
            a, b = b, a
        intervals.append((a, b))
        intervals.append((b, a + n2))

    intervals.sort(key=lambda x: (x[0], -x[1]))

    dp = []
    for _, r in intervals:
        r = -r
        pos = bisect_left(dp, r)
        if pos == len(dp):
            dp.append(r)
        else:
            dp[pos] = r

    print(len(dp))
