// 将逗号分隔值字符串转换为对象的2d 数组。
CSVToArray('a,b\nc,d') // [['a', 'b'], ['c', 'd']];
CSVToArray('a;b\nc;d', ';') // [['a', 'b'], ['c', 'd']];
CSVToArray('col1,col2\na,b\nc,d', ',', true) // [['a', 'b'], ['c', 'd']];

function CSVToArray(csv: string, delimiter = ',', omitFirstRow = false): string[][] {
  return (
    csv
      // 删除第一行(标题行)。
      .slice(omitFirstRow ? csv.indexOf('\n') + 1 : 0)
      .split('\n')
      .map(v => v.split(delimiter))
  )
}

export {}
