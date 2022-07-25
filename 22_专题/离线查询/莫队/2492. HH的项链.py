"""静态查询区间数字种数 n<=5e4"""

from collections import defaultdict
import sys

input = sys.stdin.readline
n = int(input())
nums = list(map(int, input().split()))
q = int(input())
chunkSize = max(int(n // (q**0.5)), 1)

queries = []
for i in range(q):
    left, right = map(int, input().split())
    left, right = left - 1, right - 1
    queries.append((i, left, right + 1))  # !注意这里 right + 1

# !以询问左端点所在的分块的序号为第一关键字，右端点的大小为第二关键字进行排序
queries.sort(key=lambda x: (x[1] // chunkSize, x[2]))

counter = defaultdict(int)
left, right, count = 0, 0, 0


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


res = [0] * q
for qi, qLeft, qRight in queries:
    # !窗口扩张
    while left > qLeft:
        left -= 1
        add(nums[left])
    while right < qRight:
        add(nums[right])
        right += 1

    # !窗口收缩
    while left < qLeft:
        remove(nums[left])
        left += 1
    while right > qRight:
        right -= 1
        remove(nums[right])

    res[qi] = count

for num in res:
    print(num)
