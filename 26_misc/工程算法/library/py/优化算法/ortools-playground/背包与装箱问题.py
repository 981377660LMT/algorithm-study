# 将一批物品打包到5个箱子中，每个箱子的最大容量为100，求总打包的最大值。

from ortools.linear_solver import pywraplp
import numpy as np


def binPacking() -> None:
    weights = np.array([48, 30, 42, 36, 36, 48, 42, 42, 36, 24, 30, 30, 42, 36, 36])
    values = np.array([10, 30, 25, 50, 35, 30, 15, 40, 30, 35, 45, 10, 20, 30, 25])
    bin_capacities = [100, 100, 100, 100, 100]

    solver = pywraplp.Solver.CreateSolver("SCIP")
    x = np.array(
        [
            [solver.BoolVar(f"x_{b}_{i}") for i in range(len(values))]
            for b in range(len(bin_capacities))
        ]
    )
    # 每个物品最多只能放在一个箱子里
    for item in x.sum(axis=0):
        solver.Add(item <= 1)
    # 每个箱子包装的总重量不超过它的容量
    for pack, bin_capacitie in zip((x * weights).sum(axis=1), bin_capacities):
        solver.Add(pack <= bin_capacitie)
    # 打包项的总价值
    solver.Maximize((x * values).sum())
    status = solver.Solve()
    result = np.frompyfunc(lambda x: x.solution_value(), 1, 1)(x).astype(bool)
    print("总价值：", solver.Objective().Value())
    print("总总量：", (x * weights).sum().solution_value())
    index = np.arange(values.shape[0])
    for i, row in enumerate(result):
        print(f"\n第{i}个箱子")
        print("选中的球：", index[row])
        print("选中的球的重量：", weights[row], "，总重量：", weights[row].sum())
        print("选中的球的价值：", values[row], "，总价值：", values[row].sum())


if __name__ == "__main__":
    binPacking()
