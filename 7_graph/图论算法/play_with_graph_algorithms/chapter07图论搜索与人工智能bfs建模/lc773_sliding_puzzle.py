from collections import deque


class Solution:
    def sliding_puzzle(self, board):
        """
        :type board: List[List[int]]
        :rtype: int
        """
        # 标准BFS题目
        m, n = len(board), len(board[0])
        target = '123450'
        # 核心: dirs数组的第i个元素表示它可以跟周围哪些位置交换（一维数组）
        # 比如0号下标，只能和右边的1下标和下面的3下标交换
        dirs = [
            [1, 3],
            [0, 2, 4],
            [1, 5],
            [0, 4],
            [1, 3, 5],
            [2, 4],
        ]
        
        # 初始化开始的状态
        start = ''
        for i in range(m):
            for j in range(n):
                start += str(board[i][j])

        queue = deque()
        queue.append(start)
        visited = set()

        res = 0
        while queue:
            # 注意这里重点要分层遍历才能更新res！！！！
            # 所以要for循环而不是单纯的每次popleft
            q_len = len(queue)
            for _ in range(q_len):
                curr = queue.popleft()
                if curr == target:
                    return res
                zero_inx = curr.find('0')
                for next_pos in dirs[zero_inx]:
                    # python里面string是immutable
                    # 这里只好先转成list处理
                    curr_list = list(curr)
                    curr_list[next_pos], curr_list[zero_inx] = curr_list[zero_inx], curr_list[next_pos]
                    new_str = ''.join(curr_list)
                    if new_str in visited:
                        continue
                    queue.append(new_str)
                    visited.add(new_str)
            res += 1

        # 遍历完了queue都没有找到解答
        # 说明做不到，只能返回-1
        return -1


if __name__ == '__main__':
    sol = Solution()
    data =[[4, 1, 2],[5, 0, 3]]
    print(sol.sliding_puzzle(data))