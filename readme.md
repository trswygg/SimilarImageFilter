# SimilarImageFilter

相似图片过滤

## TODO

-   [x] improve perfomence
-   [x] trim output
-   [ ] ascii preview
-   [ ] everything API on windows
-   [ ] finishing parameter

## lib

`github.com/corona10/goimagehash`
`github.com/schollz/progressbar/v3`
`github.com/codeskyblue/go-libjpeg` 

## use

-h to see help

## build

see `github.com/codeskyblue/go-libjpeg`

## update

-   2020/11/10

    -   fix goimagehash ExtPerceptionHash() Resize() parameter Bugs
        -   it will cause `fatal error: out of memory allocating heap arena metadata` when -w and -h is big

    -   testing parameter

    -   2020/11/11
        -   trim output
        -   TODO: ascii preview