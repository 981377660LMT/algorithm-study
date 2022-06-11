from typing import List
from itertools import zip_longest


class Solution:
    def validWordSquare(self, words: List[str]) -> bool:
        #  看每一列的组合
        for index, col in enumerate(zip_longest(*words, fillvalue='')):
            if words[index] != ''.join(col):
                return False
        return True


# zip相当于竖向拼接列
print(*zip_longest("abcd", "bnrt", "crm", "dt"))
print(*zip("abcd", "bnrt", "crm", "dt"))

