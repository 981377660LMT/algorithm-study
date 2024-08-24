// @ts-nocheck

import BitSet from './BitSet.ts'
import Chunk from './Chunk.ts'
import SourceMap from './SourceMap.ts'
import guessIndent from './utils/guessIndent.ts'
import getRelativePath from './utils/getRelativePath.ts'
import isObject from './utils/isObject.ts'
import getLocator from './utils/getLocator.ts'
import Mappings from './utils/Mappings.ts'

const n = '\n'

const warned = {
  insertLeft: false,
  insertRight: false,
  storeName: false
}

export default class MagicString {
  constructor(string, options = {}) {
    // 创建一个初始的 Chunk 实例，表示整个源代码字符串
    const chunk = new Chunk(0, string.length, string)

    Object.defineProperties(this, {
      original: { writable: true, value: string },
      outro: { writable: true, value: '' }, // !结束字符串（将在原始字符串后追加的内容
      intro: { writable: true, value: '' }, // !开始字符串（将在原始字符串前追加的内容）
      firstChunk: { writable: true, value: chunk },
      lastChunk: { writable: true, value: chunk },
      lastSearchedChunk: { writable: true, value: chunk },
      byStart: { writable: true, value: {} }, // !根据起始位置索引的 Chunk 实例的映射
      byEnd: { writable: true, value: {} }, // !根据结束位置索引的 Chunk 实例的映射
      filename: { writable: true, value: options.filename },
      indentExclusionRanges: { writable: true, value: options.indentExclusionRanges },
      sourcemapLocations: { writable: true, value: new BitSet() }, // !sourcemap的位置信息
      storedNames: { writable: true, value: {} },
      indentStr: { writable: true, value: undefined },
      ignoreList: { writable: true, value: options.ignoreList }
    })

    this.byStart[0] = chunk
    this.byEnd[string.length] = chunk
  }

  addSourcemapLocation(char) {
    this.sourcemapLocations.add(char)
  }

  append(content) {
    if (typeof content !== 'string') throw new TypeError('outro content must be a string')

    this.outro += content
    return this
  }

  appendLeft(index, content) {
    if (typeof content !== 'string') throw new TypeError('inserted content must be a string')

    this._split(index)

    const chunk = this.byEnd[index]

    if (chunk) {
      chunk.appendLeft(content)
    } else {
      this.intro += content
    }

    return this
  }

  appendRight(index, content) {
    if (typeof content !== 'string') throw new TypeError('inserted content must be a string')

    this._split(index)

    const chunk = this.byStart[index]

    if (chunk) {
      chunk.appendRight(content)
    } else {
      this.outro += content
    }

    return this
  }

  clone() {
    const cloned = new MagicString(this.original, { filename: this.filename })

    let originalChunk = this.firstChunk
    let clonedChunk = (cloned.firstChunk = cloned.lastSearchedChunk = originalChunk.clone())

    while (originalChunk) {
      cloned.byStart[clonedChunk.start] = clonedChunk
      cloned.byEnd[clonedChunk.end] = clonedChunk

      const nextOriginalChunk = originalChunk.next
      const nextClonedChunk = nextOriginalChunk && nextOriginalChunk.clone()

      if (nextClonedChunk) {
        clonedChunk.next = nextClonedChunk
        nextClonedChunk.previous = clonedChunk

        clonedChunk = nextClonedChunk
      }

      originalChunk = nextOriginalChunk
    }

    cloned.lastChunk = clonedChunk

    if (this.indentExclusionRanges) {
      cloned.indentExclusionRanges = this.indentExclusionRanges.slice()
    }

    cloned.sourcemapLocations = new BitSet(this.sourcemapLocations)

    cloned.intro = this.intro
    cloned.outro = this.outro

    return cloned
  }

