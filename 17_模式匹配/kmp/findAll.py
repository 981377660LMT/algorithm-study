from typing import List


def findAll(string: str, target: str) -> List[int]:
    """找到所有匹配的字符串起始位置"""
    start = 0
    res = []
    while True:
        pos = string.find(target, start)
        if pos == -1:
            break
        else:
            res.append(pos)
            start = pos + 1

    return res


if __name__ == '__main__':
    print(findAll('abcdefgabcabc', 'abc'))


# https://www.quora.com/What-is-the-time-complexity-of-std-string-find-in-C++
