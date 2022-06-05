class TextEditor:
    def __init__(self):
        self.stack1 = []
        self.stack2 = []

    def addText(self, text: str) -> None:
        ...

    def deleteText(self, k: int) -> int:
        ...

    def cursorLeft(self, k: int) -> str:
        ...

    def cursorRight(self, k: int) -> str:
        ...
