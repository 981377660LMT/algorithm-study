# from titan_pylib.io.fast_o import FastO
import os
import io


class FastO:
    """標準出力高速化ライブラリです。"""

    _output = io.StringIO()

    @classmethod
    def write(cls, *args, sep: str = " ", end: str = "\n", flush: bool = False) -> None:
        """標準出力します。次の ``FastO.flush()`` が起きると print します。"""
        wr = cls._output.write
        for i in range(len(args) - 1):
            wr(str(args[i]))
            wr(sep)
        if args:
            wr(str(args[-1]))
        wr(end)
        if flush:
            cls.flush()

    @classmethod
    def flush(cls) -> None:
        """``flush`` します。これを実行しないと ``write`` した内容が表示されないので忘れないでください。"""
        os.write(1, cls._output.getvalue().encode())
        cls._output.close()
