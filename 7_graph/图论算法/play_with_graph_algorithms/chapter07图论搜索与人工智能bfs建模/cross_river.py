from collections import deque


class Solution:

    def cross_river(self, begin):
        queue = deque()
        visited = dict()
        queue.append(begin)
        visited[''.join(str(i) for i in begin)] = None

        while queue:
            curr = queue.popleft()
            for i in range(4):
                if curr[0] != curr[i]:
                    continue
                next_ = curr[:]
                next_[0] = 1 - curr[0]
                next_[i] = 1 - curr[0]
                new_state = ''.join(str(i) for i in next_)
                if new_state not in visited:
                    if (
                        new_state[1] == new_state[2] and new_state[0] != new_state[1] or
                        new_state[2] == new_state[3] and new_state[0] != new_state[2]
                    ):
                        continue
                    queue.append(next_)
                    visited[new_state] = ''.join(str(i) for i in curr)
                    if new_state == '1111':
                        break

        path = []
        curr = '1111'
        path.append(curr)
        while visited.get(curr):
            path.append(visited.get(curr))
            curr = visited.get(curr)
        return path[::-1]


if __name__ == '__main__':
    # initial state:
    begin = [0, 0, 0, 0]
    sol = Solution()
    print(sol.cross_river(begin))