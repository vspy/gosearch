1000 batch, wikisource:
```
real  34m44.178s
user  39m9.537s
sys 3m46.542s
```

5000 batch, wikisource:
```
real  24m24.028s
user  37m39.359s
sys 2m52.240s
```

5000 batch, wikisource, buffered writer, optimized index elements:
```
real  15m49.556s
user  32m45.092s
sys 2m8.009s
```

1000 batch, wikisource, buffered writer, perfect batching:
```
real  12m39.609s
user  33m1.934s
sys 2m5.266s
```

1000, wikisource, buffered, proper batching, FSM tokenizer:
```
real  11m10.720s
user  32m24.632s
sys 1m34.062s
```

current revisions enwiki, 25m pages:
```
real  163m13.560s
user  410m6.874s
sys 8m4.561s
```

enwiki articles
```
real  77m28.131s
user  185m21.550s
sys 3m37.201s
```
