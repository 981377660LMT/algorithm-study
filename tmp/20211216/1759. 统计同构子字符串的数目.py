# 输入：s = "abbcccaa"
# 输出：13
# 解释：同构子字符串如下所列：
# "a"   出现 3 次。
# "aa"  出现 1 次。
# "b"   出现 2 次。
# "bb"  出现 1 次。
# "c"   出现 3 次。
# "cc"  出现 2 次。
# "ccc" 出现 1 次。
# 3 + 1 + 2 + 1 + 3 + 2 + 1 = 13
from itertools import groupby


class Solution:
    def countHomogenous(self, s: str) -> int:
        counter = [len(list(group)) for _, group in groupby(s)]
        return sum(count * (count + 1) // 2 for count in counter) % int(1e9 + 7)

