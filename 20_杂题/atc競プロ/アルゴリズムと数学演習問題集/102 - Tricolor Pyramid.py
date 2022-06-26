# 三种颜色的金字塔
# 如果下面两块颜色相同 则上面那块颜色也相同
# 如果下面两块颜色不同 则上面那块颜色为另一种
# !求金字塔顶端的颜色
# n<=4e5

# !上のブロックの整数は -(p1+p2) mod3
# !类似二项式定理的系数推导

MAPPING = {'B': 0, 'W': 1, 'R': 2}
RMAPPING = {0: 'B', 1: 'W', 2: 'R'}
n = int(input())
colors = list(map(MAPPING.get, input().split()))
