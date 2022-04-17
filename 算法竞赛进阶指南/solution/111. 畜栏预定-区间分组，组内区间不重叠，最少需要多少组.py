# 区间分组，组内区间不重叠，最少需要多少组

# 输出格式
# 第 1 行：输入一个整数，代表所需最小畜栏数。
# 第 2..N+1 行：第 i+1 行输入第 i 头牛被安排到的畜栏编号，
# 编号是从 1 开始的 连续 整数，只要方案合法即可。


from heapq import heappop, heappush


n = int(input())

intervals = []
for i in range(n):
    a, b = map(int, input().split())
    intervals.append((a, b, i))
intervals.sort()


pq = []
group = 1
res = [-1] * n

for start, end, index in intervals:
    preGroup = None
    if pq and start > pq[0][0]:
        preGroup = heappop(pq)[1]

    if preGroup is not None:
        heappush(pq, (end, preGroup))
        res[index] = preGroup
    else:
        heappush(pq, (end, group))
        res[index] = group
        group += 1

print(len(pq))
for num in res:
    print(num)
