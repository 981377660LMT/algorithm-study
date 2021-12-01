from typing import List

# 给你一个整数数组 perm ，它是前 n 个正整数的排列，且 n 是个 奇数 。
# encoded[i] = perm[i] ^ perm[i+1]
# 已知encode 求perm

# 总结:
# 1. perm 是前 n 个正整数的排列,我们就可以得到 perm 数组所有元素的异或和
# 2. 求出perm数组第一个数之后，转化为1720题
class Solution:
    def decode(self, encoded: List[int]) -> List[int]:
        n = len(encoded) + 1
        all_xor = 0
        for num in range(1, n + 1):
            all_xor ^= num

        all_xor_without_first = 0
        for i in range(1, n, 2):
            all_xor_without_first ^= encoded[i]

        first = all_xor_without_first ^ all_xor
        res = [first]
        for e in encoded:
            res.append(res[-1] ^ e)
        return res


print(Solution().decode(encoded=[3, 1]))
# 输出：[1,2,3]
# 解释：如果 perm = [1,2,3] ，那么 encoded = [1 XOR 2,2 XOR 3] = [3,1]
