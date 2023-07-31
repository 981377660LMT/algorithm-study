class F:
    def __init__(self) -> None:
        print("F.__init__")
        self._build()

    def _build(self) -> None:
        print("F._build")


class S(F):
    def __init__(self) -> None:
        print("S.__init__")
        super().__init__()

    def _build(self) -> None:
        print("S._build")


S()
