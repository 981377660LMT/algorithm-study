# python 默认递归调用栈1000
import sys


sys.setrecursionlimit(300000)

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
