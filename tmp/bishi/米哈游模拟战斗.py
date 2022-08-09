# OOP设计的关键:
# 0. 按照流程提取需要哪些类、设计什么方法和属性
# 1. `敌人`要加到自己属性上，而不是在每个方法里传入(依赖注入)
# 2. 如果构造函数里不好传某个属性，可以用建造者模式的思想，用addxx方法添加属性
# 例如addEnemy、addEdge
# 3. 通过一个公共的run() 方法来执行一回合中的逻辑会更好
# 4. 注意：抽象基类上需要包括状态很多时怎么办?
# 可以用一个哈希表/数组来存储多个`状态`
# `状态`实现某个接口
# 5. 用property装饰器来实现计算属性(computed)


# 当时写的答案：
from abc import ABCMeta, abstractmethod
from math import floor
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

WIN1 = "I love V2V forever!"
WIN2 = "Kalpas yame te!"
T = int(input())


class AbstractRole(metaclass=ABCMeta):
    """战斗中角色的抽象基类"""

    __slots__ = (
        "_hp",
        "_attack",
        "_defend",
        "_speed",
        "_turn",
        "_isCrazy",
        "_isSleep",
        "_hasSub",
        "_enemy",
    )

    def __init__(self, hp: int, attack: int, defend: int, speed: int) -> None:
        self._hp = hp
        self._attack = attack
        self._defend = defend
        self._speed = speed

        self._turn = 1
        self._isCrazy = False
        self._isSleep = False
        self._hasSub = True

    def addEnemy(self, enemy: "AbstractRole") -> None:
        self._enemy = enemy

    def endTurn(self) -> None:
        self._turn += 1

    def receiveAttack(self, attack: int, type: str) -> None:
        if type == "normal":
            realDamage = max(0, attack - self._defend)
            self._hp -= realDamage
        else:
            realDamage = attack
            self._hp -= realDamage

    def startTurn(self) -> None:
        if self._hasSub:
            self.excuteSubSkill()

    def attack(self) -> None:
        if self._turn % 3 == 0:
            damage = self._attack
            if self._isCrazy:
                self.receiveAttack(damage, "normal")
                self._isCrazy = False
            else:
                self._enemy.receiveAttack(damage, "normal")

        else:
            self.excuteMainSkill()

    @abstractmethod
    def excuteMainSkill(self) -> None:
        ...

    @abstractmethod
    def excuteSubSkill(self) -> None:
        ...

    @property
    def hasDied(self) -> bool:
        return self._hp <= 0

    @property
    def speed(self) -> int:
        return self._speed


class V2V(AbstractRole):
    def __init__(self, hp: int, attack: int, defend: int, speed: int) -> None:
        super().__init__(hp, attack, defend, speed)

    def excuteMainSkill(self) -> None:
        if self._turn % 3 != 0:
            return
        damage = self._attack
        self._enemy.receiveAttack(damage, "normal")  # 攻击对手并混乱
        self._enemy._isCrazy = True  # 对手下次普通攻击打自己

    def excuteSubSkill(self) -> None:
        if self._hp < 31:
            self._hp += 20
            self._enemy._hp += 20
            self._attack += 15
            self._hasSub = False


class Kalpas(AbstractRole):
    def __init__(self, hp: int, attack: int, defend: int, speed: int) -> None:
        super().__init__(hp, attack, defend, speed)

    def excuteMainSkill(self) -> None:
        if self._turn % 3 != 0:
            return
        if self._hp < 11:
            damage = self._attack
            self._enemy.receiveAttack(damage, "normal")
        else:
            self._hp -= 10
            self._enemy.receiveAttack(45, "normal")
            self._enemy.receiveAttack(20, "magic")
            self._isSleep = True

    def excuteSubSkill(self) -> None:
        delta = max(0, floor((100 - self._hp) / 5))
        self._attack += delta  # 这里有问题 不能每个回合都叠加


####################################################################
for _ in range(T):
    hp1, att1, def1, spd1 = list(map(int, input().split()))
    hp2, att2, def2, spd2 = list(map(int, input().split()))
    v2v, kalpas = V2V(hp1, att1, def1, spd1), Kalpas(hp2, att2, def2, spd2)
    v2v.addEnemy(kalpas)
    kalpas.addEnemy(v2v)

    while True:
        # 速度相同时 千劫先动
        role1, role2 = kalpas, v2v
        if role2.speed > role1.speed:
            role1, role2 = role2, role1

        if role1._isSleep:
            role1._isSleep = False
            role1._turn += 1
        else:
            role1.startTurn()
            role1.attack()
            role1.endTurn()
        if role2.hasDied:
            print(WIN1)
            break

        if role2._isSleep:
            role2._isSleep = False
            role2._turn += 1
        else:
            role2.startTurn()
            role2.attack()
            role2.endTurn()
        if role1.hasDied:
            print(WIN2)
            break