  generateDecodedMap(options) {
    options = options || {}

    const sourceIndex = 0
    const names = Object.keys(this.storedNames)
    const mappings = new Mappings(options.hires)

    const locate = getLocator(this.original)

    if (this.intro) {
      mappings.advance(this.intro)
    }

    this.firstChunk.eachNext(chunk => {
      const loc = locate(chunk.start)

      if (chunk.intro.length) mappings.advance(chunk.intro)

      if (chunk.edited) {
        mappings.addEdit(sourceIndex, chunk.content, loc, chunk.storeName ? names.indexOf(chunk.original) : -1)
      } else {
        mappings.addUneditedChunk(sourceIndex, chunk, this.original, loc, this.sourcemapLocations)
      }

      if (chunk.outro.length) mappings.advance(chunk.outro)
    })

    return {
      file: options.file ? options.file.split(/[/\\]/).pop() : undefined,
      sources: [options.source ? getRelativePath(options.file || '', options.source) : options.file || ''],
      sourcesContent: options.includeContent ? [this.original] : undefined,
      names,
      mappings: mappings.raw,
      x_google_ignoreList: this.ignoreList ? [sourceIndex] : undefined
    }
  }

  generateMap(options) {
    return new SourceMap(this.generateDecodedMap(options))
  }

  _ensureindentStr() {
    if (this.indentStr === undefined) {
      this.indentStr = guessIndent(this.original)
    }
  }

  _getRawIndentString() {
    this._ensureindentStr()
    return this.indentStr
  }

  getIndentString() {
    this._ensureindentStr()
    return this.indentStr === null ? '\t' : this.indentStr
  }

  indent(indentStr, options) {
    const pattern = /^[^\r\n]/gm

    if (isObject(indentStr)) {
      options = indentStr
      indentStr = undefined
    }

    if (indentStr === undefined) {
      this._ensureindentStr()
      indentStr = this.indentStr || '\t'
    }

    if (indentStr === '') return this // noop

    options = options || {}

    // Process exclusion ranges
    const isExcluded = {}

    if (options.exclude) {
      const exclusions = typeof options.exclude[0] === 'number' ? [options.exclude] : options.exclude
      exclusions.forEach(exclusion => {
        for (let i = exclusion[0]; i < exclusion[1]; i += 1) {
          isExcluded[i] = true
        }
      })
    }

    let shouldIndentNextCharacter = options.indentStart !== false
    const replacer = match => {
      if (shouldIndentNextCharacter) return `${indentStr}${match}`
      shouldIndentNextCharacter = true
      return match
    }

    this.intro = this.intro.replace(pattern, replacer)

    let charIndex = 0
    let chunk = this.firstChunk

    while (chunk) {
      const end = chunk.end

      if (chunk.edited) {
        if (!isExcluded[charIndex]) {
          chunk.content = chunk.content.replace(pattern, replacer)

          if (chunk.content.length) {
            shouldIndentNextCharacter = chunk.content[chunk.content.length - 1] === '\n'
          }
        }
      } else {
        charIndex = chunk.start

        while (charIndex < end) {
          if (!isExcluded[charIndex]) {
            const char = this.original[charIndex]

            if (char === '\n') {
              shouldIndentNextCharacter = true
            } else if (char !== '\r' && shouldIndentNextCharacter) {
              shouldIndentNextCharacter = false

              if (charIndex === chunk.start) {
                chunk.prependRight(indentStr)
              } else {
                this._splitChunk(chunk, charIndex)
                chunk = chunk.next
                chunk.prependRight(indentStr)
              }
            }
          }

          charIndex += 1
        }
      }

      charIndex = chunk.end
      chunk = chunk.next
    }

    this.outro = this.outro.replace(pattern, replacer)

    return this
  }

  insert() {
    throw new Error('magicString.insert(...) is deprecated. Use prependRight(...) or appendLeft(...)')
  }

  insertLeft(index, content) {
    if (!warned.insertLeft) {
      console.warn('magicString.insertLeft(...) is deprecated. Use magicString.appendLeft(...) instead') // eslint-disable-line no-console
      warned.insertLeft = true
    }

    return this.appendLeft(index, content)
  }

