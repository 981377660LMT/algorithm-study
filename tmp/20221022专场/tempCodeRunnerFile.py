            diff = need[char] - visited[char]
            while stack and stack[-1] > char and remain[stack[-1]] >= diff:
                last = stack.pop()
                visited[last] -= 1
            stack.append(char)
            visited[char] += 1
