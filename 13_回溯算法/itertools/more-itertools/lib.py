from itertools import filterfalse, compress, starmap

from more_itertools import (
    chunked,
    ichunked,
    unzip,
    chunked_even,
    prepend,
    powerset,
    ncycles,
    nth,
    nth_product,
    nth_permutation,
    nth_combination,
    nth_combination_with_replacement,
    filter_map,
    grouper,
    consume,
    take,
    tail,
    partitions,
    product_index,
    combination_index,
    permutation_index,
    combination_with_replacement_index,
    gray_product,
    outer_product,
    powerset_of_sets,
    random_product,
    random_permutation,
    random_combination,
    random_combination_with_replacement,
    side_effect,
    collapse,
    split_at,
    split_before,
    split_after,
    split_into,
    split_when,
    bucket,
    map_reduce,
    sort_together,
    seekable,
    set_partitions,
    filter_except,
    unique_to_each,
    numeric_range,
    make_decorator,
    tabulate,
    repeat_each,
    repeat_last,
    repeatfunc,
    batched,
    flatten,
    quantify,
    first,
    first_true,
    all_equal,
    unique,
    unique_justseen,
    unique_everseen,
    sliding_window,
    windowed,
    roundrobin,
    partition,
    subslices,
    substrings,
    iter_index,
    iter_except,
    sum_of_squares,
    reshape,
    convolve,
    sliced,
    distribute,
    divide,
    transpose,
    spy,
    peekable,
    stagger,
    windowed_complete,
    triplewise,
    count_cycle,
    intersperse,
    adjacent,
    pad_none,
    interleave_longest,
    interleave_evenly,
    partial_product,
    value_chain,
    ilen,
    dft,
    sample,
    consecutive_groups,
    run_length,
    join_mappings,
    is_sorted,
    all_unique,
    map_if,
    minmax,
    only,
    strip,
    iter_suppress,
    nth_or_last,
    last,
    unique_in_window,
    longest_common_prefix,
    distinct_permutations,
    distinct_combinations,
    nth_product,
    nth_permutation,
    nth_combination_with_replacement,
    circular_shifts,
    always_iterable,
    always_reversible,
    countable,
    with_iter,
    locate,
    difference,
    SequenceView,
    time_limited,
    replace,
    iter_index,
    iterate,
)


args = [(1, 2), (3, 4), (5, 6)]
print(list(starmap(lambda x, y: x + y, args)))


dates = [
    "2020-01-01",
    "2020-02-04",
    "2020-02-01",
    "2020-01-24",
    "2020-01-08",
    "2020-02-10",
    "2020-02-15",
    "2020-02-11",
]
counts = [1, 4, 3, 8, 0, 7, 9, 2]
print(list(filterfalse(lambda x: x < 5, counts)))
print(list(compress(dates, [c < 5 for c in counts])))

tree = [40, [25, [10, 3, 17], [32, 30, 38]], [78, 50, 93]]
print(list(collapse(tree)))

lines = [
    "erhgedrgh",
    "erhgedrghed",
    "esdrhesdresr",
    "ktguygkyuk",
    "-------------",
    "srdthsrdt",
    "waefawef",
    "ryjrtyfj",
    "-------------",
    "edthedt",
    "awefawe",
]
print(list(split_at(lines, lambda x: "-------------" in x)))


data = "This is example sentence for seeking back and forth".split()
it = seekable(data)
print(it.peek(3))


nums = [1, 2, 3, 4, 5]
tit = take(3, nums)
print(*tit)


print(*roundrobin("ABC", "D", "EF"))


print(*chunked_even([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 3))
