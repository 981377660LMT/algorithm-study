# from titan_pylib.algorithm.permutation import Permutation

from typing import List, Any


class Permutation:
    """順列のライブラリです。"""

    @staticmethod
    def next_permutation(a: List[Any], l: int = 0, r: int = -1) -> bool:
        """列 ``a[l, r)`` を辞書順で次の順列にします。

        Args:
          a (List[Any])
          l (int, optional)
          r (int, optional)

        Returns:
          bool: 辞書順で次の順列が存在する場合は ``True`` 、存在しない場合は ``False`` を返します。
        """
        if r == -1:
            r = len(a)
        for i in range(r - 2, l - 1, -1):
            if a[i] < a[i + 1]:
                for j in range(r - 1, i, -1):
                    if a[i] < a[j]:
                        a[i], a[j] = a[j], a[i]
                        p = i + 1
                        q = r - 1
                        while p < q:
                            a[p], a[q] = a[q], a[p]
                            p += 1
                            q -= 1
                        return True
        return False

    @staticmethod
    def prev_permutation(a: List[Any]) -> bool:
        """列 ``a`` を辞書順で後の順列にします。

        Args:
          a (List[Any])

        Returns:
          bool: 辞書順で後の順列が存在する場合は ``True`` 、存在しない場合は ``False`` を返します。
        """
        l = 0
        r = len(a)
        for i in range(r - 2, l - 1, -1):
            if a[i] > a[i + 1]:
                for j in range(r - 1, i, -1):
                    if a[i] > a[j]:
                        a[i], a[j] = a[j], a[i]
                        p = i + 1
                        q = r - 1
                        while p < q:
                            a[p], a[q] = a[q], a[p]
                            p += 1
                            q -= 1
                        return True
        return False
