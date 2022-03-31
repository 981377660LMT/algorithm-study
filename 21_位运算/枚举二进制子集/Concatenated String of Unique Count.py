from itertools import chain, combinations

# n,m<=20
# 选择word使其无相同字符，求最长长度
class Solution:
    def powerset(self, iterable):
        # powerset([1,2,3]) --> () (1,) (2,) (3,) (1,2) (1,3) (2,3) (1,2,3)
        s = list(iterable)
        return chain.from_iterable(combinations(s, count) for count in range(len(s) + 1))

    def solve(self, words):
        res = 0
        for sub in self.powerset(words):
            string = "".join(sub)
            if len(string) == len(set(string)):
                res = max(res, len(string))
        return res
