interface JSONSerializer {
  toJSON(): string
}

interface JSONDeserializer {
  fromJSON(json: string): this
}

export {}
