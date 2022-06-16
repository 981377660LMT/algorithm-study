from typing import List
from itertools import accumulate

# 1 <= nums.length <= 105
# nums[i] 要么是 0 ，要么是 1 。
# https://leetcode-cn.com/problems/minimum-adjacent-swaps-for-k-consecutive-ones/solution/de-dao-lian-xu-k-ge-1-de-zui-shao-xiang-fqbhp/

# 答案转换为： 需要把连续的k个ai, 交换到一起的最小交换次数
# 假设交换前k个1的序号为[o0,o1,...,ok-1]，交换后序号为[x,x+1,...,x+k-1]
# 那么要求的就是|o0-x|+|o1-(x+1)|+...+|ok-1-(x+k-1)|，变形得
# |o0-x|+|(o1-1)-x|+|(ok-1-(k-1))-x| 即对k个原来的1 需要找到一个x使其最小 这个x就是他们的中位数mid


# mid左边的和： a[mid] - a[l] + a[mid] - a[l+1] + ... + a[mid] - a[mid - 1]
# 	             = a[mid] * (mid - l) - (a[l] + a[l + 1] + ... + a[mid - 1])
# 				 = a[mid] * (mid - l) - (sum[mid - 1] - sum[l - 1])  前缀和

# mid右边的和： a[r] - a[mid] + a[r - 1] - a[mid] + ... + a[mid + 1] - a[mid]
# 				 = sum[r] - sum[mid]  - (r - mid ) * a[mid]

# 每次求连续k个ai的和的最小值, 需要快速的求`每个数与中位数的差的绝对值的和`
# 需要预处理下前缀和


class Solution:
    def minMoves(self, nums: List[int], k: int) -> int:
        ones = []
        for i in range(len(nums)):
            if nums[i] == 1:
                ones.append(i - len(ones))
        preSum = [0] + list(accumulate(ones))

        res = int(1e20)
        #    left+k-1<len(ones)
        for left in range(len(ones) + 1 - k):
            right = left + k - 1
            mid = (left + right) // 2
            # mid左右两边的和
            leftSum = ones[mid] * (mid - left) - (preSum[mid] - preSum[left])
            rightSum = preSum[right + 1] - preSum[mid + 1] - (right - mid) * ones[mid]
            res = min(res, leftSum + rightSum)

        return res


print(Solution().minMoves(nums=[1, 0, 0, 1, 0, 1], k=2))
# 输出：1
# 解释：在第一次操作时，nums 可以变成 [1,0,0,0,1,1] 得到连续两个 1 。

