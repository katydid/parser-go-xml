# Parsarray

Parse array parses all XML nodes as arrays

For example 

```xml
<a k1="v1" k2="v2">b</a>
```

Is parsed as:
```
[{"a": [{"k1": "v1"}, {"k2": v2"}, "b"]}]
```

This is technically more correct, but less easy to swap out with other formats.