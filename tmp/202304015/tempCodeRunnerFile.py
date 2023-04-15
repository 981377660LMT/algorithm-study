
        for i in range(len(sufMax) - 2, -1, -1):
            if sufMax[i] < sufMax[i + 1]:
                sufMax[i] = sufMax[i + 1]