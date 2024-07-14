
        #     state, count = queue2.popleft()
        #     if state in visited1:
        #         print(count + visited1[state])
        #         exit()
        #     for i in range(N + 1):
        #         if state[i] == 2 and state[i + 1] == 2:
        #             for j in range(N + 1):
        #                 if i == j:
        #                     continue
        #                 newState = list(state)
        #                 tmpI1, tmpI2 = newState[i], newState[i + 1]
        #                 tmpJ1, tmpJ2 = newState[j], newState[j + 1]
        #                 newState[i], newState[j] = state[j], state[i]
        #                 newState[i + 1], newState[j + 1] = state[j + 1], state[i + 1]
        #                 newState = tuple(newState)
        #                 if newState not in visited2:
        #                     visited2[newState] = count + 1
        #                     queue2.append((newS