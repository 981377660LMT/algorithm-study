# 某公司雇有 N 名员工，每名员工可以负责多个项目，但一个项目只能交由一名员工负责。现在该公司接到 M 个项目，令 A(i,j)A(i,j) 表示第 i 名员工负责第 j
# 个项目所带来的收益，那么如果项目分配得当，总收益最大是多少?
# 1 ≤ N，M ≤ 1000

n, m = [int(i) for i in input().split()]
matrix = []
for _ in range(n):
    row = [int(i) for i in input().split()]
    matrix.append(row)

res = 0
for col in zip(*matrix):
    res += max(col)
print(res)
