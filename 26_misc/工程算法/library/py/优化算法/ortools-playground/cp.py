# CP 基于可行性（找到可行的解决方案）而非优化（找出最佳解决方案），
# 并且侧重于约束和变量，而非目标函数。
# !事实上，CP 问题可能甚至没有目标函数
# 目标是通过为问题添加约束条件，将大量可能的解决方案缩小为更易于管理的子集。

from typing import List
from ortools.sat.python import cp_model


def demo() -> None:
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

    # 判断是否存在有效的解
    status = solver.Solve(model)
    if status == cp_model.OPTIMAL or status == cp_model.FEASIBLE:  # type: ignore
        print("x =", solver.Value(x), ", y =", solver.Value(y), ", z =", solver.Value(z))
    else:
        print("未找到结果")


# 二元一次方程
def solveEquation() -> None:
    model = cp_model.CpModel()

    x = model.NewIntVar(cp_model.INT32_MIN, cp_model.INT32_MAX, "x")
    y = model.NewIntVar(cp_model.INT32_MIN, cp_model.INT32_MAX, "y")
    model.Add(x - y == 3)
    model.Add(3 * x - 8 * y == 4)
    solver = cp_model.CpSolver()
    status = solver.Solve(model)
    if status == cp_model.OPTIMAL or status == cp_model.FEASIBLE:  # type: ignore
        print("x =", solver.Value(x), ", y =", solver.Value(y))


# 线性多项式约束
def solveLinear() -> None:
    model = cp_model.CpModel()

    x = model.NewIntVar(2, 7, "x")
    y = model.NewIntVar(0, 3, "y")
    model.Add(x + 2 * y == 7)
    solver = cp_model.CpSolver()
    status = solver.Solve(model)
    if status == cp_model.OPTIMAL or status == cp_model.FEASIBLE:  # type: ignore
        print("x =", solver.Value(x), ", y =", solver.Value(y))


# 八皇后问题
def solveEightQueens() -> None:
    from ortools.sat.python import cp_model

    class NQueenSolutionPrinter(cp_model.CpSolverSolutionCallback):
        def __init__(self, queens):
            cp_model.CpSolverSolutionCallback.__init__(self)
            self.__queens = queens
            self.__solution_count = 0

        def solution_count(self):
            return self.__solution_count

        def on_solution_callback(self):
            self.__solution_count += 1
            print(f" --------- {self.solution_count()} ---------- ")

            for queen in self.__queens:
                for i in range(8):
                    if i == self.Value(queen):
                        print(end="Q  ")
                    else:
                        print(end="*  ")
                print()

    model = cp_model.CpModel()

    # 每个皇后必须在不同的行中，记录每行对应的皇后对应的列位置
    queens = [model.NewIntVar(0, 7, f"Q_{i}") for i in range(8)]
    # 每列最多一个皇后
    model.AddAllDifferent(queens)
    # 对角线约束
    model.AddAllDifferent([queens[i] + i for i in range(8)])
    model.AddAllDifferent([queens[i] - i for i in range(8)])

    solver = cp_model.CpSolver()
    solution_printer = NQueenSolutionPrinter(queens)
    solver.parameters.enumerate_all_solutions = True
    solver.Solve(model, solution_callback=solution_printer)

    statistics = f"""
    Statistics
        conflicts : {solver.NumConflicts()}
        branches  : {solver.NumBranches()}
        耗时      : {solver.WallTime():.4f} s
        共找到    : {solution_printer.solution_count()}个解"""
    print(statistics)


def solveSudoku(board: List[List[int]]) -> List[List[int]]:
    model = cp_model.CpModel()

    # 9x9 整数变量矩阵,数值范围1-9
    X = [[model.NewIntVar(1, 9, "x") for _ in range(9)] for _ in range(9)]
    # 每行每个数字最多出现一次
    for i in range(9):
        model.AddAllDifferent(X[i])
    # 每列每个数字最多出现一次
    for j in range(9):
        model.AddAllDifferent([X[i][j] for i in range(9)])
    # 每个 3x3 方格每个数字最多出现一次
    for i0 in range(3):
        for j0 in range(3):
            constraint = [X[3 * i0 + i][3 * j0 + j] for i in range(3) for j in range(3)]
            model.AddAllDifferent(constraint)
    for i in range(9):
        for j in range(9):
            if board[i][j] != 0:
                model.Add(X[i][j] == board[i][j])

    solver = cp_model.CpSolver()
    status = solver.Solve(model)
    if status != cp_model.OPTIMAL and status != cp_model.FEASIBLE:  # type: ignore
        print("无正确的解")
        return []
    r = [[solver.Value(X[i][j]) for j in range(9)] for i in range(9)]
    return r


if __name__ == "__main__":
    solveEquation()
    solveLinear()
    solveEightQueens()
    print(
        solveSudoku(
            [
                [5, 3, 0, 0, 7, 0, 0, 0, 0],
                [6, 0, 0, 1, 9, 5, 0, 0, 0],
                [0, 9, 8, 0, 0, 0, 0, 6, 0],
                [8, 0, 0, 0, 6, 0, 0, 0, 3],
                [4, 0, 0, 8, 0, 3, 0, 0, 1],
                [7, 0, 0, 0, 2, 0, 0, 0, 6],
                [0, 6, 0, 0, 0, 0, 2, 8, 0],
                [0, 0, 0, 4, 1, 9, 0, 0, 5],
                [0, 0, 0, 0, 8, 0, 0, 7, 9],
            ]
        )
    )
