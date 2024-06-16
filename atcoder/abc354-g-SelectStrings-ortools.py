# G - Select Strings(dag上的最大带权独立集)
# https://atcoder.jp/contests/abc354/tasks/abc354_g
# 给定n个字符串,每个字符串有一个权值.
# 从中选出若干字符串,使得选出的字符串互不包含,且权值和最大.
# n<=100.|s|<=5000.


from ortools.sat.python import cp_model


def max2(a: int, b: int) -> int:
    return a if a > b else b


if __name__ == "__main__":
    n = int(input())
    words = [input() for _ in range(n)]
    weights = list(map(int, input().split()))

    mp = dict()
    for i, word in enumerate(words):
        mp[word] = max2(mp.get(word, 0), weights[i])
    words = list(mp.keys())
    weights = list(mp.values())
    n = len(words)

    model = cp_model.CpModel()
    xs = [model.NewBoolVar(f"x{i}") for i in range(n)]
    for i, si in enumerate(words):
        for j, sj in enumerate(words):
            if i != j and si in sj:
                model.AddAtMostOne([xs[i], xs[j]])

    model.Maximize(sum(x * w for x, w in zip(xs, weights)))

    solver = cp_model.CpSolver()
    status = solver.Solve(model)
    if status == cp_model.OPTIMAL:
        print(sum(w for x, w in zip(xs, weights) if solver.Value(x)))
