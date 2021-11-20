# # 不断弹出（片段内）无重复的上升、下降序列，直到用尽字母。终于看懂了……

from collections import Counter


class Solution:
    def sortString(self, s: str) -> str:
        arr = sorted([num, freq] for num, freq in Counter(s).items())
        res = []

        while len(res) < len(s):
            for i in range(len(arr)):
                if arr[i][1] > 0:
                    res.append(arr[i][0])
                    arr[i][1] -= 1

            for i in range(len(arr)):
                if arr[~i][1] > 0:
                    res.append(arr[~i][0])
                    arr[~i][1] -= 1

        return ''.join(res)


print(Solution().sortString("aaaabbbbcccc"))
