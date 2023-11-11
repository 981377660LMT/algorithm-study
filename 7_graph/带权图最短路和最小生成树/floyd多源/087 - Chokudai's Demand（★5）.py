# 使得 i到j距离<=p 的 (i,j) 对正好k对(i<j)
# 权值等于-1表示权值为x 否则为输入的权值
# 求x的值得个数 如果存在无数个 输出'Infinity'


# n<=40
# p<=1e9
# k>=0
# !正好k对等价于 `<=k对` 减去 `<=k-1对` 的情况
# !无数个:x的值对对数无影响 即去除那些边剩下的对数为k

from itertools import combinations, product
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n, p, k = map(int, input().split())
dists = [[0] * n for _ in range(n)]

for i in range(n):
    nums = list(map(int, input().split()))
    for j, v in enumerate(nums):
        dists[i][j] = v


def countNGT(mid: int) -> int:
    """x取mid时费用,距离不超过p的对数"""
    curDists = [[0] * n for _ in range(n)]
    for i, j in product(range(n), repeat=2):
        curDists[i][j] = mid if dists[i][j] == -1 else dists[i][j]

    for k, i, j in product(range(n), repeat=3):
        curDists[i][j] = min(curDists[i][j], curDists[i][k] + curDists[k][j])

    return sum(curDists[i][j] <= p for i, j in combinations(range(n), 2))


# 与x无影响
if countNGT(int(1e18)) == k:
    print('Infinity')
    exit(0)


def cal(upper: int) -> int:
    # 不超过upper对时,-1代表的边权的最小值
    left, right = 0, p + 1
    while left <= right:
        mid = (left + right) // 2
        cur = countNGT(mid)
        if cur <= upper:
            right = mid - 1
        else:
            left = mid + 1
    return left


res1, res2 = cal(k), cal(k - 1)
print(res2 - res1)

# 有问题 todo
