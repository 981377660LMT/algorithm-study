# 先用permutaions生成cookies的全排列，然后用combinations生成挡板。
# 这两步操作就相当于暴力枚举所有饼干分给所有孩子的操作，然后就简单计算了~

# 不过即便是简单粗暴的暴力，其实也有一些优化的点，比如：
# 可以用前缀和来优化每次累加
# permutations中的每个组，其实不要求有序。那么可以各种combinations


from itertools import combinations, pairwise, permutations
from typing import List

# n<=8 暴力


class Solution:
    def distributeCookies(self, cookies: List[int], k: int) -> int:
        n = len(cookies)
        splits = list(range(1, n))
        res = int(1e20)
        for perm in permutations(cookies):
            for splitComb in combinations(splits, k - 1):
                cand = max(sum(perm[pre:cur]) for pre, cur in pairwise((0,) + splitComb + (n,)))
                if cand < res:
                    res = cand
        return res

