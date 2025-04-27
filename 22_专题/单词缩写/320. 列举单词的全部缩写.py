# 1 <= word.length <= 15
from typing import List


class Solution:
    def generateAbbreviations(self, word: str) -> List[str]:
        path: List[str] = []

        def bt(index: int, count: int) -> None:
            if index == len(word):
                if count:
                    path.append(str(count))
                res.append("".join(path))
                if count:
                    path.pop()
                return

            bt(index + 1, count + 1)

            if count > 0:
                path.append(str(count))
            path.append(word[index])
            bt(index + 1, 0)
            path.pop()
            if count > 0:
                path.pop()

        res = []
        bt(0, 0)
        return res


print(Solution().generateAbbreviations("word"))
