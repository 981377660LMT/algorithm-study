# AcWing 2492. HH的项链

# 静态查询区间数字种数

# 三种做法，莫队，离线树状数组，主席树
from collections import defaultdict
import sys

input = sys.stdin.readline
n = int(input())
nums = list(map(int, input().split()))
qCount = int(input())
chunkLen = max(int(n // (qCount**0.5)), 1)

queries = []
for i in range(qCount):
    left, right = map(int, input().split())
    left, right = left - 1, right - 1
    queries.append((i, left, right))

# !以询问左端点所在的分块的序号为第一关键字，右端点的大小为第二关键字进行排序
queries.sort(key=lambda x: (x[1] // chunkLen, x[2]))

counter = defaultdict(int)
left, right, count = 0, 0, 1
counter[nums[0]] = 1


def add(num: int) -> None:
    global count
    counter[num] += 1
    if counter[num] == 1:
        count += 1


def remove(num: int) -> None:
    global count
    counter[num] -= 1
    if counter[num] == 0:
        count -= 1


res = [0] * qCount
for qi, qLeft, qRight in queries:
    while left < qLeft:
        remove(nums[left])
        left += 1
    while left > qLeft:
        add(nums[left - 1])
        left -= 1
    while right < qRight:
        add(nums[right + 1])
        right += 1
    while right > qRight:
        remove(nums[right])
        right -= 1
    res[qi] = count

for num in res:
    print(num)
