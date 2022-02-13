# 把一个字符串的大写字母放到字符串的后面，各个字符的相对位置不变，且不能申请额外的空间。
# 注意：字符串不可变的语言是不可能不用额外空间的

# 题目的输入应该改为字符串数组，不用额外空间
from typing import List


def main(arr: List[str]) -> None:
    n = len(arr)
    for i in range(n - 1):
        for j in range(n - 1 - i):
            # 前面大写后面小写就要交换
            if arr[j].isupper() and arr[j + 1].islower():
                arr[j], arr[j + 1] = arr[j + 1], arr[j]


arr = ['a', 'b', 'A', 'c', 'D', 'e', 'F', 'G']
main(arr)
print(arr)
