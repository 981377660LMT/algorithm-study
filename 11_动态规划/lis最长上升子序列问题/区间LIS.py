# Static Range LIS Query
# 区间LIS查询
# TODO
# https://judge.yosupo.jp/problem/static_range_lis_query

n, q = map(int, input().split())
perm = list(map(int, input().split()))  # 0-n-1的排列
for _ in range(q):
    left, right = map(int, input().split())
