"""
An object is hashable if it has a hash value which never changes during its lifetime (it needs a __hash__() method), 
and can be compared to other objects (it needs an __eq__() method). 
Hashable objects which compare equal must have the same hash value.
"""
# https://docs.python.org/3/glossary.html#term-hashable

# !User defined classes have __hash__ by default that calls id(self)


from typing import Hashable


class Apple:
    def __init__(self, weight: int):
        self.weight = weight

    def __repr__(self):
        return f"Apple({self.weight})"


apple_a = Apple(1)
apple_b = Apple(1)
apple_c = Apple(2)

apple_dictionary = {apple_a: 3, apple_b: 4, apple_c: 5}

print(apple_dictionary[apple_a])  # 3
print(apple_dictionary)  # {Apple(1): 3, Apple(1): 4, Apple(2): 5}


class H:
    def __init__(self, value: int):
        self.value = value

    def __repr__(self):
        return f"H({self.value})"

    def __hash__(self):
        return hash(self.value)

    def __eq__(self, other):
        return self.value == other.value
