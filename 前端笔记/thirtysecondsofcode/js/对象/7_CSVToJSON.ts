// 将逗号分隔值字符串转换为对象的2d 数组。
console.log(CSVToJSON('col1,col2\na,b\nc,d'))
// [{'col1': 'a', 'col2': 'b'}, {'col1': 'c', 'col2': 'd'}];
console.log(CSVToJSON('col1;col2\na;b\nc;d', ';'))
// [{'col1': 'a', 'col2': 'b'}, {'col1': 'c', 'col2': 'd'}];

function CSVToJSON(csv: string, delimiter = ','): Record<PropertyKey, any>[] {
  const titles = csv.slice(0, csv.indexOf('\n')).split(delimiter)
  return csv
    .slice(csv.indexOf('\n') + 1)
    .split('\n')
    .map(row => {
      const vals = row.split(delimiter)
      return titles.reduce<Record<PropertyKey, any>>((obj, title, index) => {
        obj[title] = vals[index]
        return obj
      }, {})
    })
}

export {}
