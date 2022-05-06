# 与运算为0的二元组个数
# nums[i]<=1e6
# n<=1e6

# 对于a[i]，与他匹配的a[j]一定都贡献进a[i]的补集，那么答案就是sigma(f[补ai])


from math import ceil, log2
from typing import List

N = ceil(log2(1e5))
UPPER = 1 << N


def count(nums: List[int]) -> int:
    """与运算为0的二元组个数
    
    n<=1e5,nums[i]<=1e5

    sosdp 计算每个状态的高维前缀和即可
    对于a[i]，与他匹配的a[j]一定都贡献进a[i]的补集，那么答案就是sum(f[补ai])
    """
    sosdp = [0] * UPPER
    for num in nums:
        sosdp[num] += 1

    for i in range(N):
        for state in range(UPPER):
            if (state >> i) & 1:
                sosdp[state] += sosdp[state ^ (1 << i)]

    res = 0
    for num in nums:
        comp = (UPPER - 1) ^ num
        res += sosdp[comp]
    return res


print(count([1, 8]))