  insertRight(index, content) {
    if (!warned.insertRight) {
      console.warn('magicString.insertRight(...) is deprecated. Use magicString.prependRight(...) instead') // eslint-disable-line no-console
      warned.insertRight = true
    }

    return this.prependRight(index, content)
  }

  move(start, end, index) {
    if (index >= start && index <= end) throw new Error('Cannot move a selection inside itself')

    this._split(start)
    this._split(end)
    this._split(index)

    const first = this.byStart[start]
    const last = this.byEnd[end]

    const oldLeft = first.previous
    const oldRight = last.next

    const newRight = this.byStart[index]
    if (!newRight && last === this.lastChunk) return this
    const newLeft = newRight ? newRight.previous : this.lastChunk

    if (oldLeft) oldLeft.next = oldRight
    if (oldRight) oldRight.previous = oldLeft

    if (newLeft) newLeft.next = first
    if (newRight) newRight.previous = last

    if (!first.previous) this.firstChunk = last.next
    if (!last.next) {
      this.lastChunk = first.previous
      this.lastChunk.next = null
    }

    first.previous = newLeft
    last.next = newRight || null

    if (!newLeft) this.firstChunk = first
    if (!newRight) this.lastChunk = last

    return this
  }

  overwrite(start, end, content, options) {
    options = options || {}
    return this.update(start, end, content, { ...options, overwrite: !options.contentOnly })
  }

  // update 方法会首先查找起始索引和结束索引之间的所有 Chunk 实例，然后在这些 Chunk 实例上进行相应的操作
  update(start, end, content, options) {
    if (typeof content !== 'string') throw new TypeError('replacement content must be a string')

    if (this.original.length !== 0) {
      while (start < 0) start += this.original.length
      while (end < 0) end += this.original.length
    }

    if (end > this.original.length) throw new Error('end is out of bounds')
    if (start === end) throw new Error('Cannot overwrite a zero-length range – use appendLeft or prependRight instead')

    // 将替换范围拆分成多个 Chunk 实例
    this._split(start)
    this._split(end)

    if (options === true) {
      if (!warned.storeName) {
        console.warn('The final argument to magicString.overwrite(...) should be an options object. See https://github.com/rich-harris/magic-string') // eslint-disable-line no-console
        warned.storeName = true
      }

      options = { storeName: true }
    }
    const storeName = options !== undefined ? options.storeName : false
    const overwrite = options !== undefined ? options.overwrite : false

    if (storeName) {
      const original = this.original.slice(start, end)
      Object.defineProperty(this.storedNames, original, {
        writable: true,
        value: true,
        enumerable: true
      })
    }

    const first = this.byStart[start]
    const last = this.byEnd[end]

    // 如果存在first，说明在范围内发生了拆分，所以需要迭代所有的Chunk实例，将它们的内容置为空字符串，表示删除内容。然后，将第一个Chunk实例的内容更新为给定的content，使用给定的storeName和overwrite参数。
    // 如果没有first，则说明在范围的末尾插入了新的内容。在这种情况下，创建一个新的Chunk实例newChunk，表示要插入的内容。

    if (first) {
      // 如果存在第一个Chunk实例，则在范围内进行替换
      let chunk = first
      while (chunk !== last) {
        if (chunk.next !== this.byStart[chunk.end]) {
          throw new Error('Cannot overwrite across a split point')
        }
        chunk = chunk.next
        chunk.edit('', false)
      }

      // 在第一个Chunk实例上进行编辑
      first.edit(content, storeName, !overwrite)
    } else {
      // 如果没有第一个Chunk实例，则表示在范围之后追加内容
      // must be inserting at the end
      const newChunk = new Chunk(start, end, '').edit(content, storeName)

      // TODO last chunk in the array may not be the last chunk, if it's moved...
      last.next = newChunk
      newChunk.previous = last
    }

    return this
  }

