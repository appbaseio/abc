# CSV

CSV adaptor works for csv files.

A basic config.env looks like the following.
We have an additional parameter `typename` in csv adaptor because csv files only have data and no concept of tables / types. 
So we need to define it manually.

```ini
src.type=csv
src.uri=/full/local/path/to/file.csv
src.typename=type_name_to_use

dest.type=elasticsearch
dest.uri=https://USER:PASSWORD@SERVER/INDEX
```
