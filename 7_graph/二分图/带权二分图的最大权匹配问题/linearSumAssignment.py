# 任务分配/任务调度
# 分配问题(assignment problemlinearSumAssignment/)

# n个人分配n项任务，一个人只能分配一项任务，一项任务只能分配给一个人，
# 将一项任务分配给一个人是需要支付报酬，如何分配任务，保证支付的报酬总数最小。
# 简单的说：就是n*n矩阵中，选取n个元素，每行每列各有一个元素，使得和最小。
# 1<=n<=500 -1e9<=Aij<=1e9
# https://judge.yosupo.jp/problem/assignment


from scipy.optimize import linear_sum_assignment

INF = int(1e18)


if __name__ == "__main__":
    n = int(input())
    cost_matrix = [list(map(int, input().split())) for _ in range(n)]
    row_ind, col_ind = linear_sum_assignment(cost_matrix, maximize=False)
    cost = sum(cost_matrix[r][c] for r, c in zip(row_ind, col_ind))

    print(cost)
    print(*col_ind)
    print(row_ind, col_ind)
