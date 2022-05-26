from typing import List

# 自除数 是指可以被它包含的每一位数除尽的数。


# 例如，128 是一个自除数，因为 128 % 1 == 0，128 % 2 == 0，128 % 8 == 0。
# 自除数不允许包含 0 。
class Solution:
    def selfDividingNumbers(self, left: int, right: int) -> List[int]:
        res = []
        for num in range(left, right + 1):
            s = str(num)
            if '0' in s:
                continue
            for char in s:
                if num % int(char) != 0:
                    break
            else:
                res.append(num)

        return res


class Solution:
    def isPathCrossing(self, path: str) -> bool:
        x, y = 0, 0
        s = set((x, y))
        for i in path:
            if i == "N":
                y += 1
            elif i == "E":
                x += 1
            elif i == "S":
                y -= 1
            elif i == "W":
                x -= 1

            if (x, y) in s:
                return True
            else:
                s.add((x, y))
        return False
