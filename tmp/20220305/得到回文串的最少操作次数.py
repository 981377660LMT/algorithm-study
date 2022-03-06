from typing import List, Tuple
from collections import defaultdict

# 每一次 操作 ，你可以选择 s 中两个 相邻 的字符，并将它们交换。
# 请你返回将 s 变成回文串的 最少操作次数 。

# 1 <= s.length <= 2000

# 双指针

# 从左到右依次处理。对于每个字母，从右向左找第一个与他匹配的字母：
# 如果找到，且不是它本身，则将找到的这个字母移动到右侧的对称位置；
# 如果找到，但是当前字母本身，这说明当前字母需要最后放在字符串的正中间。注意这里不能一步到位，而只能先将其向右移动一步，以避免进行多余操作。操作完成后，我们需要继续处理当前位置（所以在下面的实现中，有 i-- 这一步）。


class Solution:
    def minMovesToMakePalindrome(self, s: str) -> int:
        """从字符串头(head)和尾(tail)各取一个字符，判断是否相同，
      若相同则去除，不同则比较两个字符变换成相同所需的最小代价，并将相应的字符去掉，循环直至head==tail为止。
      """

        def dfs(s):
            if len(s) < 2:
                return 0
            n = len(s)
            for j in range(n - 1, 0, -1):
                if s[j] == s[0]:
                    return n - 1 - j + dfs(s[1:j] + s[j + 1 :])
            for i in range(n - 1):
                if s[i] == s[-1]:
                    return i + dfs(s[:i] + s[i + 1 : -1])

        return dfs(s)


# 163,42
print(Solution().minMovesToMakePalindrome("scpcyxprxxsjyjrww"))
