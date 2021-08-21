import itertools
import pprint

a = itertools.groupby("alexe")
b = itertools.groupby("aaleex")
# 当参数长度不一时，zip和较短的保持一致，itertools.zip_longest()和较长的保持一致。
# 可以使用fillvalue来制定那些缺失值的默认值。
c = itertools.zip_longest(a, b, fillvalue=666)

pprint.pprint([(item0, item1) for item0, item1 in c])
