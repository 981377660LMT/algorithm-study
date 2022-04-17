from typing import List


def getNext(pattren: str) -> List[int]:
    next = [0] * len(pattren)
    j = 0

    for i in range(1, len(pattren)):
        while j and pattren[i] != pattren[j]:
            j = next[j - 1]

        if pattren[i] == pattren[j]:
            j += 1

        next[i] = j

    return next


# next[i]表示[:i+1]这一段字符串中最长公共前后缀(不是原串)的长度
next = getNext('aabaabaabaab')
print(next)
