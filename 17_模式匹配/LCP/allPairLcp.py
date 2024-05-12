# 所有数对的最长公共前缀之和


from typing import List


def allPairLcp(words: List[str]) -> int:
    root = dict()  # char -> (count, children)
    res = 0
    for word in words:
        cur = root
        for c in word:
            if c in cur:
                v = cur[c]
                res += v[0]
                v[0] += 1
            else:
                cur[c] = [1, dict()]
            cur = cur[c][1]
    return res


if __name__ == "__main__":
    n = int(input())
    words = input().split()
    print(allPairLcp(words))
