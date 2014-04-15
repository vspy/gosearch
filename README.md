```
time ./index -stopwords=./src/gosearch/stopwords.txt enwiki-latest-pages-articles.xml

real  78m45.226s
user  195m0.323s
sys 3m58.667s
```

```
time ./search 駿

without actually searching something:
real  0m17.019s
user  0m15.878s
sys 0m1.122s

with search and printing results:
real  0m17.142s
user  0m15.977s
sys 0m1.125s
```
