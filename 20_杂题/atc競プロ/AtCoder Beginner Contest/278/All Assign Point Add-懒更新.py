# !区间整体赋值 单点加
# 三种操作
# ! 1 x 将区间所有数变为x
# ! 2 index x 在index处加上x
# ! 3 index 输出index处的值

from collections import defaultdict


n = int(input())
nums = list(map(int, input().split()))
q = int(input())

add = defaultdict(int)  # ! 每个点的增量
curRange = -1  # ! 当前区间覆盖的值

for _ in range(q):
    op, *args = list(map(int, input().split()))
    if op == 1:
        x = args[0]
        curRange = x
        add.clear()
    elif op == 2:
        i, x = args
        add[i] += x
    else:
        i = args[0]
        if curRange == -1:
            print(nums[i - 1] + add[i])
        else:
            print(curRange + add[i])
