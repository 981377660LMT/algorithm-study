console.log(JSONtoCSV([{ a: 1, b: 2 }, { a: 3, b: 4, c: 5 }, { a: 6 }, { b: 7 }], ['a', 'b'])) // 'a,b\n"1","2"\n"3","4"\n"6",""\n"","7"'
console.log(JSONtoCSV([{ a: 1, b: 2 }, { a: 3, b: 4, c: 5 }, { a: 6 }, { b: 7 }], ['a', 'b'], ';')) // 'a,b\n"1","2"\n"3","4"\n"6",""\n"","7"'

function JSONtoCSV(
  arr: Record<string, number>[],
  columns: [...titles: string[]],
  delimiter = ','
): string {
  return [
    columns.join(delimiter),
    ...arr.map(obj =>
      columns.reduce((pre, key) => `${pre}${pre.length ? delimiter : ''}"${obj[key] ?? ''}"`, '')
    ),
  ].join('\n')
}

export {}
