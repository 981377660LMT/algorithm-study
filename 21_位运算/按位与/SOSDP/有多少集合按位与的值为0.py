# https://www.luogu.com.cn/problem/solution/CF449D

from math import ceil, log2
from typing import List

N = ceil(log2(1e6))
UPPER = 1 << N
MOD = int(1e9 + 7)

POW2 = [1] * UPPER
for i in range(1, UPPER):
    POW2[i] = (POW2[i - 1] << 1) % MOD


def solve(nums: List[int]) -> int:
    """多少集合按位与为0
    
    n<=1e6,nums[i]<=1e6
    sosdp 计算每个状态的高维前缀和即可
    选的一个集合满足条件，显然当且仅当每一位都有人是 0。于是我们考虑每个数能够满足哪些位（就是取反）。
    正着不好算，考虑反过来去掉不合法的
    记g[i]表示按位与后结果所有位上至少有i个1的方案数，答案为恰好有0个1不好直接算
    根据容斥原理，ans=g[0]-g[1]+g[2]-g[3]+g[4]...
    """
    sosdp = [0] * UPPER
    for i, num in enumerate(nums, start=1):
        sosdp[num] += 1  # 初始化每个数的贡献

    for i in range(N):
        for state in range(UPPER):
            if (state >> i) & 1:
                # 统计超集前缀和  state ^ (1 << i) 和哪些集合与不为0
                sosdp[state ^ (1 << i)] += sosdp[state]

    res = 0
    for subset in range(UPPER):
        if not sosdp[subset]:
            continue
        # 不能全不选
        count = POW2[sosdp[subset]] - 1
        if bin(subset).count('1') & 1:
            res -= count
        else:
            res += count
        res %= MOD

    return res


print(solve([2, 3, 3]))  # 0
print(solve([0, 1, 2, 3]))  # 10
print(solve([5, 2, 0, 5, 2, 1]))  # 53
