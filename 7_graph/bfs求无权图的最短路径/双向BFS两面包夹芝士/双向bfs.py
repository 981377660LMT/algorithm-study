from typing import Callable, Iterable, List, TypeVar


S = TypeVar("S")


def biBfs(start: S, target: S, getNextStates: Callable[[S], Iterable[S]]) -> int:
    """双向bfs.如果不存在则返回-1."""
    queue1, queue2 = set([start]), set([target])
    visited = set()
    steps = 0

    while queue1 and queue2:
        if len(queue1) > len(queue2):
            queue1, queue2 = queue2, queue1

        nextQueue = set()
        for cur in queue1:
            if cur in queue2:
                return steps
            if cur in visited:
                continue
            visited.add(cur)
            for next in getNextStates(cur):
                nextQueue.add(next)

        steps += 1
        queue1, queue2 = queue2, nextQueue

    return -1


def biBfsPath(start: S, target: S, getNextStates: Callable[[S], Iterable[S]]) -> List[S]:
    """双向bfs,返回路径."""
    queue = [set([start]), set([target])]
    pre = [dict(), dict()]
    visited = [set([start]), set([target])]

    def restorePath(mid: S) -> List[S]:
        pre1, path1, cur1 = pre[0], [mid], mid
        while True:
            p = pre1.get(cur1)
            if p is None:
                break
            cur1 = p
            path1.append(cur1)
        pre2, path2, cur2 = pre[1], [], mid
        while True:
            p = pre2.get(cur2)
            if p is None:
                break
            cur2 = p
            path2.append(cur2)
        return path1[::-1] + path2

    curQueue, curVisited, curPre, otherQueue = None, None, None, None
    while queue[0] and queue[1]:
        qi = 1 if len(queue[0]) > len(queue[1]) else 0
        nextQueue = set()
        curQueue, curVisited, curPre = queue[qi], visited[qi], pre[qi]
        otherQueue = queue[qi ^ 1]
        for cur in curQueue:
            if cur in otherQueue:
                return restorePath(cur)
            for next in getNextStates(cur):
                if next in curVisited:
                    continue
                curVisited.add(next)
                nextQueue.add(next)
                curPre[next] = cur
        queue[qi] = nextQueue
    return []


if __name__ == "__main__":
    # 127. 单词接龙
    # https://leetcode.cn/problems/word-ladder/description/
    class Solution:
        def ladderLength(self, beginWord: str, endWord: str, wordList: List[str]) -> int:
            def getNextStates(word: str) -> Iterable[str]:
                for i in range(len(word)):
                    for c in "abcdefghijklmnopqrstuvwxyz":
                        newWord = word[:i] + c + word[i + 1 :]
                        if newWord in wordSet:
                            yield newWord

            wordSet = set(wordList)
            if endWord not in wordSet:
                return 0

            # path = biBfsPath(beginWord, endWord, getNextStates)
            # return len(path)
            res = biBfs(beginWord, endWord, getNextStates)
            return res + 1 if res != -1 else 0
