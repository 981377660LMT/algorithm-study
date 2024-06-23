type RowId = number
type ColId = number
type CellMatrix = { rowIds: RowId[]; colIds: ColId[] }
interface IDiff {
  removed: CellMatrix
  added: CellMatrix
  colSortChanged: boolean
  rowSortChanged: boolean
  meta?: {
    tooManyChanges?: boolean
  }
}
