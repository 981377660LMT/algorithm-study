# python 默认递归调用栈1000
import os
import sys


sys.setrecursionlimit(int(7e5))

i = 0


def recur():
    global i
    i += 1
    recur()


try:
    recur()
except Exception as e:
    i += 1
    print(f'Maximum stack size is {i}')
