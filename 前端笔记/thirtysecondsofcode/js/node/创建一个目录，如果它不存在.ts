import { existsSync, mkdirSync, PathLike } from 'fs'

const createDirIfNotExists = (dir: PathLike) => (!existsSync(dir) ? mkdirSync(dir) : undefined)
createDirIfNotExists('test')
// creates the directory 'test', if it doesn't exist
