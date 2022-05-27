# https://codeforces.com/problemset/problem/1310/A
# 给你两个数组 a 和 t，长度相同且不超过 2e5，1<=a[i]<=1e9, 1<=t[i]<=1e5。
# 每次操作，你可以给某个 a[i]+=1，花费为 t[i]。你可以操作任意多次。
# 请问要使 a 中所有数字均不相同，最小花费是多少？
# 输入
# a=[3,7,9,7,8]
# t=[5,2,5,7,5]
# 输出 6
# 解释 把第二个数字加三次，得到 [3,10,9,7,8]，花费为 2*3=6

# https://www.luogu.com.cn/blog/endlesscheng/solution-cf1310a

# 按 a[i]a[i] 从小到大遍历，把 a[i]a[i] 相同的元素分为一组。

# 比如 a=[1,1,1,2,9,9]a=[1,1,1,2,9,9]，那么一开始把等于 1 的元素的花费丢到一个大根堆中，其中有一个 1 是不需要增大的，
# 把最大的花费弹出，表示对应的那个 1 不需要增大。
# 然后把等于 2 的元素的花费丢到大根堆中，同样把最大的花费弹出，表示对应的那个元素不需要继续增大。
# 然后把等于 3 的元素的花费丢到大根堆中，由于此时没有等于 3 的元素，直接从堆中弹出一个最大的花费。

# 模拟上述过程，同时用一个变量 costSum 表示堆中的花费之和，每次弹出最大花费后，累加 costSum，即为答案。
from heapq import heappop, heappush


nums = list(map(int, input().split()))
costs = list(map(int, input().split()))
pairs = sorted(zip(nums, costs), key=lambda x: x[0])
n = len(nums)


pq = []
res = 0
# 每次curVal增大,所有等于 cur 的元素的花费入堆，然后出去一个最大的花费，再累加curSum
curVal = 0
runningSum = 0
index = 0

while pq or index < n:  # 没看完所有的数
    if not pq:
        curVal = pairs[index][0]
    while index < n and pairs[index][0] == curVal:
        heappush(pq, -pairs[index][1])
        runningSum += pairs[index][1]
        index += 1
    runningSum -= -heappop(pq)
    res += runningSum
    curVal += 1


print(res)
