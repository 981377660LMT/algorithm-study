# https://atcoder.jp/contests/abc301/tasks/abc301_c


# 给定两个长度相等的字符串s1,s2包含小写字母和@，
# 问能否将@替换成atcoder中的一个字母，可以通过对其中一个字符串排序，使得两者相同。


from collections import Counter


def atcoderCards(s: str, t: str) -> bool:
    at1, at2 = 0, 0
    C1, C2 = Counter(), Counter()
    for c in s:
        if c == "@":
            at1 += 1
        else:
            C1[c] += 1
    for c in t:
        if c == "@":
            at2 += 1
        else:
            C2[c] += 1

    diff1 = C1 - C2  # 差集
    diff2 = C2 - C1
    diff = diff1 + diff2
    if not all(c in "atcoder" for c in diff):
        return False
    return at1 + at2 >= sum(diff.values())


if __name__ == "__main__":
    s = input()
    t = input()
    res = atcoderCards(s, t)
    print("Yes" if res else "No")
