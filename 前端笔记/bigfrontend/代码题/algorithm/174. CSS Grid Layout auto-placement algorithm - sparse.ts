type Grid = number[][]
type Item = {
  id: number
  style?: {
    gridRowStart?: number | string
    gridRowEnd?: number | string
    gridColumnStart?: number | string
    gridColumnEnd?: number | string
  }
}

// Are you familiar with CSS Grid Layout ?
// grid-auto-flow controls how the auto-placement algorithm works.
// Default packing mode is sparse, please create a function to illustrate how it works.
function layout(rows: number, columns: number, items: Array<Item>): Grid {
  // your code here
  const res: Grid = Array.from({ length: rows }, () => Array(columns).fill(0))
  return res
}

// TODO
