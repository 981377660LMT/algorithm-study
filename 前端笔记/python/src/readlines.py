# 有一个jsonline格式的文件file.txt大小约为10K


from typing import Any


def process(*arg: Any):
    ...


def get_lines():
    with open('file.txt', 'rb') as f:
        return f.readlines()


if __name__ == '__main__':
    for e in get_lines():
        process(e)  # 处理每一行数据


# 现在要处理一个大小为10G的文件，但是内存只有4G，如果在只修改get_lines 函数而其他代码保持不变的情况下，应该如何实现？需要考虑的问题都有那些？
def get_lines_2():
    with open('file.txt', 'rb') as f:
        yield from f.readlines()
