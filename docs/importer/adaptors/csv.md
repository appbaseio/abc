# CSV

CSV adaptor works for csv files.

A basic config.env looks like the following.
We have an additional parameter `typename` in csv adaptor because csv files only have data and no concept of tables / types. 
So we need to define it manually.

```ini
src_type=csv
src_uri=/full/local/path/to/file.csv
typename=type_name_to_use

dest_type=elasticsearch
dest_uri=https://USER:PASSWORD@SERVER/INDEX
```
