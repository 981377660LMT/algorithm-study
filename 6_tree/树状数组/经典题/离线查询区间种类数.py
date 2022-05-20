# 离线查询区间种类数
# 把所有数据读入后排序，从左到右依次处理。
from BIT import BIT1

n = int(input())
nums = list(map(int, input().split()))
m = int(input())
queries = []
for i in range(m):
    left, right = map(int, input().split())
    queries.append((left, right, i))

bit = BIT1(int(5e4 + 10))

