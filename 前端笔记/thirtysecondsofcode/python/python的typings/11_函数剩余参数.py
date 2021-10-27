def stars(*args: int, **kwargs: float) -> None:
    # 'args' has type 'Tuple[int, ...]' (a tuple of ints)
    # 'kwargs' has type 'Dict[str, float]' (a dict of strs to floats)
    for arg in args:
        print(arg)
    for key, value in kwargs:
        print(key, value)
