from typing import List

# 你需要找出并删除 不是按字典序升序排列的 列
# 返回你需要删除的列数。


class Solution:
    def minDeletionSize(self, strs: List[str]) -> int:
        return sum(list(col) != sorted(col) for col in zip(*strs))


print(Solution().minDeletionSize(strs=["cba", "daf", "ghi"]))
# 输出：1
# 解释：网格示意如下：
#   cba
#   daf
#   ghi
# 列 0 和列 2 按升序排列，但列 1 不是，所以只需要删除列 1 。
