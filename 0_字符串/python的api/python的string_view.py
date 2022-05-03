# python的string_view
import string

nums = ['11', '12', '12']
s = 'abcaab'
view = memoryview(bytes(s, encoding='utf-8'))
# 用memoryview做字符串哈希
print(view[:2] == view[-2:])
