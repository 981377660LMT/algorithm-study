frozenSet([1, 2, 3, 1, 2])
// 将新创建对象的添加、删除和清除方法设置为未定义，这样就不能使用这些方法，实际上会冻结对象。
// Set { 1, 2, 3, add: undefined, delete: undefined, clear: undefined }

function frozenSet(iterable) {
  const set = new Set(iterable)
  set.add = undefined
  set.delete = undefined
  set.clear = undefined
  return set
}
