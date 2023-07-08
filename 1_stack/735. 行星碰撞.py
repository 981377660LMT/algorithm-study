# 735. 行星碰撞
# https://leetcode.cn/problems/asteroid-collision/
# 找出碰撞后剩下的所有行星。
# !碰撞规则：两个行星相互碰撞，较小的行星会爆炸。如果两颗行星大小相同，则两颗行星都会爆炸。
# 两颗移动方向相同的行星，永远不会发生碰撞。


from typing import List


class Solution:
    def asteroidCollision(self, asteroids: List[int]) -> List[int]:
        stack = []
        for v in asteroids:
            stack.append(v)
            while len(stack) >= 2 and stack[-2] > 0 and stack[-1] < 0:
                pre, cur = stack[-2], -stack[-1]
                if pre > cur:
                    stack.pop()
                elif pre < cur:
                    top = stack.pop()
                    stack.pop()
                    stack.append(top)
                else:
                    stack.pop()
                    stack.pop()
        return stack
