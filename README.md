# gxjpath
XJPath for Go

The basic use case of this library is to lookup values from deeps of decoded
JSON structures into interface{}. It happens very frequently that unknown 
JSON structure may be received, however, you might be interested to check if there is
any value exist. Using this library you can drammatically simplify the way how you look
for those values.

Here is a bunch of examples.

```Go
// Lookup for a key "anymapkey" and then a keey "k1"
v, err := := gxjpath.LookupRawPath("anymapkey.k1", data)

// Lookup for a key "anyarraykey" and then for a last array element.
v, err := := gxjpath.LookupRawPath("anyarraykey.@last", data)

// Lookup for a key "anyarraykey" and then for a last-1.
v, err := := gxjpath.LookupRawPath("anyarraykey.@-2", data)

// Lookup for a key "anyarraykey" and then for a key "@-2".
v, err := := gxjpath.LookupRawPath("anyarraykey.\\@-2", data)

// Long path example.
// 1. Key1 -> key2 -> last array element -> first array element-> key3 -> third array element -> key.123
v, err := := gxjpath.LookupRawPath("key1.key2.@last.@first.key3.@3.key\\.123", data)
```

# Motive
People ask me why did I implement this thing again. Here is why:
  
  - Performance. A data lookup(especially chached one) is much faster.
  - I've implemented the same thing in Python, but now I need performance. So, compatibility need.
  - I can look-up last array element. This is a killer feature for me.

Additional extentions are welcome.
