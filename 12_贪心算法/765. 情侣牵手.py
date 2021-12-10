from typing import List

# len(row) 是偶数且数值在 [4, 60]范围内。
# N 对情侣坐在连续排列的 2N 个座位上，想要牵到对方的手。
# 计算最少交换座位的次数，以便每对情侣可以并肩坐在一起。
# 一次交换可选择任意两人，让他们站起来交换座位。


# 配对技巧 x的对象是x^1 => (0,1) (2,3)
class Solution:
    # 贪心:我们遍历每个偶数位置 2 * i ，把它的对象安排到它右边的奇数位置 2 * i + 1。
    def minSwapsCouples(self, row: List[int]) -> int:

        n = len(row)
        res = 0
        for i in range(0, n, 2):
            if row[i] == row[i + 1] ^ 1:
                continue
            # 找到i的伴侣
            for j in range(i + 1, n):
                if row[i] == row[j] ^ 1:
                    row[i + 1], row[j] = row[j], row[i + 1]
            res += 1
        return res


print(Solution().minSwapsCouples(row=[0, 2, 1, 3]))
# 输出: 1
# 解释: 我们只需要交换row[1]和row[2]的位置即可。
