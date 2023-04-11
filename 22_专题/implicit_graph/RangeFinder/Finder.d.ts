interface PrevFinder {
  prev(x: number): number | null
  erase(x: number): void
}

interface NextFinder {
  next(x: number): number | null
  erase(x: number): void
}

interface Finder extends PrevFinder, NextFinder {}
