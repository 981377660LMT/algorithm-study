# 给你一个字符串 num
# 其中，num 表示一个很大的整数，字符串中的每个字符依次对应整数上的各个 数位 。
# 你可以交换这个整数相邻数位的数字 最多 k 次。
# 请你返回你能得到的最小整数，并以字符串形式返回。

# 1 <= num.length <= 30000
# 1 <= k <= 10^9

# 贪心：找的数字越小越靠前
# 直接暴力过了
class Solution:
    def minInteger1(self, num: str, k: int) -> str:
        if k <= 0 or not num:
            return num
        for i in range(0, 10):
            index = num.find(str(i))
            if index < 0:
                continue
            if index <= k:
                return str(i) + self.minInteger1(num[0:index] + num[index + 1 :], k - index)

        return num


print(Solution().minInteger1("4321", 4))
