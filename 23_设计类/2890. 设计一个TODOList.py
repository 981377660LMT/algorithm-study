# 2890. 设计一个Todo List

# 工程类设计:
# !1. 面向对象设计, 对象上有什么方法和属性
# !2. 重视id, 用哈希表保存 id => 对象
# !3. 代码健壮性, 多关注边界条件

from typing import List
from collections import defaultdict
from sortedcontainers import SortedList
import itertools


class TodoList:
    def __init__(self):
        self._idToTodo = dict()
        self._userTodo = defaultdict(lambda: SortedList(key=lambda id: self._idToTodo[id].dueDate))
        self._idPool = itertools.count(1)

    def addTask(self, userId: int, taskDescription: str, dueDate: int, tags: List[str]) -> int:
        """注册一个任务, 返回taskId, taskId 从 1 开始全局自增."""
        id = next(self._idPool)
        todo = Todo(id, taskDescription, dueDate, tags)
        self._idToTodo[id] = todo
        self._userTodo[userId].add(id)
        return id

    def getAllTasks(self, userId: int) -> List[str]:
        """返回按照 dueDate 升序排列的所有未完成的任务."""
        res = []
        for id in self._userTodo[userId]:
            todo = self._idToTodo[id]
            res.append(todo.description)
        return res

    def getTasksForTag(self, userId: int, tag: str) -> List[str]:
        """按照 dueDate 升序排列返回所有未完成的任务, 且任务中包含 tag."""
        res = []
        for id in self._userTodo[userId]:
            todo = self._idToTodo[id]
            if tag in todo.tags:
                res.append(todo.description)
        return res

    def completeTask(self, userId: int, taskId: int) -> None:
        """删除 taskId 任务."""
        if taskId not in self._idToTodo or taskId not in self._userTodo[userId]:
            return
        self._userTodo[userId].remove(taskId)
        self._idToTodo.pop(taskId)

    def __repr__(self):
        res = []
        for id, todo in self._idToTodo.items():
            res.append(f"{id}: {todo}")
        userInfo = []
        for userId, todoIds in self._userTodo.items():
            userInfo.append(f"{userId}: {list(todoIds)}")
        return f"AllTodo: {res}\nUserTodo: {userInfo}"


class Todo:
    __slots__ = ("id", "description", "dueDate", "tags")

    def __init__(self, id: int, description: str, dueDate: int, tags: List[str]):
        self.id = id
        self.description = description
        self.dueDate = dueDate
        self.tags = set(tags)

    def __repr__(self) -> str:
        return f"Todo(id:{self.id}, description:{self.description}, dueDate:{self.dueDate}, tags:{self.tags})"


# Your TodoList object will be instantiated and called as such:
# obj = TodoList()
# param_1 = obj.addTask(userId,taskDescription,dueDate,tags)
# param_2 = obj.getAllTasks(userId)
# param_3 = obj.getTasksForTag(userId,tag)
# obj.completeTask(userId,taskId)
# Input
# ["TodoList", "addTask", "addTask", "getAllTasks", "getAllTasks", "addTask", "getTasksForTag", "completeTask", "completeTask", "getTasksForTag", "getAllTasks"]
# [[], [1, "Task1", 50, []], [1, "Task2", 100, ["P1"]], [1], [5], [1, "Task3", 30, ["P1"]], [1, "P1"], [5, 1], [1, 2], [1, "P1"], [1]]
# Output
# [null, 1, 2, ["Task1", "Task2"], [], 3, ["Task3", "Task2"], null, null, ["Task3"], ["Task3", "Task1"]]

if __name__ == "__main__":
    todoList = TodoList()
    todoList.addTask(1, "Task1", 50, [])
    todoList.addTask(1, "Task2", 100, ["P1"])
    print(todoList.getAllTasks(1))
    todoList.addTask(1, "Task3", 30, ["P1"])
    print(todoList.getTasksForTag(1, "P1"))
    todoList.completeTask(5, 1)
    todoList.completeTask(1, 2)
    print(todoList.getTasksForTag(1, "P1"))
    print(todoList.getAllTasks(1))

    print(todoList)
