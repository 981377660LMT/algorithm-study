# 寻找扑克牌中的对子
# 请你写一个函数帮他找到手牌中所有的对子
# （炸即大小王不是对子，本题中多于两张的不算对子），

from collections import defaultdict


order1 = {"H": 0, "S": 1, "C": 2, "D": 3}
order2 = {
    "3": 0,
    "4": 1,
    "5": 2,
    "6": 3,
    "7": 4,
    "8": 5,
    "9": 6,
    "10": 7,
    "J": 8,
    "Q": 9,
    "K": 10,
    "A": 11,
    "2": 12,
    "L": 13,
    "B": 14,
}


def findPair(s: str) -> str:
    arr = list(s)
    colors, counts = arr[::2], arr[1::2]
    mp = defaultdict(list)
    for color, count in zip(colors, counts):
        mp[count].append((color, count))
    res = [c for cards in mp.values() for c in cards if len(cards) == 2]
    res.sort(key=lambda x: (-order2[x[1]], order1[x[0]]))
    return "".join(map(lambda x: x[0] + x[1], res))


# https://leetcode.cn/contest/season/2022-fall/?modal=invite-vote&teamSlug=ting-qu-washeng-yi-pian
print(findPair("C4HAS5HQC6CJS3D3S6H3HKS7D7C9HJH2H8"))
