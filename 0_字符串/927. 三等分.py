from typing import List, Sequence

# !将数组分成 3 个非空的部分，使得所有这些部分表示相同的二进制值。
# 如果无法做到，就返回 [-1, -1]。
# 前导零也是被允许的，所以 [0,1,1] 和 [1,1] 表示相同的值。


# !如果要将数组进行三等分，那么每组中1的数目一定要相同
# https://leetcode.com/problems/three-equal-parts/discuss/1343709/2-clean-Python-linear-solutions


def useArrayHasher(nums: Sequence[int], mod=10**11 + 7, base=1313131):
    n = len(nums)
    prePow = [1] * (n + 2)
    preHash = [0] * (n + 1)
    for i in range(n):
        prePow[i + 1] = (prePow[i] * base) % mod
        preHash[i] = (preHash[i - 1] * base + nums[i]) % mod

    def sliceHash(left: int, right: int):
        """切片 `nums[left:right]` 的哈希值"""
        right -= 1
        return (preHash[right] - preHash[left - 1] * prePow[right - left + 1]) % mod

    return sliceHash


class Solution:
    def threeEqualParts(self, arr: List[int]) -> List[int]:
        ones = [i for i, d in enumerate(arr) if d == 1]
        n = len(ones)
        if n == 0:
            return [0, 2]
        if n % 3 != 0:
            return [-1, -1]

        # 三个部分的起始位置
        i, j, k = ones[0], ones[n // 3], ones[n // 3 * 2]
        groupSize = len(arr) - k

        if arr[i : i + groupSize] == arr[j : j + groupSize] == arr[k : k + groupSize]:
            return [i + groupSize - 1, j + groupSize]

        # !可以字符串哈希(负)优化
        # hasher = useArrayHasher(arr)
        # if hasher(i, i + groupSize) == hasher(j, j + groupSize) == hasher(k, k + groupSize):
        #     return [i + groupSize - 1, j + groupSize]

        return [-1, -1]


print(Solution().threeEqualParts([1, 0, 1, 0, 1]))
# 输出：[0,3]
print(Solution().threeEqualParts([0, 0, 0, 0, 0]))
