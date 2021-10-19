from math import floor, comb
from collections import namedtuple, Counter
from itertools import chain, combinations


def sort_by_indexes(lst, indexes, reverse=False):
    return [val for (_, val) in sorted(zip(indexes, lst), key=lambda x: x[0], reverse=reverse)]


a = ['eggs', 'bread', 'oranges', 'jam', 'apples', 'milk']
b = [3, 2, 6, 4, 1, 5]
sort_by_indexes(a, b)  # ['apples', 'bread', 'eggs', 'jam', 'milk', 'oranges']


def map_dictionary(itr, fn):
    return dict(zip(itr, map(fn, itr)))


map_dictionary([1, 2, 3], lambda x: x * x)  # { 1: 1, 2: 4, 3: 9 }


def symmetric_difference_by(a, b, fn):
    (_a, _b) = (set(map(fn, a)), set(map(fn, b)))
    return [item for item in a if fn(item) not in _b] + [item for item in b if fn(item) not in _a]


symmetric_difference_by([2.1, 1.2], [2.3, 3.4], floor)  # [1.2, 3.4]


Point = namedtuple('Point', ['x', 'y', 'z'], defaults=[1])
a = Point(1, 1, 0)
# a.x = 1, a.y = 1, a.z = 0

# Default value used for `z`
b = Point(2, 2)
# b.x = 2, b.y = 2, b.z = 1 (default)


def is_anagram(s1, s2):
    return Counter(c.lower() for c in s1 if c.isalnum()) == Counter(
        c.lower() for c in s2 if c.isalnum()
    )


is_anagram('#anagram', 'Nag a ram!')  # True


def powerset(iterable):
    s = list(iterable)
    return list(chain.from_iterable(combinations(s, r) for r in range(len(s) + 1)))


powerset([1, 2])  # [(), (1,), (2,), (1, 2)]


# 如果 num 在范围(a，b)内，则返回 num。
# 否则，返回范围内最接近的数字。
def clamp_number(num, a, b):
    return max(min(num, max(a, b)), min(a, b))


clamp_number(2, 3, 5)  # 3
clamp_number(1, -1, -5)  # -1


# capitalize_every_word
def capitalize_every_word(s: str):
    return s.title()


capitalize_every_word('hello world!')  # 'Hello World!'


# 二项式系数
def binomial_coefficient(n, k):
    return comb(n, k)


binomial_coefficient(8, 2)  # 28

