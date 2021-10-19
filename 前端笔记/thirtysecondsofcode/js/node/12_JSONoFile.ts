import { writeFileSync } from 'fs'

const JSONToFile = (obj: object, filename: string) =>
  writeFileSync(`${filename}.json`, JSON.stringify(obj, null, 2))
JSONToFile({ test: 'is passed' }, 'testJsonFile')
// writes the object to 'testJsonFile.json'
