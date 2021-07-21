from collections import deque


class Solution:

    def open_lock(self, deadends, target):
        deadends = set(deadends)
        
        if target in deadends:
            return -1

        if '0000' in deadends:
            return -1

        if '0000' == target:
            return 0

        queue = deque()
        queue.append('0000')
        visited = dict()
        visited['0000'] = 0

        while queue:
            curs = queue.popleft()
            nexts = []
            for i in range(4):
                int_ch = int(curs[i])
                new_int_ch1 = (int_ch + 1) % 10
                new_int_ch2 = (int_ch + 9) % 10
                nexts.append(curs[:i] + str(new_int_ch1) + curs[i + 1:])
                nexts.append(curs[:i] + str(new_int_ch2) + curs[i + 1:])

            for next_ in nexts:
                if next_ not in deadends and visited.get(next_) is None:
                    queue.append(next_)
                    visited[next_] = visited[curs] + 1
                    if next_ == target:
                        return visited[next_]
        
        return -1


if __name__ == '__main__':
    sol = Solution()
    data = []
    print(sol.open_lock(['0201','0101','0102','1212','2002'], target = '0202'))