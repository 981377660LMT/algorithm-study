from math import ceil


def minimize_max_diff(nums):
    n = len(nums)
    max_adjacent_known_diff = 0
    min_adjacent = None
    max_adjacent = None
    prev_known_index = None
    segments = []

    # Step 1: Identify known elements and calculate max_adjacent_known_diff
    for i in range(n):
        if nums[i] != -1:
            if prev_known_index is not None:
                diff = abs(nums[i] - nums[prev_known_index])
                max_adjacent_known_diff = max(max_adjacent_known_diff, diff)
                # Collect segments of -1s between known numbers
                if i - prev_known_index > 1:
                    segments.append((nums[prev_known_index], nums[i], i - prev_known_index - 1))
            prev_known_index = i
        else:
            # Collect adjacent known values to -1s
            if i > 0 and nums[i - 1] != -1:
                val = nums[i - 1]
                if min_adjacent is None:
                    min_adjacent = val
                    max_adjacent = val
                else:
                    min_adjacent = min(min_adjacent, val)
                    max_adjacent = max(max_adjacent, val)
            if i + 1 < n and nums[i + 1] != -1:
                val = nums[i + 1]
                if min_adjacent is None:
                    min_adjacent = val
                    max_adjacent = val
                else:
                    min_adjacent = min(min_adjacent, val)
                    max_adjacent = max(max_adjacent, val)

    if min_adjacent is None:
        # All elements are -1
        return 0

    # Step 2: Determine minimal possible maximum difference k for segments
    segment_max_diff = 0
    for a, b, m in segments:
        k_candidate = ceil(abs(a - b) / (m + 1))
        segment_max_diff = max(segment_max_diff, k_candidate)

    # Combine max differences
    k = max(max_adjacent_known_diff, segment_max_diff)

    # Step 3: Calculate replacement value m
    m = (min_adjacent + max_adjacent) // 2

    # Step 4: Replace -1s with m and calculate actual maximum difference
    nums_replaced = nums[:]
    for i in range(n):
        if nums_replaced[i] == -1:
            nums_replaced[i] = m

    actual_max_diff = 0
    for i in range(1, n):
        diff = abs(nums_replaced[i] - nums_replaced[i - 1])
        actual_max_diff = max(actual_max_diff, diff)

    # The final result is the minimal possible maximum difference
    result = max(k, actual_max_diff)
    return result


# [1,2,-1,10,8] => 4
print(minimize_max_diff([1, 2, -1, 10, 8]))
# [-1,-1,-1] => 0
print(minimize_max_diff([-1, -1, -1]))
# [-1,10,-1,8] => 1
print(minimize_max_diff([-1, 10, -1, 8]))
# [1,12] => 11
print(minimize_max_diff([1, 12]))
# [14,-1,-1,46] => 11
print(minimize_max_diff([14, -1, -1, 46]))
# [40,-1,-1,-1,79] => 13
print(minimize_max_diff([40, -1, -1, -1, 79]))
