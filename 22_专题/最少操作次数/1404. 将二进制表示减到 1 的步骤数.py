class Solution:
    def numSteps2(self, s: str) -> int:
        num = int(s, 2)
        res = 0
        while num != 1:
            if num & 1:
                num += 1
            else:
                num >>= 1
            res += 1
        return res

    def numSteps(self, s: str) -> int:
        carry = count = 0

        for i in range(len(s) - 1, 0, -1):
            if s[i] == '0':
                count += 1 + carry
            else:
                count += 2 - carry
                carry = 1

        return count + carry


print(Solution().numSteps("1101"))
# 输入：s = "1101"
# 输出：6
# 解释："1101" 表示十进制数 13 。
# Step 1) 13 是奇数，加 1 得到 14
# Step 2) 14 是偶数，除 2 得到 7
# Step 3) 7  是奇数，加 1 得到 8
# Step 4) 8  是偶数，除 2 得到 4
# Step 5) 4  是偶数，除 2 得到 2
# Step 6) 2  是偶数，除 2 得到 1

