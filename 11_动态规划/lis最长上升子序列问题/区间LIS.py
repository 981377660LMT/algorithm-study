# Static Range LIS Query
# 区间LIS查询 O(n^2logn)
# TODO
# https://judge.yosupo.jp/problem/static_range_lis_query


from bisect import bisect_right


n, q = map(int, input().split())
nums = list(map(int, input().split()))


dp = [[0] * n for _ in range(n)]
for i in range(n):
    lis = []
    for j in range(i, n):
        pos = bisect_right(lis, nums[j])  # 非严格单增
        if pos == len(lis):
            lis.append(nums[j])
        else:
            lis[pos] = nums[j]
        dp[i][j] = len(lis)

for _ in range(q):
    l, r = map(int, input().split())
    print(dp[l][r - 1])
