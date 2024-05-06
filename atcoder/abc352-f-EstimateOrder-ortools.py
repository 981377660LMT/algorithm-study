# F - EstimateOrder
# https://atcoder.jp/contests/abc352/tasks/abc352_f
# n个人，排名唯一，给定关于这些人的m条排名信息，问每个人的名次是否能唯一确定。
# 每条排名消息形如： a - b = c，表示a比b名次高c名。
# n<=16.

# ortools、z3、pulp 调库解法
# https://atcoder.jp/contests/abc352/submissions/53115584 ortools 619 ms
# https://atcoder.jp/contests/abc352/submissions/53125345 z3 1847 ms
# https://atcoder.jp/contests/abc352/submissions/53130190 pulp 844 ms
#
# https://blog.51cto.com/u_11866025/5833945
# https://developers.google.cn/optimization?hl=zh-cn


import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")


from copy import deepcopy


if __name__ == "__main__":
    N, M = map(int, input().split())

    model = cp_model.CpModel()

    xs = [model.NewIntVar(1, N, f"x_{i}") for i in range(N)]

    for _ in range(M):
        a, b, c = map(int, input().split())
        a -= 1
        b -= 1
        model.AddLinearConstraint(xs[a] - xs[b], c, c)

    model.AddAllDifferent(xs)

    solver = cp_model.CpSolver()

    ret = []
    for i in range(N):
        feasible_ranks: list[int] = []

        for r in range(1, N + 1):
            tmp_model = deepcopy(model)
            tmp_model.AddLinearConstraint(xs[i], r, r)

            status = solver.Solve(tmp_model)
            if status == cp_model.OPTIMAL:
                feasible_ranks.append(r)

            if len(feasible_ranks) > 1:
                break

        ret.append(feasible_ranks[0] if len(feasible_ranks) == 1 else -1)

    print(*ret)