  prepend(content) {
    if (typeof content !== 'string') throw new TypeError('outro content must be a string')

    this.intro = content + this.intro
    return this
  }

  prependLeft(index, content) {
    if (typeof content !== 'string') throw new TypeError('inserted content must be a string')

    this._split(index)

    const chunk = this.byEnd[index]

    if (chunk) {
      chunk.prependLeft(content)
    } else {
      this.intro = content + this.intro
    }

    return this
  }

  prependRight(index, content) {
    if (typeof content !== 'string') throw new TypeError('inserted content must be a string')

    this._split(index)

    const chunk = this.byStart[index]

    if (chunk) {
      chunk.prependRight(content)
    } else {
      this.outro = content + this.outro
    }

    return this
  }

  remove(start, end) {
    if (this.original.length !== 0) {
      while (start < 0) start += this.original.length
      while (end < 0) end += this.original.length
    }

    if (start === end) return this

    if (start < 0 || end > this.original.length) throw new Error('Character is out of bounds')
    if (start > end) throw new Error('end must be greater than start')

    this._split(start)
    this._split(end)

    let chunk = this.byStart[start]

    while (chunk) {
      chunk.intro = ''
      chunk.outro = ''
      chunk.edit('')

      chunk = end > chunk.end ? this.byStart[chunk.end] : null
    }

    return this
  }

  reset(start, end) {
    if (this.original.length !== 0) {
      while (start < 0) start += this.original.length
      while (end < 0) end += this.original.length
    }

    if (start === end) return this

    if (start < 0 || end > this.original.length) throw new Error('Character is out of bounds')
    if (start > end) throw new Error('end must be greater than start')

    this._split(start)
    this._split(end)

    let chunk = this.byStart[start]

    while (chunk) {
      chunk.reset()

      chunk = end > chunk.end ? this.byStart[chunk.end] : null
    }

    return this
  }

  lastChar() {
    if (this.outro.length) return this.outro[this.outro.length - 1]
    let chunk = this.lastChunk
    do {
      if (chunk.outro.length) return chunk.outro[chunk.outro.length - 1]
      if (chunk.content.length) return chunk.content[chunk.content.length - 1]
      if (chunk.intro.length) return chunk.intro[chunk.intro.length - 1]
    } while ((chunk = chunk.previous))
    if (this.intro.length) return this.intro[this.intro.length - 1]
    return ''
  }

  lastLine() {
    let lineIndex = this.outro.lastIndexOf(n)
    if (lineIndex !== -1) return this.outro.substr(lineIndex + 1)
    let lineStr = this.outro
    let chunk = this.lastChunk
    do {
      if (chunk.outro.length > 0) {
        lineIndex = chunk.outro.lastIndexOf(n)
        if (lineIndex !== -1) return chunk.outro.substr(lineIndex + 1) + lineStr
        lineStr = chunk.outro + lineStr
      }

      if (chunk.content.length > 0) {
        lineIndex = chunk.content.lastIndexOf(n)
        if (lineIndex !== -1) return chunk.content.substr(lineIndex + 1) + lineStr
        lineStr = chunk.content + lineStr
      }

      if (chunk.intro.length > 0) {
        lineIndex = chunk.intro.lastIndexOf(n)
        if (lineIndex !== -1) return chunk.intro.substr(lineIndex + 1) + lineStr
        lineStr = chunk.intro + lineStr
      }
    } while ((chunk = chunk.previous))
    lineIndex = this.intro.lastIndexOf(n)
    if (lineIndex !== -1) return this.intro.substr(lineIndex + 1) + lineStr
    return this.intro + lineStr
  }

