from typing import List

# 你需要找到 s 中最多数目的非空子字符串，满足如下条件：
# 1.互不重叠
# 2.如果一个子字符串包含字符 char ，那么 s 中所有 char 字符都应该在这个子字符串中。(要么不取,要么取完)


class Solution:
    def maxNumOfSubstrings(self, s: str) -> List[str]:
        ...


print(Solution().maxNumOfSubstrings(s="adefaddaccc"))
# 输出：["e","f","ccc"]
# 解释：下面为所有满足第二个条件的子字符串：
# [
#   "adefaddaccc"
#   "adefadda",
#   "ef",
#   "e",
#   "f",
#   "ccc",
# ]
# 如果我们选择第一个字符串，那么我们无法再选择其他任何字符串，所以答案为 1 。如果我们选择 "adefadda" ，剩下子字符串中我们只可以选择 "ccc" ，它是唯一不重叠的子字符串，所以答案为 2 。同时我们可以发现，选择 "ef" 不是最优的，因为它可以被拆分成 2 个子字符串。所以最优解是选择 ["e","f","ccc"] ，答案为 3 。不存在别的相同数目子字符串解。

# 太难了
