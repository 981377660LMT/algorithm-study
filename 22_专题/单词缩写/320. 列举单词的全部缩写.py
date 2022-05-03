# 1 <= word.length <= 15
from typing import List


class Solution:
    def generateAbbreviations(self, word: str) -> List[str]:
        def bt(index: int, count: int, path: List[str]) -> None:
            if index == len(word):
                if count:
                    path.append(str(count))
                res.append(''.join(path))
                if count:
                    path.pop()
                return

            bt(index + 1, count + 1, path)

            path.append(str(count) if count else '')
            path.append(word[index])
            bt(index + 1, 0, path)
            path.pop()
            path.pop()

        res = []
        bt(0, 0, [])
        return res


print(Solution().generateAbbreviations("word"))
