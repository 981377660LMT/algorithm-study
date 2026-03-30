import json
import sys
import pandas as pd
from pathlib import Path


def xlsx_to_json(src: str, dst: str = None):
    """
    xlsx 转 json
    - src: 源文件路径
    - dst: 目标文件路径（默认同目录 .json）
    """
    p = Path(src)
    if not p.exists():
        raise FileNotFoundError(f"文件不存在: {src}")

    # 使用 pandas 按字符串读取，保留所有数据完整性
    sheets = pd.read_excel(p, sheet_name=None, dtype=str, engine="openpyxl")

    result = {}
    for name, df in sheets.items():
        # NaN → 空字符串
        df = df.fillna("")
        result[name] = df.to_dict(orient="records")

    # 单 sheet 直接输出数组，多 sheet 保留层级结构
    out = result[list(result.keys())[0]] if len(result) == 1 else result

    if dst is None:
        dst = str(p.with_suffix(".json"))

    with open(dst, "w", encoding="utf-8") as f:
        json.dump(out, f, ensure_ascii=False, indent=2)

    cnt = sum(len(v) if isinstance(v, list) else 0 for v in result.values())
    print(f"✓ 转换完成: {dst}")
    print(f"✓ 共 {cnt} 条数据")
    return out


def main():
    if len(sys.argv) < 2:
        print("用法: python xlsx_to_json.py <xlsx文件路径> [输出json路径]")
        print("示例: python xlsx_to_json.py data.xlsx output.json")
        sys.exit(1)

    src = sys.argv[1]
    dst = sys.argv[2] if len(sys.argv) > 2 else None

    try:
        xlsx_to_json(src, dst)
    except Exception as e:
        print(f"✗ 错误: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
