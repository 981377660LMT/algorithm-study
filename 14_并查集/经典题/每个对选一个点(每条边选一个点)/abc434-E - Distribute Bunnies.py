# E - Distribute Bunnies (分发兔子)
# https://atcoder.jp/contests/abc434/tasks/abc434_e
#  数轴上有 $N$ 只兔子，编号为 $1$ 到 $N$。兔子 $i$ 位于坐标 $X_i$。
# 同一个坐标可能有多只兔子。
# 每只兔子都有一个“跳跃力”参数，兔子 $i$ 的跳跃力为 $R_i$。
# 接下来，所有兔子都要跳跃恰好一次。位于坐标 $x$ 且跳跃力为 $r$ 的兔子跳跃后，会移动到坐标 $x+r$ 或 $x-r$。
# 每只兔子可以自由选择跳到哪个坐标。
# 请计算所有兔子跳跃后，兔子所在的坐标种类数（即有多少个不同的坐标上有兔子）的最大可能值。
#
# !对于每只兔子 $i$，我们有两个选择（坐标 $X_i - R_i$ 和 $X_i + R_i$），我们需要从中选一个，使得最终选出的不同坐标数量最大化。
# 节点：所有可能的落点坐标。
# 边：每只兔子连接它的两个可能落点。
# 目标：每条边选一个端点，使得选出的不同端点数最多。

from SelectOneFromEachPair import SelectOneFromEachPairMap


N = int(input())
X, R = [0] * N, [0] * N

coords = set()
for i in range(N):
    X[i], R[i] = list(map(int, input().split()))
    coords.add(X[i] - R[i])
    coords.add(X[i] + R[i])

sortedCoords = sorted(list(coords))
mp = {v: i for i, v in enumerate(sortedCoords)}

uf = SelectOneFromEachPairMap()


for i in range(N):
    u = mp[X[i] - R[i]]
    v = mp[X[i] + R[i]]
    uf.union(u, v)

print(uf.solve())
