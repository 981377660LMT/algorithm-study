# 求整个num数组中找出长度最长的金字塔(山脉)数组，如果金字塔(山脉)数组不存在，请输出0
class Solution:
    def getMaxLength(self, n, num):
        # write code here
        up = [0] * n
        down = [0] * n

        for i in range(n - 1):
            if num[i] < num[i + 1]:
                up[i + 1] = up[i] + 1

        for i in range(n - 2, -1, -1):
            if num[i] > num[i + 1]:
                down[i] = down[i + 1] + 1

        return max(up[i] + down[i] + 1 for i in range(n))


print(Solution().getMaxLength(8, [1, 2, 3, 4, 5, 4, 4, 2]))

