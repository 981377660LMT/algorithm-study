# https://ssrs-cp.github.io/cp_library/other/poker_hands.hpp
# poker役判定
# 德州扑克比大小
# https://zh.wikipedia.org/wiki/%E6%92%B2%E5%85%8B%E7%89%8C%E5%9E%8B


import enum
from typing import Any, Literal, Tuple

Suit = Literal["S", "H", "D", "C"]
"""S: spade(黑桃♠) > H: heart(红桃♥) > D: diamond(梅花♣) > C: club(方块♦) """

Rank = Literal[2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14]
"""2-10对应2-10,11对应J,12对应Q,13对应K,14对应A"""


class Card:

    _MAPPING = {"S": 3, "H": 2, "D": 1, "C": 0}

    @staticmethod
    def fromStr(s: str) -> "Card":
        """https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2535"""
        assert len(s) == 2
        suit, rank = s[0], 0
        if s[1] == "A":
            rank = 14
        elif s[1] == "K":
            rank = 13
        elif s[1] == "Q":
            rank = 12
        elif s[1] == "J":
            rank = 11
        elif s[1] == "T":
            rank = 10
        else:
            rank = int(s[1])
        res = Card(suit=suit, rank=rank)  # type: ignore
        return res

    __slots__ = "suit", "rank"

    def __init__(self, suit: "Suit", rank: "Rank"):
        self.suit = suit
        self.rank = rank

    def __lt__(self, other: "Card") -> bool:
        return self.rank < other.rank or (self.rank == other.rank and self.suit < other.suit)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, Card):
            return False
        return self.rank == other.rank and self.suit == other.suit

    def __hash__(self) -> int:
        return self.rank * 4 + Card._MAPPING[self.suit]

    def __repr__(self) -> str:
        return f"Card({self.rank}, {self.suit})"


class PokerHand(enum.Enum):
    HIGH_CARD = 0
    """高牌/散牌,无对子."""
    ONE_PAIR = 1
    TWO_PAIR = 2

    THREE_OF_A_KIND = 3
    """三条,三张相同点数的牌."""

    STRAIGHT = 4
    """顺子,五张连续点数的牌."""

    FLUSH = 5
    """同花,五张同花色的牌."""

    FULL_HOUSE = 6
    """葫芦,三条加一对."""

    FOUR_OF_A_KIND = 7
    """铁支,四张相同点数的牌. """

    STRAIGHT_FLUSH = 8
    """同花顺,五张连续点数的同花色牌."""

    ROYAL_STRAIGHT_FLUSH = 9
    """皇家同花顺,同花顺中最大的牌(10,J,Q,K,A)."""


Hand = Tuple["Card", "Card", "Card", "Card", "Card"]


def evalHand(hand: "Hand") -> "Tuple[Any, ...]":
    """
    评价手牌,返回一个衡量手牌大小的元组.
    手牌的大小通过`返回的元组的字典序`来比较.
    """
    C = sorted(hand)
    isFlush = C[0].suit == C[1].suit == C[2].suit == C[3].suit == C[4].suit
    if isFlush and C[4].rank == 14 and C[0].rank == 10:
        return (PokerHand.ROYAL_STRAIGHT_FLUSH,)
    if isFlush and C[4].rank - C[0].rank == 4:
        return (PokerHand.STRAIGHT_FLUSH, C[4].rank)
    if isFlush and C[3].rank == 5 and C[4].rank == 14:
        return (PokerHand.STRAIGHT_FLUSH, 5)
    if C[0].rank == C[3].rank:
        return (PokerHand.FOUR_OF_A_KIND, C[0].rank, C[4].rank)
    if C[1].rank == C[4].rank:
        return (PokerHand.FOUR_OF_A_KIND, C[1].rank, C[0].rank)
    if C[0].rank == C[2].rank and C[3].rank == C[4].rank:
        return (PokerHand.FULL_HOUSE, C[0].rank, C[3].rank)
    if C[2].rank == C[4].rank and C[0].rank == C[1].rank:
        return (PokerHand.FULL_HOUSE, C[2].rank, C[0].rank)
    if isFlush:
        return (PokerHand.FLUSH, C[4].rank, C[3].rank, C[2].rank, C[1].rank, C[0].rank)
    if (
        C[1].rank - C[0].rank == 1
        and C[2].rank - C[1].rank == 1
        and C[3].rank - C[2].rank == 1
        and C[4].rank - C[3].rank == 1
    ):
        return (PokerHand.STRAIGHT, C[4].rank)
    if C[0].rank == 2 and C[1].rank == 3 and C[2].rank == 4 and C[3].rank == 5 and C[4].rank == 14:
        return (PokerHand.STRAIGHT, 5)
    if C[0].rank == C[2].rank:
        return (PokerHand.THREE_OF_A_KIND, C[0].rank, C[4].rank, C[3].rank)
    if C[1].rank == C[3].rank:
        return (PokerHand.THREE_OF_A_KIND, C[1].rank, C[4].rank, C[0].rank)
    if C[2].rank == C[4].rank:
        return (PokerHand.THREE_OF_A_KIND, C[2].rank, C[1].rank, C[0].rank)
    if C[0].rank == C[1].rank and C[2].rank == C[3].rank:
        return (PokerHand.TWO_PAIR, C[2].rank, C[0].rank, C[4].rank)
    if C[0].rank == C[1].rank and C[3].rank == C[4].rank:
        return (PokerHand.TWO_PAIR, C[3].rank, C[0].rank, C[2].rank)
    if C[1].rank == C[2].rank and C[3].rank == C[4].rank:
        return (PokerHand.TWO_PAIR, C[3].rank, C[1].rank, C[0].rank)
    if C[0].rank == C[1].rank:
        return (PokerHand.ONE_PAIR, C[0].rank, C[4].rank, C[3].rank, C[2].rank)
    if C[1].rank == C[2].rank:
        return (PokerHand.ONE_PAIR, C[1].rank, C[4].rank, C[3].rank, C[0].rank)
    if C[2].rank == C[3].rank:
        return (PokerHand.ONE_PAIR, C[2].rank, C[4].rank, C[1].rank, C[0].rank)
    if C[3].rank == C[4].rank:
        return (PokerHand.ONE_PAIR, C[3].rank, C[2].rank, C[1].rank, C[0].rank)
    return (PokerHand.HIGH_CARD, C[4].rank, C[3].rank, C[2].rank, C[1].rank, C[0].rank)


if __name__ == "__main__":
    # 你的任务是编写一个程序，计算你赢得游戏的概率，假设回合和河流是一致随机从剩余的卡片选择。
    # 你和对手总是要选择最强的那只手。联系应包括在计算中，即应计为损失。
    hands = tuple([Card("S", 2), Card("S", 3), Card("S", 4), Card("S", 5), Card("S", 6)])
    print(evalHand(hands))
    # https://ssrs-cp.github.io/cp_library/test/aoj/other/2535.test.cpp
