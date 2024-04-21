# from titan_pylib.data_structures.dict.hash_dict import HashDict
import random
from typing import List, Iterator, Tuple, Any

random.seed(0)
_titan_pylib_HashDict_K: int = 0x517CC1B727220A95


class HashDict:
    """ハッシュテーブルです。
    組み込み辞書の ``dict`` よりやや遅いです。
    """

    def __init__(self, e: int = -1, default: Any = 0, reserve: int = -1):
        """
        Args:
          e (int, optional): ``int`` 型で ``key`` として使用しない値です。
                              ``key`` を ``int`` 型以外のもので指定したいときは ``_hash(key) -> int`` 関数をいじってください。
          default (Any, optional): 存在しないキーにアクセスしたときの値です。
        """
        # e: keyとして使わない値
        # default: valのdefault値
        self._keys: List[int] = [e]
        self._vals: List[Any] = [default]
        self._msk: int = 0
        self._xor: int = random.getrandbits(1)
        if reserve > 0:
            self._keys: List[int] = [e] * (1 << (reserve.bit_length()))
            self._vals: List[Any] = [default] * (1 << (reserve.bit_length()))
            self._msk = (1 << (len(self._keys) - 1).bit_length()) - 1
            self._xor = random.getrandbits((len(self._keys) - 1).bit_length())
        self._e: int = e
        self._len: int = 0
        self._default: Any = default

    def _rebuild(self) -> None:
        old_keys, old_vals, _e = self._keys, self._vals, self._e
        self._keys = [_e] * (2 * len(old_keys))
        self._vals = [self._default] * len(self._keys)
        self._len = 0
        self._msk = (1 << (len(self._keys) - 1).bit_length()) - 1
        self._xor = random.getrandbits((len(self._keys) - 1).bit_length())
        for i in range(len(old_keys)):
            if old_keys[i] != _e:
                self.set(old_keys[i], old_vals[i])

    def _hash(self, key: int) -> int:
        return (
            ((((key >> 32) & self._msk) ^ (key & self._msk) ^ self._xor))
            * (_titan_pylib_HashDict_K & self._msk)
        ) & self._msk

    def get(self, key: int, default: Any = None) -> Any:
        """
        キーが ``key`` の値を返します。
        存在しない場合、引数 ``default`` に ``None`` 以外を指定した場合は ``default`` が、
        そうでない場合はコンストラクタで設定した ``default`` が返ります。

        期待 :math:`O(1)` です。
        """
        assert (
            key != self._e
        ), f"KeyError: HashDict.get({key}, {default}), {key} cannot be equal to {self._e}"
        l, _keys, _e = len(self._keys), self._keys, self._e
        h = self._hash(key)
        while True:
            x = _keys[h]
            if x == _e:
                return self._vals[h] if default is None else default
            if x == key:
                return self._vals[h]
            h = 0 if h == l - 1 else h + 1

    def add(self, key: int, val: Any, default: Any) -> None:
        assert (
            key != self._e
        ), f"KeyError: HashDict.add({key}, {default}), {key} cannot be equal to {self._e}"
        l, _keys, _e = len(self._keys), self._keys, self._e
        h = self._hash(key)
        while True:
            x = _keys[h]
            if x == _e:
                self._vals[h] = val
                return
            if x == key:
                self._vals[h] += val
                return
            h = 0 if h == l - 1 else h + 1

    def set(self, key: int, val: Any) -> None:
        """キーを ``key`` として ``val`` を格納します。
        ``key`` が既に存在している場合は上書きされます。

        期待 :math:`O(1)` です。
        """
        assert (
            key != self._e
        ), f"KeyError: HashDict.set({key}, {val}), {key} cannot be equal to {self._e}"
        l, _keys, _e = len(self._keys), self._keys, self._e
        l -= 1
        h = self._hash(key)
        while True:
            x = _keys[h]
            if x == _e:
                _keys[h] = key
                self._vals[h] = val
                self._len += 1
                if 2 * self._len > len(self._keys):
                    self._rebuild()
                return
            if x == key:
                self._vals[h] = val
                return
            h = 0 if h == l else h + 1

    def __contains__(self, key: int) -> bool:
        """存在判定です。

        期待 :math:`O(1)` です。

        Returns:
          bool: ``key`` が存在すれば ``True`` を、そうでなければ ``False`` を返します。
        """
        assert key != self._e, f"KeyError: {key} in HashDict, {key} cannot be equal to {self._e}"
        l, _keys, _e = len(self._keys), self._keys, self._e
        h = self._hash(key)
        while True:
            x = _keys[h]
            if x == _e:
                return False
            if x == key:
                return True
            h += 1
            if h == l:
                h = 0

    __getitem__ = get
    __setitem__ = set

    def keys(self) -> Iterator[int]:
        """``key 集合`` を列挙するイテレータです。"""
        _keys, _e = self._keys, self._e
        for i in range(len(_keys)):
            if _keys[i] != _e:
                yield _keys[i]

    def values(self) -> Iterator[Any]:
        """``val 集合`` を列挙するイテレータです。"""
        _keys, _vals, _e = self._keys, self._vals, self._e
        for i in range(len(_keys)):
            if _keys[i] != _e:
                yield _vals[i]

    def items(self) -> Iterator[Tuple[int, Any]]:
        """``key とそれに対応する val のタプル`` を列挙するイテレータです。"""
        _keys, _vals, _e = self._keys, self._vals, self._e
        for i in range(len(_keys)):
            if _keys[i] != _e:
                yield _keys[i], _vals[i]

    def __str__(self):
        return "{" + ", ".join(map(lambda x: f"{x[0]}: {x[1]}", self.items())) + "}"

    def __len__(self):
        return self._len