  slice(start = 0, end = this.original.length) {
    if (this.original.length !== 0) {
      while (start < 0) start += this.original.length
      while (end < 0) end += this.original.length
    }

    let result = ''

    // find start chunk
    let chunk = this.firstChunk
    while (chunk && (chunk.start > start || chunk.end <= start)) {
      // found end chunk before start
      if (chunk.start < end && chunk.end >= end) {
        return result
      }

      chunk = chunk.next
    }

    if (chunk && chunk.edited && chunk.start !== start) throw new Error(`Cannot use replaced character ${start} as slice start anchor.`)

    const startChunk = chunk
    while (chunk) {
      if (chunk.intro && (startChunk !== chunk || chunk.start === start)) {
        result += chunk.intro
      }

      const containsEnd = chunk.start < end && chunk.end >= end
      if (containsEnd && chunk.edited && chunk.end !== end) throw new Error(`Cannot use replaced character ${end} as slice end anchor.`)

      const sliceStart = startChunk === chunk ? start - chunk.start : 0
      const sliceEnd = containsEnd ? chunk.content.length + end - chunk.end : chunk.content.length

      result += chunk.content.slice(sliceStart, sliceEnd)

      if (chunk.outro && (!containsEnd || chunk.end === end)) {
        result += chunk.outro
      }

      if (containsEnd) {
        break
      }

      chunk = chunk.next
    }

    return result
  }

  // TODO deprecate this? not really very useful
  snip(start, end) {
    const clone = this.clone()
    clone.remove(0, start)
    clone.remove(end, clone.original.length)

    return clone
  }

  _split(index) {
    // 如果指定索引位置的起始或结束处已经存在 Chunk 实例，直接返回
    if (this.byStart[index] || this.byEnd[index]) return

    // 初始化搜索方向和起始 Chunk 实例
    let chunk = this.lastSearchedChunk
    const searchForward = index > chunk.end

    // 在 Chunk 实例链表中查找包含指定索引位置的 Chunk 实例
    while (chunk) {
      // 如果当前 Chunk 实例包含指定索引位置，调用 _splitChunk 方法进行拆分
      if (chunk.contains(index)) return this._splitChunk(chunk, index)

      // 根据搜索方向更新当前 Chunk 实例
      chunk = searchForward ? this.byStart[chunk.end] : this.byEnd[chunk.start]
    }
  }

  _splitChunk(chunk, index) {
    // 如果已经编辑过的 chunk 不为空，且拆分的位置处有内容，抛出错误
    if (chunk.edited && chunk.content.length) {
      // zero-length edited chunks are a special case (overlapping replacements)
      const loc = getLocator(this.original)(index)
      throw new Error(`Cannot split a chunk that has already been edited (${loc.line}:${loc.column} – "${chunk.original}")`)
    }

    // 在指定索引位置拆分 chunk，并获取新的 chunk 实例(右边的 chunk)
    const newChunk = chunk.split(index)

    // 更新索引，将原始 chunk 和新 chunk 实例添加到索引中
    this.byEnd[index] = chunk
    this.byStart[index] = newChunk
    this.byEnd[newChunk.end] = newChunk

    // 如果拆分的是最后一个 chunk，更新 this.lastChunk 为新的 chunk 实例
    if (chunk === this.lastChunk) this.lastChunk = newChunk

    this.lastSearchedChunk = chunk
    return true
  }

  toString() {
    let str = this.intro

    let chunk = this.firstChunk
    while (chunk) {
      str += chunk.toString()
      chunk = chunk.next
    }

    return str + this.outro
  }

  isEmpty() {
    let chunk = this.firstChunk
    do {
      if ((chunk.intro.length && chunk.intro.trim()) || (chunk.content.length && chunk.content.trim()) || (chunk.outro.length && chunk.outro.trim()))
        return false
    } while ((chunk = chunk.next))
    return true
  }

  length() {
    let chunk = this.firstChunk
    let length = 0
    do {
      length += chunk.intro.length + chunk.content.length + chunk.outro.length
    } while ((chunk = chunk.next))
    return length
  }

  trimLines() {
    return this.trim('[\\r\\n]')
  }

  trim(charType) {
    return this.trimStart(charType).trimEnd(charType)
  }

