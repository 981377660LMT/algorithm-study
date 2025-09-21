# abc424-g-pulp优化(Hall定理)
# https://atcoder.jp/contests/abc424/tasks/abc424_g
#
# G - 歌曲列表
#
# 问题描述
# 有一个由 N 名偶像（偶像 1, 偶像 2, ..., 偶像 N）组成的团体。
# 另外，有 M 首候选歌曲，编号分别为歌曲 1, 歌曲 2, ..., 歌曲 M。
# 高桥君希望从这些歌曲中选择若干首（0 首也可以）来举办一场演唱会。 但是，同一首歌曲最多只能选择一次。
# 偶像 i (1≤i≤N) 最多可以跳 A_i 首歌。 在演唱会中，歌曲 j (1≤j≤M) 需要 B_j 名偶像来跳舞，并且这首歌的“气氛值”为 C_j（与谁跳舞无关）。
# 高桥君需要选择演唱会要表演的歌曲以及每首歌由哪些偶像来跳舞，同时不能超过每个偶像可跳舞的次数限制。请找出所选歌曲的气氛值总和可能的最大值。
# 约束条件
#
# 1 ≤ N ≤ 100
# 1 ≤ M ≤ 100
# 0 ≤ A_i ≤ M
# 0 ≤ B_j ≤ N
# 0 ≤ C_j ≤ 10^9
# 所有输入均为整数
#
# https://atcoder.jp/contests/abc424/submissions/69478040

from pulp import (
    LpProblem,
    LpVariable,
    LpBinary,
    LpMaximize,
    lpSum,
    PULP_CBC_CMD,
)


N, M = map(int, input().split())
A = list(map(int, input().split()))  # 偶像 i 最多能跳的歌曲数
A.sort(reverse=True)
B = []
C = []
for _ in range(M):
    b, c = map(int, input().split())
    B.append(b)
    C.append(c)

prob = LpProblem(sense=LpMaximize)

vars = [
    LpVariable(f"x_{i}", lowBound=0, upBound=1, cat=LpBinary) for i in range(M)
]  # 是否选择歌曲i

prob.setObjective(lpSum(C[j] * vars[j] for j in range(M)))

# 原始问题的对偶形式的松弛，确保了对于任意数量的偶像子集，其总跳舞能力都足以支撑分配给他们的跳舞任务
for thres in range(N):
    vs = [vars[i] * (B[i] - thres) for i in range(M) if B[i] > thres]
    if vs:
        prob.addConstraint(lpSum(vs) <= sum(A[thres:]))

prob.solve(PULP_CBC_CMD(msg=False))  # CBC求解器
print(round(prob.objective.value()))
