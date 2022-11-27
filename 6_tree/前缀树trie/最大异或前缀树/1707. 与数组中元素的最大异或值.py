from typing import List
from XORTrieArray import useXORTrie


class Solution:
    def maximizeXor(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        # offline，排序，Trie树 贪心
        n = len(nums)
        nums.sort()

        qn = len(queries)
        Q = [(qv, limit, qi) for qi, (qv, limit) in enumerate(queries)]
        Q.sort(key=lambda x: x[1])

        res = [0] * qn
        xorTrie = useXORTrie(int(1e9))
        ni = 0

        for qv, limit, qi in Q:
            while ni < n and nums[ni] <= limit:  # 不超过m的都入树
                xorTrie.insert(nums[ni])
                ni += 1

            if ni == 0:  # 树中没有元素
                res[qi] = -1
            else:
                res[qi] = xorTrie.search(qv)

        return res
