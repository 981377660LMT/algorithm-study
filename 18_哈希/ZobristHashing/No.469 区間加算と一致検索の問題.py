# 区间加,区间一致查询

# 初始时序列为全 0, 有 q 个操作, 操作有两种类型:
# ! L R K : 将区间 [L,R] 的每个数加上 K
# ? : 查询是否出现过当前数组, 如果出现过, 输出最早出现是在第几个操作后

# 解:
# 需要这样的哈希函数:
# - H(x)+H(y) = H(x+y)
# - `k*H(x)` = `H(k*x)`
# 现在数组的哈希值为h ,区间加后k后变为 h+(hash(R)-hash(L))*k

from collections import defaultdict
from random import randint


pool = defaultdict(lambda: randint(1, (1 << 61) - 1))
curHash = 0
history = {0: 0}

n, q = map(int, input().split())
for i in range(q):
    op, *args = input().split()
    if op == "?":
        print(history.get(curHash, -1))
    else:
        left, right, add = map(int, args)
        curHash += (pool[right] - pool[left]) * add  # 区间加
        history.setdefault(curHash, i + 1)
