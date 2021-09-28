class PhoneDirectory:
    def __init__(self, maxNumbers: int):
        """
        Initialize your data structure here
        @param maxNumbers - The maximum numbers that can be stored in the phone directory.
        """
        self.exist = set(range(maxNumbers))

    def get(self) -> int:
        """
        Provide a number which is not assigned to anyone.
        @return - Return an available number. Return -1 if none is available.
        """
        if not self.exist:
            return -1
        return self.exist.pop()

    def check(self, number: int) -> bool:
        """
        Check if a number is available or not.
        """
        return number in self.exist

    def release(self, number: int) -> None:
        """
        Recycle or release a number.
        """
        if number not in self.exist:
            self.exist.add(number)


print(set([1, 2, 3]).pop())

