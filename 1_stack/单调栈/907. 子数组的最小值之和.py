from typing import List

MOD = int(1e9 + 7)

# 找到 min(b) 的总和，其中 b 的范围为 arr 的每个（连续）子数组。
# 需要维护区间极值:单调栈
# https://leetcode-cn.com/problems/sum-of-subarray-minimums/solution/python3-tong-84ti-zui-da-zhi-fang-tu-by-5ersw/

# 思路：考虑每个极小值在多少个子数组里产生贡献


class Solution:
    def sumSubarrayMins(self, arr: List[int]) -> int:
        arr.append(-int(1e20))
        stack = [-1]
        res = 0

        for i in range(len(arr)):
            while stack and arr[stack[-1]] > arr[i]:
                j = stack.pop()
                k = stack[-1]
                # 在(stack[-2], i)范围内（exclusive）的最小值都是stack[-1]。
                # i-j j-k 表示开头选择*结尾选择
                res += arr[j] * (i - j) * (j - k)
            stack.append(i)
        return res % MOD


print(Solution().sumSubarrayMins(arr=[3, 1, 2, 4]))
# 解释：
# 子数组为 [3]，[1]，[2]，[4]，[3,1]，[1,2]，[2,4]，[3,1,2]，[1,2,4]，[3,1,2,4]。
# 最小值为 3，1，2，4，1，1，2，1，1，1，和为 17。

