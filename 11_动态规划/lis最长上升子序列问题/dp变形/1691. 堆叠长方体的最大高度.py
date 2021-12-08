from typing import List

# 如果 widthi <= widthj 且 lengthi <= lengthj 且 heighti <= heightj ，
# 你就可以将长方体 i 堆叠在长方体 j 上。你可以通过旋转把长方体的长宽高重新排列，以将它放在另一个长方体上。
# 请你从 cuboids 选出一个 子集 ，并将它们堆叠起来。
# 返回 堆叠长方体 cuboids 可以得到的 最大高度 。

# 1 <= n <= 100
# 面试题 08.13. 堆箱子(LIS)
# https://leetcode-cn.com/problems/maximum-height-by-stacking-cuboids/solution/li-kou-jia-jia-zui-chang-shang-sheng-zi-o4u5b/


class Solution:
    def maxHeight(self, cuboids: List[List[int]]) -> int:
        n = len(cuboids)
        # 前面的箱子放在上面
        cuboids = sorted(map(sorted, cuboids))
        dp = [max(cuboid) for cuboid in cuboids]
        for i in range(n):
            for j in range(i):
                if all(cuboids[j][k] <= cuboids[i][k] for k in range(3)):
                    dp[i] = max(dp[i], dp[j] + max(cuboids[i]))

        return max(dp)


print(Solution().maxHeight(cuboids=[[50, 45, 20], [95, 37, 53], [45, 23, 12]]))
# 输出：190
# 解释：
# 第 1 个长方体放在底部，53x37 的一面朝下，高度为 95 。
# 第 0 个长方体放在中间，45x20 的一面朝下，高度为 50 。
# 第 2 个长方体放在上面，23x12 的一面朝下，高度为 45 。
# 总高度是 95 + 50 + 45 = 190 。
