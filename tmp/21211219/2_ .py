from typing import List


class Solution:
    def addSpaces(self, s: str, spaces: List[int]) -> str:
        arr = list(s)
        j = len(spaces) - 1
        for i in range(len(s) - 1, -1, -1):
            if i == spaces[j]:
                arr[i] = ' ' + arr[i]
                j -= 1
        return ''.join(arr)


print(Solution().addSpaces(s="LeetcodeHelpsMeLearn", spaces=[8, 13, 15]))
print(Solution().addSpaces(s="icodeinpython", spaces=[1, 5, 7, 9]))
