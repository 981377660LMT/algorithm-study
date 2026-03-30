#!/usr/bin/env python3
import argparse
import pathlib
import subprocess
import sys
from typing import Dict, List, Tuple

INVALID_CHARS = {
    ":": " -",
    '"': "",
    "|": "_",
    "?": "",
    "*": "_",
    "<": "(",
    ">": ")",
}

RESERVED = {
    "CON",
    "PRN",
    "AUX",
    "NUL",
    "COM1",
    "COM2",
    "COM3",
    "COM4",
    "COM5",
    "COM6",
    "COM7",
    "COM8",
    "COM9",
    "LPT1",
    "LPT2",
    "LPT3",
    "LPT4",
    "LPT5",
    "LPT6",
    "LPT7",
    "LPT8",
    "LPT9",
}


def run(cmd: List[str]) -> str:
    proc = subprocess.run(
        cmd,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    if proc.returncode != 0:
        raise RuntimeError(f"Command failed: {' '.join(cmd)}\n{proc.stderr.strip()}")
    return proc.stdout


def git_root() -> pathlib.Path:
    out = run(["git", "rev-parse", "--show-toplevel"]).strip()
    return pathlib.Path(out)


def is_dirty() -> bool:
    out = run(["git", "status", "--porcelain"])
    return bool(out.strip())


def list_paths() -> List[str]:
    out = run(["git", "-c", "core.quotepath=false", "ls-tree", "-r", "--name-only", "HEAD"])
    return [line for line in out.splitlines() if line]


def has_invalid_segment(seg: str) -> bool:
    if seg in ("", ".", ".."):
        return False
    if any(ch in seg for ch in INVALID_CHARS):
        return True
    if seg.endswith(" ") or seg.endswith("."):
        return True
    stem = seg.split(".")[0].upper()
    if stem in RESERVED:
        return True
    return False


def sanitize_segment(seg: str) -> str:
    s = seg
    for old, new in INVALID_CHARS.items():
        s = s.replace(old, new)
    s = s.rstrip(" .")
    if not s:
        s = "_"

    stem, dot, ext = s.partition(".")
    if stem.upper() in RESERVED:
        stem = stem + "_"
    s = stem + (dot + ext if dot else "")
    return s


def sanitize_path(p: str) -> str:
    parts = p.split("/")
    new_parts = [sanitize_segment(seg) for seg in parts]
    return "/".join(new_parts)


def resolve_collisions(pairs: List[Tuple[str, str]]) -> List[Tuple[str, str]]:
    used = set()
    out: List[Tuple[str, str]] = []

    for old, new in pairs:
        if new not in used:
            used.add(new)
            out.append((old, new))
            continue

        candidate = new
        path = pathlib.PurePosixPath(candidate)
        stem = path.stem
        suffix = "".join(path.suffixes)
        parent = str(path.parent)
        if parent == ".":
            parent = ""

        i = 2
        while True:
            name = f"{stem} ({i}){suffix}"
            cand = f"{parent}/{name}" if parent else name
            if cand not in used:
                used.add(cand)
                out.append((old, cand))
                break
            i += 1

    return out


def build_mapping(paths: List[str]) -> List[Tuple[str, str]]:
    raw: List[Tuple[str, str]] = []
    for p in paths:
        parts = p.split("/")
        if any(has_invalid_segment(seg) for seg in parts):
            raw.append((p, sanitize_path(p)))

    return resolve_collisions(raw)


def apply_git_mv(mapping: List[Tuple[str, str]]) -> None:
    for old, new in mapping:
        run(["git", "mv", old, new])


def write_output(mapping: List[Tuple[str, str]], output: str) -> None:
    lines = [f"{old}\t{new}" for old, new in mapping]
    data = "\n".join(lines) + ("\n" if lines else "")
    pathlib.Path(output).write_text(data, encoding="utf-8")


def main() -> int:
    parser = argparse.ArgumentParser(description="Sanitize Windows-incompatible Git paths.")
    parser.add_argument("--mode", choices=["preview", "apply"], default="preview")
    parser.add_argument("--output", default="")
    args = parser.parse_args()

    _ = git_root()
    paths = list_paths()
    mapping = build_mapping(paths)

    if args.output:
        write_output(mapping, args.output)

    if not mapping:
        print("No invalid paths found.")
        return 0

    print(f"Found {len(mapping)} paths to sanitize.")
    for old, new in mapping[:50]:
        print(f"{old} -> {new}")
    if len(mapping) > 50:
        print(f"... ({len(mapping) - 50} more)")

    if args.mode == "preview":
        return 0

    if is_dirty():
        print("Working tree is dirty. Commit or stash changes before apply mode.", file=sys.stderr)
        return 2

    apply_git_mv(mapping)
    print("Applied renames with git mv. Review and commit changes.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
