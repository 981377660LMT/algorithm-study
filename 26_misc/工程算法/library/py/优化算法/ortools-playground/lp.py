from ortools.linear_solver import pywraplp


def demo() -> None:
    # 创建默认的GLOP求解器
    solver = pywraplp.Solver.CreateSolver("GLOP")
    # 搜索时间限制为15秒
    solver.set_time_limit(15)
    x = solver.NumVar(0, solver.infinity(), "x")
    y = solver.NumVar(0, solver.infinity(), "y")
    # 定义约束条件
    solver.Add(x + 2 * y <= 14.0)
    solver.Add(3 * x - y >= 0.0)
    solver.Add(x - y <= 2.0)
    print("变量数量：", solver.NumVariables())
    print("约束数量：", solver.NumConstraints())
    # 定义目标函数
    solver.Maximize(3 * x + 4 * y)
    # 调用求解器
    status = solver.Solve()
    if status == pywraplp.Solver.OPTIMAL:
        print(
            f"x={x.solution_value():.2f}, "
            f"y={y.solution_value():.2f}, "
            f"max(3x+4y)={solver.Objective().Value():.2f}"
        )
    else:
        print("无最优解")


# 某县新建一家医院，根据各个科室要求需要配备护士，周一到周日分别最小需要34、25、36、30、28、31、32人.
# 按照规定，一个护士一周要连续上班5天。这家医院至少需要多少个护士？
def 医院每天护士人数分配() -> None:
    min_nums = [34, 25, 36, 30, 28, 31, 32]

    solver = pywraplp.Solver.CreateSolver("SCIP")

    x = [solver.IntVar(0, 18, f"x{i}") for i in range(1, 8)]  # 可以预估每天所需的护士数最大为36/2=18
    solver.Minimize(sum(x))
    for i in range(len(x)):
        solver.Add(sum(x) - x[(i + 1) % 7] - x[(i + 2) % 7] >= min_nums[i])
    # 调用求解器
    status = solver.Solve()
    if status == solver.OPTIMAL:
        print("有最优解")
    elif status == solver.FEASIBLE:
        print("有可行解")
    else:
        print("无最优解")
    print(
        f"最少护士人数 z={solver.Objective().Value():.0f}",
    )
    print("周1到周日开始上班的护士人数分别为：", [x[i].solution_value() for i in range(7)])
    print(
        "周一到周日上班人数分别为：",
        [(sum(x) - x[(i + 1) % 7] - x[(i + 2) % 7]).solution_value() for i in range(7)],
    )


def 数据包络分析() -> None:
    import numpy as np

    # 创建GLOP求解器
    solver = pywraplp.Solver.CreateSolver("GLOP")
    # 搜索时间限制为15秒
    solver.set_time_limit(15)

    data = [
        [20, 149, 1300, 636, 1570],
        [18, 152, 1500, 737, 1730],
        [23, 140, 1500, 659, 1320],
        [22, 142, 1500, 635, 1420],
        [22, 129, 1200, 626, 1660],
        [25, 142, 1600, 775, 1590],
    ]

    # 定义6个变量
    x = [solver.NumVar(0, 1, f"x{i}") for i in range(1, 7)]
    # 定义期望E
    e = solver.NumVar(0, 1, "e")
    solver.Minimize(e)
    # 办事处1的数据
    office1 = data[0]
    # 各办事处的加权平均值
    office_wavg = np.sum(data * np.array(x)[:, None], axis=0)
    # 权重之和为1
    solver.Add(sum(x) == 1)
    # 投入更少
    for i in range(2):
        solver.Add(office_wavg[i] <= office1[i] * e)
    # 产出更多
    for i in range(2, len(data[0])):
        solver.Add(office_wavg[i] >= office1[i])
    # 调用求解器
    status = solver.Solve()
    if status == solver.OPTIMAL:
        print("有最优解")
    elif status == solver.FEASIBLE:
        print("有可行解")
    else:
        print("无最优解")
    print(
        f"目标函数的最小值z={solver.Objective().Value()}，此时目标函数的决策变量为:",
        {v.name(): v.solution_value() for v in solver.variables()},
    )
    print("组合后的投入和产出：", [f"{office_wavg[i].solution_value():.2f}" for i in range(len(data[0]))])


if __name__ == "__main__":
    demo()
    print("---")
    医院每天护士人数分配()
    print("---")
    数据包络分析()
