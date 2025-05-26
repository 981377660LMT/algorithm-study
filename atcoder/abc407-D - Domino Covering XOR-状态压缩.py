# D - Domino Covering XOR
# !https://atcoder.jp/contests/abc407/editorial/13104
# !在一个网格图上放置不重叠的多米诺骨牌，最大化未覆盖格子的权值异或和.
#
# 有一个格子图，最多 HW≤20 个顶点，每个顶点有一个 0≤A<2⁶⁰ 的权重。
# 我们可以选择若干对相邻顶点放置骨牌（匹配），不能重叠。
# 匹配中未被覆盖的顶点上的权值做按位 XOR ，记为得分。求最大可能得分。
#
# 1. 生成所有可能的多米诺骨牌边的掩码
# 2. dp 处理出所有可能的状态
# 3. 处理子集异或和，求答案

DIR4 = [(-1, 0), (1, 0), (0, -1), (0, 1)]

if __name__ == "__main__":
    H, W = map(int, input().split())
    A = [list(map(int, input().split())) for _ in range(H)]
    N = H * W

    def idx(i: int, j: int):
        return i * W + j

    edges = []
    for x in range(H):
        for y in range(W):
            id1 = idx(x, y)
            if (x + 1) < H:
                id2 = idx(x + 1, y)
                edges.append((1 << id1) | (1 << id2))
            if (y + 1) < W:
                id2 = idx(x, y + 1)
                edges.append((1 << id1) | (1 << id2))

    canReach = [False] * (1 << N)
    canReach[0] = True
    for pre in range(1 << N):
        if not canReach[pre]:
            continue
        for cur in edges:
            if pre & cur == 0:
                canReach[pre | cur] = True

    subsetXor = [0] * (1 << N)
    for x in range(H):
        for y in range(W):
            id1 = idx(x, y)
            for pre in range(1 << N):
                if pre & (1 << id1):
                    subsetXor[pre] ^= A[x][y]

    res = 0
    allXor = subsetXor[-1]
    for s in range(1 << N):
        if canReach[s]:
            res = max(res, allXor ^ subsetXor[s])

    print(res)