  trimEndAborted(charType) {
    const rx = new RegExp((charType || '\\s') + '+$')

    this.outro = this.outro.replace(rx, '')
    if (this.outro.length) return true

    let chunk = this.lastChunk

    do {
      const end = chunk.end
      const aborted = chunk.trimEnd(rx)

      // if chunk was trimmed, we have a new lastChunk
      if (chunk.end !== end) {
        if (this.lastChunk === chunk) {
          this.lastChunk = chunk.next
        }

        this.byEnd[chunk.end] = chunk
        this.byStart[chunk.next.start] = chunk.next
        this.byEnd[chunk.next.end] = chunk.next
      }

      if (aborted) return true
      chunk = chunk.previous
    } while (chunk)

    return false
  }

  trimEnd(charType) {
    this.trimEndAborted(charType)
    return this
  }
  trimStartAborted(charType) {
    const rx = new RegExp('^' + (charType || '\\s') + '+')

    this.intro = this.intro.replace(rx, '')
    if (this.intro.length) return true

    let chunk = this.firstChunk

    do {
      const end = chunk.end
      const aborted = chunk.trimStart(rx)

      if (chunk.end !== end) {
        // special case...
        if (chunk === this.lastChunk) this.lastChunk = chunk.next

        this.byEnd[chunk.end] = chunk
        this.byStart[chunk.next.start] = chunk.next
        this.byEnd[chunk.next.end] = chunk.next
      }

      if (aborted) return true
      chunk = chunk.next
    } while (chunk)

    return false
  }

  trimStart(charType) {
    this.trimStartAborted(charType)
    return this
  }

  hasChanged() {
    return this.original !== this.toString()
  }

  _replaceRegexp(searchValue, replacement) {
    function getReplacement(match, str) {
      if (typeof replacement === 'string') {
        return replacement.replace(/\$(\$|&|\d+)/g, (_, i) => {
          // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/replace#specifying_a_string_as_a_parameter
          if (i === '$') return '$'
          if (i === '&') return match[0]
          const num = +i
          if (num < match.length) return match[+i]
          return `$${i}`
        })
      } else {
        return replacement(...match, match.index, str, match.groups)
      }
    }
    function matchAll(re, str) {
      let match
      const matches = []
      while ((match = re.exec(str))) {
        matches.push(match)
      }
      return matches
    }
    if (searchValue.global) {
      const matches = matchAll(searchValue, this.original)
      matches.forEach(match => {
        if (match.index != null) {
          const replacement = getReplacement(match, this.original)
          if (replacement !== match[0]) {
            this.overwrite(match.index, match.index + match[0].length, replacement)
          }
        }
      })
    } else {
      const match = this.original.match(searchValue)
      if (match && match.index != null) {
        const replacement = getReplacement(match, this.original)
        if (replacement !== match[0]) {
          this.overwrite(match.index, match.index + match[0].length, replacement)
        }
      }
    }
    return this
  }

  _replaceString(string, replacement) {
    const { original } = this
    const index = original.indexOf(string)

    if (index !== -1) {
      this.overwrite(index, index + string.length, replacement)
    }

    return this
  }

  replace(searchValue, replacement) {
    if (typeof searchValue === 'string') {
      return this._replaceString(searchValue, replacement)
    }

    return this._replaceRegexp(searchValue, replacement)
  }

  _replaceAllString(string, replacement) {
    const { original } = this
    const stringLength = string.length
    for (let index = original.indexOf(string); index !== -1; index = original.indexOf(string, index + stringLength)) {
      const previous = original.slice(index, index + stringLength)
      if (previous !== replacement) this.overwrite(index, index + stringLength, replacement)
    }

    return this
  }

  replaceAll(searchValue, replacement) {
    if (typeof searchValue === 'string') {
      return this._replaceAllString(searchValue, replacement)
    }

    if (!searchValue.global) {
      throw new TypeError('MagicString.prototype.replaceAll called with a non-global RegExp argument')
    }

    return this._replaceRegexp(searchValue, replacement)
  }
}
