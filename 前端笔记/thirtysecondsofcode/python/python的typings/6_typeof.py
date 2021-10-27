from typing import Type


a = 3  # Has type 'int'
b = int  # Has type 'Type[int]'
c = type(a)  # Also has type 'Type[int]'


class User:
    ...


class BasicUser(User):
    ...


class ProUser(User):
    ...


class TeamUser(User):
    ...


# Accepts User, BasicUser, ProUser, TeamUser, ...
def make_new_user(user_class: Type[User]) -> User:
    # ...
    return user_class()
