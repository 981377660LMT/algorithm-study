from typing import List


print(48 ^ 1)
print(49 ^ 1)
# 01互换

# // 康托对角线。
# // 和第0个串的[0]不同，和第1个串的[1]不同，和第2个串的[2]不同，……
# // 由此构造出来的串就和所有已有的串都不同了
class Solution:
    def findDifferentBinaryString(self, nums: List[str]) -> str:
        return ''.join([chr(ord(num[i]) ^ 1) for i, num in enumerate(nums)])


# 输入：nums = ["111","011","001"]
# 输出："101"
# 解释："101" 没有出现在 nums 中。"000"、"010"、"100"、"110" 也是正确答案。
