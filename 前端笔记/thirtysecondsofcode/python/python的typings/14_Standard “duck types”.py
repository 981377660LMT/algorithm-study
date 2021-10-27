from typing import Mapping, MutableMapping, Sequence, Iterable, List, Set


# Use Iterable for generic iterables (anything usable in "for"),
# and Sequence where a sequence (supporting "len" and "__getitem__") is
# required
def f(ints: Iterable[int]) -> List[str]:
    return [str(x) for x in ints]


f(range(1, 3))


# Mapping describes a dict-like object (with "__getitem__") that we won't
# mutate, and MutableMapping one (with "__setitem__") that we might
def f2(my_mapping: Mapping[int, str]) -> List[int]:
    my_mapping[5] = 'maybe'  # if we try this, mypy will throw an error...
    return list(my_mapping.keys())


f2({3: 'yes', 4: 'no'})


def f3(my_mapping: MutableMapping[int, str]) -> Set[str]:
    my_mapping[5] = 'maybe'  # ...but mypy is OK with this.
    return set(my_mapping.values())


f3({3: 'yes', 4: 'no'})
