from typing import List

# 将数组分成 3 个非空的部分，使得所有这些部分表示相同的二进制值。
# 如果无法做到，就返回 [-1, -1]。
# 前导零也是被允许的，所以 [0,1,1] 和 [1,1] 表示相同的值。


# 总结：如果要将数组进行三等分，那么每组中1的数目一定要相同
# https://leetcode.com/problems/three-equal-parts/discuss/1343709/2-clean-Python-linear-solutions
class Solution:
    def threeEqualParts(self, arr: List[int]) -> List[int]:
        ones = [i for i, d in enumerate(arr) if d == 1]
        print(ones)
        if not ones:
            return [0, 2]
        elif len(ones) % 3 != 0:
            return [-1, -1]

        # get the start indices of the 3 groups
        i, j, k = ones[0], ones[len(ones) // 3], ones[len(ones) // 3 * 2]
        groupSize = len(arr) - k

        # compare the three groups
        if arr[i : i + groupSize] == arr[j : j + groupSize] == arr[k : k + groupSize]:
            return [i + groupSize - 1, j + groupSize]

        return [-1, -1]


print(Solution().threeEqualParts([1, 0, 1, 0, 1]))
# 输出：[0,3]
