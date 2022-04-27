from typing import List, Optional


class Manacher:
    '''马拉车算法 O(n)

    https://github.com/981377660LMT/algorithm-study/blob/master/17_%E6%A8%A1%E5%BC%8F%E5%8C%B9%E9%85%8D/%E9%A9%AC%E6%8B%89%E8%BD%A6%E6%8B%89%E9%A9%AC/Manacher.py
    '''

    def __init__(self, s: str):
        self._s = s
        self._n = len(s)
        self._odds = self._getOddRadius()
        self._evens = self._getEvenRadius()
        self._maxOdd1: Optional[List[int]] = None
        self._maxOdd2: Optional[List[int]] = None
        self._maxEven1: Optional[List[int]] = None
        self._maxEven2: Optional[List[int]] = None

    def getLongestOddStartsAt(self, index: int) -> int:
        """以s[index]开头的最长奇回文子串的长度
        """
        if self._maxOdd1 is None:
            self._maxOdd1 = [1] * self._n
            self._maxOdd2 = [1] * self._n
            for i, radius in enumerate(self._odds):
                start, end = i - radius + 1, i + radius - 1
                length = 2 * radius - 1
                self._maxOdd1[start] = max(self._maxOdd1[start], length)
                self._maxOdd2[end] = max(self._maxOdd2[end], length)

            # 根据左右更新端点
            for i in range(self._n):
                if i - 1 >= 0:
                    self._maxOdd1[i] = max(self._maxOdd1[i], self._maxOdd1[i - 1] - 2)
                if i + 1 < self._n:
                    self._maxOdd2[i] = max(self._maxOdd2[i], self._maxOdd2[i + 1] - 2)
        return self._maxOdd1[index]

    def getLongestOddEndsAt(self, index: int) -> int:
        """以s[index]结尾的最长奇回文子串的长度
        """
        if self._maxOdd2 is None:
            self._maxOdd1 = [1] * self._n
            self._maxOdd2 = [1] * self._n
            for i, radius in enumerate(self._odds):
                start, end = i - radius + 1, i + radius - 1
                length = 2 * radius - 1
                self._maxOdd1[start] = max(self._maxOdd1[start], length)
                self._maxOdd2[end] = max(self._maxOdd2[end], length)

            # 根据左右更新端点
            for i in range(self._n):
                if i - 1 >= 0:
                    self._maxOdd1[i] = max(self._maxOdd1[i], self._maxOdd1[i - 1] - 2)
                if i + 1 < self._n:
                    self._maxOdd2[i] = max(self._maxOdd2[i], self._maxOdd2[i + 1] - 2)
        return self._maxOdd2[index]

    def getLongestEvenStartsAt(self, index: int) -> int:
        """以s[index]开头的最长偶回文子串的长度
        """
        if self._maxEven1 is None:
            self._maxEven1 = [0] * self._n
            self._maxEven2 = [0] * self._n
            for i, radius in enumerate(self._evens):
                if radius == 0:
                    continue
                start = i - radius
                end = start + 2 * radius - 1
                length = 2 * radius
                self._maxEven1[start] = max(self._maxEven1[start], length)
                self._maxEven2[end] = max(self._maxEven2[end], length)

            # 根据左右更新端点
            for i in range(self._n):
                if i - 1 >= 0:
                    self._maxEven1[i] = max(self._maxEven1[i], self._maxEven1[i - 1] - 2)
                if i + 1 < self._n:
                    self._maxEven2[i] = max(self._maxEven2[i], self._maxEven2[i + 1] - 2)
        return self._maxEven1[index]

    def getLongestEvenEndsAt(self, index: int) -> int:
        """以s[index]结尾的最长偶回文子串的长度
        """
        if self._maxEven2 is None:
            self._maxEven1 = [0] * self._n
            self._maxEven2 = [0] * self._n
            for i, radius in enumerate(self._evens):
                if radius == 0:
                    continue
                start = i - radius
                end = start + 2 * radius - 1
                length = 2 * radius
                self._maxEven1[start] = max(self._maxEven1[start], length)
                self._maxEven2[end] = max(self._maxEven2[end], length)

            # 根据左右更新端点
            for i in range(self._n):
                if i - 1 >= 0:
                    self._maxEven1[i] = max(self._maxEven1[i], self._maxEven1[i - 1] - 2)
                if i + 1 < self._n:
                    self._maxEven2[i] = max(self._maxEven2[i], self._maxEven2[i + 1] - 2)
        return self._maxEven2[index]

    def _getOddRadius(self) -> List[int]:
        '''获取每个中心点的奇回文半径`radius`

        回文为`[pos-radius+1:pos+radius]`
        '''
        res = [0] * self._n
        left, right = 0, -1
        for i in range(self._n):
            k = 1 if i > right else min(res[left + right - i], right - i + 1)
            while 0 <= i - k and i + k < self._n and self._s[i - k] == self._s[i + k]:
                k += 1
            res[i] = k
            k -= 1
            if i + k > right:
                left = i - k
                right = i + k
        return res

    def _getEvenRadius(self) -> List[int]:
        '''获取每个中心点的偶回文半径`radius`

        回文为`[pos-radius:pos+radius]`
        '''
        res = [0] * self._n
        left, right = 0, -1
        for i in range(self._n):
            k = 0 if i > right else min(res[left + right - i + 1], right - i + 1)
            while 0 <= i - k - 1 and i + k < self._n and self._s[i - k - 1] == self._s[i + k]:
                k += 1
            res[i] = k
            k -= 1
            if i + k > right:
                left = i - k - 1
                right = i + k
        return res

    def __len__(self) -> int:
        return self._n


if __name__ == '__main__':
    m1 = Manacher('aaaab')
    assert m1.getLongestEvenStartsAt(0) == 4
    assert m1.getLongestEvenEndsAt(3) == 4
    assert m1.getLongestOddStartsAt(0) == 3
    assert m1.getLongestOddEndsAt(4) == 1

    m2 = Manacher("ggbswiymmlevedhkbdhntnhdbkhdevelmmyiwsbgg")
    res = [m2.getLongestOddStartsAt(i) for i in range(len(m2))]
    target = [
        41,
        39,
        37,
        35,
        33,
        31,
        29,
        27,
        25,
        23,
        21,
        19,
        17,
        15,
        13,
        11,
        9,
        7,
        5,
        3,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        3,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
    ]
    assert res == target

