# 工作量约束的任务分配
# https://blog.51cto.com/u_11866025/5833945
import pandas as pd
import numpy as np
from ortools.linear_solver import pywraplp

# !每个工人处理的成本表如下：
costs = np.array(
    [
        [90, 76, 75, 70, 50, 74, 12, 68],
        [35, 85, 55, 65, 48, 101, 70, 83],
        [125, 95, 90, 105, 59, 120, 36, 73],
        [45, 110, 95, 115, 104, 83, 37, 71],
        [60, 105, 80, 75, 59, 62, 93, 88],
        [45, 65, 110, 95, 47, 31, 81, 34],
        [38, 51, 107, 41, 69, 99, 115, 48],
        [47, 85, 57, 71, 92, 77, 109, 36],
        [39, 63, 97, 49, 118, 56, 92, 61],
        [47, 101, 71, 60, 88, 109, 52, 90],
    ]
)
num_workers = len(costs)
num_tasks = len(costs[0])

# !每个任务的工作量大小如下：
task_sizes = [10, 7, 3, 12, 15, 4, 11, 5]

# !每个工人能够工作的最大工作量为15。
total_size_max = 15

# 创建一个mip求解器
solver = pywraplp.Solver.CreateSolver("SCIP")
mask = np.array(
    [[solver.BoolVar(f"x_{i}_{j}") for j in range(num_tasks)] for i in range(num_workers)]
)

# 每个工人的工作量不得超过限制
task_sizes_mask = task_sizes * mask
for num in task_sizes_mask.sum(axis=1):
    solver.Add(num <= total_size_max)
# 每个任务必须分配给一个工人
for num in mask.sum(axis=0):
    solver.Add(num == 1)
# 目标函数：总成本最小
costs_mask = costs * mask
solver.Minimize(costs_mask.sum())

# 求解结果
status = solver.Solve()
status = {solver.OPTIMAL: "最优解", solver.FEASIBLE: "可行解"}.get(status)
if status is not None:
    print(status)
    print("最低总成本：", solver.Objective().Value())
    print("分配方案：")
    result = pd.DataFrame(
        costs_mask,
        index=[f"Worker {i}" for i in range(10)],
        columns=[f"Task {i}" for i in range(8)],
    )
    result = result.applymap(lambda x: x.solution_value()).astype("int16")
    print(result)
else:
    print("无解")
