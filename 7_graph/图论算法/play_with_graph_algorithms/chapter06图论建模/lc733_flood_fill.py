from collections import deque


class Solution:

    def flood_fill(self, image, sr, sc, new_color):

        if not image or not image[0]:
            return image
        
        m, n = len(image), len(image[0])
        if not 0 <= sr < m or not 0 <= sc < n:
            return image
        
        old_color = image[sr][sc]
        queue = deque()
        queue.append((sr, sc))
        
        dirs = [(0, 1), (1, 0), (0, -1), (-1, 0)]
        
        while queue:
            ci, cj = queue.popleft()
            image[ci][cj] = new_color
            for di, dj in dirs:
                newi, newj = ci + di, cj + dj
                # 这里别忘了如果新点不是旧颜色（old_color）
                # 或者已经是新颜色（new_color）时候，都不需要入队
                # 尤其是后者，如果已经是新颜色还入队的话
                # 会使得死循环(反复将这个点入队只因为它不是旧颜色)
                if not 0 <= newi < m or not 0 <= newj < n or \
                    image[newi][newj] != old_color or \
                    image[newi][newj] == new_color:
                    continue
                queue.append((newi, newj))
                image[newi][newj] = new_color
        
        return image
 
                