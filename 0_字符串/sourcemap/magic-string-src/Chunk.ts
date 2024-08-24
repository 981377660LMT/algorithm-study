// @ts-nocheck

/**
 * 链表节点.
 */
export default class Chunk {
  constructor(start, end, content) {
    this.start = start
    this.end = end
    this.original = content

    this.intro = ''
    this.outro = ''

    this.content = content
    this.storeName = false
    this.edited = false

    this.previous = null
    this.next = null
  }

  toString() {
    return this.intro + this.content + this.outro
  }

  appendLeft(content) {
    this.outro += content
  }

  prependLeft(content) {
    this.outro = content + this.outro
  }

  appendRight(content) {
    this.intro += content
  }

  prependRight(content) {
    this.intro = content + this.intro
  }

  clone() {
    const chunk = new Chunk(this.start, this.end, this.original)

    chunk.intro = this.intro
    chunk.outro = this.outro
    chunk.content = this.content
    chunk.storeName = this.storeName
    chunk.edited = this.edited

    return chunk
  }

  contains(index) {
    return this.start < index && index < this.end
  }

  eachNext(fn) {
    let chunk = this
    while (chunk) {
      fn(chunk)
      chunk = chunk.next
    }
  }

  eachPrevious(fn) {
    let chunk = this
    while (chunk) {
      fn(chunk)
      chunk = chunk.previous
    }
  }

  edit(content, storeName, contentOnly) {
    this.content = content
    if (!contentOnly) {
      this.intro = ''
      this.outro = ''
    }
    this.storeName = storeName

    this.edited = true

    return this
  }

  reset() {
    this.intro = ''
    this.outro = ''
    if (this.edited) {
      this.content = this.original
      this.storeName = false
      this.edited = false
    }
  }

  split(index) {
    const sliceIndex = index - this.start

    const originalBefore = this.original.slice(0, sliceIndex)
    const originalAfter = this.original.slice(sliceIndex)

    this.original = originalBefore

    // 创建新的 chunk 实例，表示拆分点后的部分
    const newChunk = new Chunk(index, this.end, originalAfter)
    newChunk.outro = this.outro
    this.outro = ''

    // 更新当前 chunk 的结束位置
    this.end = index

    // 如果当前 chunk 已经被编辑过，将新 chunk 设置为空字符串，否则保持原始内容
    if (this.edited) {
      // after split we should save the edit content record into the correct chunk
      // to make sure sourcemap correct
      // For example:
      // '  test'.trim()
      //     split   -> '  ' + 'test'
      //   ✔️ edit    -> '' + 'test'
      //   ✖️ edit    -> 'test' + ''
      // TODO is this block necessary?...
      newChunk.edit('', false)
      this.content = ''
    } else {
      this.content = originalBefore
    }

    newChunk.next = this.next
    if (newChunk.next) newChunk.next.previous = newChunk
    newChunk.previous = this
    this.next = newChunk

    return newChunk
  }

  trimEnd(rx) {
    this.outro = this.outro.replace(rx, '')
    if (this.outro.length) return true

    const trimmed = this.content.replace(rx, '')

    if (trimmed.length) {
      if (trimmed !== this.content) {
        this.split(this.start + trimmed.length).edit('', undefined, true)
        if (this.edited) {
          // save the change, if it has been edited
          this.edit(trimmed, this.storeName, true)
        }
      }
      return true
    } else {
      this.edit('', undefined, true)

      this.intro = this.intro.replace(rx, '')
      if (this.intro.length) return true
    }
  }

  trimStart(rx) {
    this.intro = this.intro.replace(rx, '')
    if (this.intro.length) return true

    const trimmed = this.content.replace(rx, '')

    if (trimmed.length) {
      if (trimmed !== this.content) {
        const newChunk = this.split(this.end - trimmed.length)
        if (this.edited) {
          // save the change, if it has been edited
          newChunk.edit(trimmed, this.storeName, true)
        }
        this.edit('', undefined, true)
      }
      return true
    } else {
      this.edit('', undefined, true)

      this.outro = this.outro.replace(rx, '')
      if (this.outro.length) return true
    }
  }
}
