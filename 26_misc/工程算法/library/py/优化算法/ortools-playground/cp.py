from ortools.sat.python import cp_model

# 声明模型
model = cp_model.CpModel()
# 创建变量，同时限制变量范围
x = model.NewIntVar(0, 2, "x")
y = model.NewIntVar(0, 2, "y")
z = model.NewIntVar(0, 2, "z")
# 创建约束
model.Add(x != y)
# 调用规划求解器
solver = cp_model.CpSolver()
# 设置规划求解器的最长计算时间
# solver.parameters.max_time_in_seconds = 10
status = solver.Solve(model)

# 判断是否存在有效的解
if status == cp_model.OPTIMAL or status == cp_model.FEASIBLE:
    print("x =", solver.Value(x), ", y =", solver.Value(y), ", z =", solver.Value(z))
else:
    print("未找到结果")
