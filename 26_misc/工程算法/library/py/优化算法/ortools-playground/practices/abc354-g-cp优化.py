from ortools.sat.python import cp_model

n = int(input())
words = [input() for _ in range(n)]
weights = list(map(int, input().split()))

model = cp_model.CpModel()
xs = [model.NewBoolVar(f"x{i}") for i in range(n)]  # 最大独立集，每个点是否选择

for i, si in enumerate(words):
    for j, sj in enumerate(words):
        if len(si) > len(sj) or i == j:
            continue

        ni, nj = len(si), len(sj)

        for s in range(nj - ni + 1):
            if sj[s : s + ni] == si:
                model.AddAtMostOne([xs[i], xs[j]])

model.Maximize(sum([xs[i] * weights[i] for i in range(n)]))

solver = cp_model.CpSolver()
status = solver.Solve(model)
opt = solver.ObjectiveValue()

print(int(round(opt)))
