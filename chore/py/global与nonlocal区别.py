# Global variables in recursion

# count=0
def foo(s):
    count = 0

    def bar(s):
        global count  # 注意此时全局变量里没有count 会报错 要用nonlocal
        if len(s) != 0:
            count += 1

    bar(s)
    print(count)


foo("as")


# !global查找比nonlocal查找快
# 参考1932. 合并多棵二叉搜索树
