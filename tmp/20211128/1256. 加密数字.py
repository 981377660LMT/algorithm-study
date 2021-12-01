class Solution:
    def encode(self, num: int) -> str:
        return bin(num + 1)[3:]


print(Solution().encode(num=107))
# 输入：num = 107
# 输出："101100"
print(bin(2))
