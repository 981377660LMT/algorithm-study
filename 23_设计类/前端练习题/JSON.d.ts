type JsonValue = JsonPrimitive | JsonObject | JsonArray
type JsonPrimitive = string | number | boolean | null
type JsonObject = { [Key in string]?: JsonValue }
type JsonArray = JsonValue[]
