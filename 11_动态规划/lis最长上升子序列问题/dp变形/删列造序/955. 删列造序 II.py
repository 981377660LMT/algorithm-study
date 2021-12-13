from typing import List

# 选取一个删除索引序列，对于 strs 中的每个字符串，删除对应每个索引处的字符。
# 在执行删除操作之后，最终得到的数组的元素是按 字典序（strs[0] <= strs[1] <= strs[2] ... <= strs[n - 1]）排列的，
# 然后请你返回 answer.length 的最小可能值。


class Solution:
    def minDeletionSize(self, strs: List[str]) -> int:
        row = len(strs)
        res = 0
        curCol = [''] * row
        for col in zip(*strs):
            mergedCol = list(zip(curCol, col))
            if mergedCol == sorted(mergedCol):
                curCol = mergedCol
            else:
                res += 1
        return res


print(Solution().minDeletionSize(strs=["ca", "bb", "ac"]))
# 输出：1
# 解释：
# 删除第一列后，strs = ["a", "b", "c"]。
# 现在 strs 中元素是按字典排列的 (即，strs[0] <= strs[1] <= strs[2])。
# 我们至少需要进行 1 次删除，因为最初 strs 不是按字典序排列的，所以答案是 1。

