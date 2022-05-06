# 给出一个长<=1e5的数组 数组里的数每个数<=1e5
# 判断是否数组里的每个数
# 都至少存在一个数与它二进制按位与的结果为0


# 1. trie 可以吗 (不行)
# 一个序列i位是0 的时候，trie树0 1都是合法的
# 异或只需搜trie的一条路 但是与运算要搜两条路 所以不行

# 2. 如果数据范围是1e5的话枚举补集的子集就好了
# 如果之前算过子集就跳过
# 求子集的过程类似记忆化
# 不会这个，你也可以用枚举子集那个做，但是复杂度会高变成 3^n 。
# SoS dp 是 n 2^n。这个会把 每一个 s 对应的所有满足 k&s = k 的数字 k 的个数算出来。记它是 dp[s]。
# 之后枚举每一个数字 a 你其实就是看 dp[mask^a] 有没有值。

# SOS DP
# https://blog.csdn.net/weixin_38686780/article/details/100109753


from collections import defaultdict
from math import ceil, log2
from typing import List

N = ceil(log2(1e5))
UPPER = 1 << N


def queryPairwiseAnd1(nums: List[int]) -> bool:
    """判断是否数组里的每个数都至少存在一个数与它二进制按位与的结果为0
    
    n<=1e5,nums[i]<=1e5
    sosdp 计算每个状态的高维前缀和即可
    """
    # sosdp = [0] * UPPER
    sosdp = defaultdict(int)
    for num in nums:
        sosdp[num] += 1

    for i in range(N):
        for state in range(UPPER):
            if (state >> i) & 1:
                sosdp[state] += sosdp[state ^ (1 << i)]

    for num in nums:
        comp = (UPPER - 1) ^ num
        if sosdp[comp] == 0:
            return False
    return True


def bruteForce(nums: List[int]) -> bool:
    for i in range(len(nums)):
        if all(nums[i] & nums[j] != 0 for j in range(len(nums)) if i != j):
            return False
    return True


if __name__ == '__main__':
    cands = [[1, 2, 3, 4], [23, 9, 12], [123, 77, 2], [223, 99]]
    for nums in cands:
        assert queryPairwiseAnd1(nums) == bruteForce(nums)

