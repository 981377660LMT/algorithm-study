# 每次区间[left,right]的值加上delta
# !求每次变化后相邻绝对值差之和


from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n, q = map(int, input().split())
nums = list(map(int, input().split()))
diff = defaultdict(int)  # 差分 每个位置处的值减去前一个位置处的值
for i in range(1, n):
    diff[i] = nums[i] - nums[i - 1]
res = sum(abs(d) for d in diff.values())

for _ in range(q):
    left, right, delta = map(int, input().split())
    left, right = left - 1, right - 1
    pre = abs(diff[left]) + abs(diff[right + 1])  # 两个端点处的值
    if left > 0:
        diff[left] += delta
    if right + 1 < n:
        diff[right + 1] -= delta

    cur = abs(diff[left]) + abs(diff[right + 1])
    res += cur - pre
    print(res)
