# design git (设计git,可持久化栈)

# 现在给定一个空数组A 和一系列操作,操作有四种

# 1. ADD x : 在A的末尾添加一个元素x (git commit)
# 2. DELETE : 删除A的末尾元素 (git revert)
# 3. SAVE y : 保存当前数组A的状态到y分支 (git checkout -b)
# 4. LOAD z : 将z分支加载到当前数组A (git merge)

# !每个操作结束后,输出A数组的末尾元素(当前状态)
# !如果A为空,输出-1


from collections import defaultdict
from PersistentStack import PersistentStack


if __name__ == "__main__":
    git = defaultdict(lambda: PersistentStack.default())  # !各个分支上的版本
    curStack = PersistentStack.default()  # !当前版本
    q = int(input())
    for _ in range(q):
        kind, *args = input().split()
        if kind == "ADD":
            value = int(args[0])
            curStack = curStack.push(value)
        elif kind == "DELETE":
            curStack = curStack.pop()
        elif kind == "SAVE":
            branch = args[0]
            git[branch] = curStack
        elif kind == "LOAD":
            branch = args[0]
            curStack = git[branch]

        curValue = curStack.top()
        print(curValue if curValue is not None else -1)
