# WordSeg
WordSeg是一个使用Go语言实现的简单中文分词工具。

## 安装
```sh
$ go get github.com/ling0322/wordseg
$ go install github.com/ling0322/wordseg/...
```

## 模型

WordSeg的模型文件包括一个json格式的配置文件以及它指向的其他数据文件。可以从GitHub中WordSeg的[Release Page](https://github.com/ling0322/wordseg/releases/tag/v0.1)获取 (model.tar.gz). 也可以自己生成

### 生成模型

生成模型需要一个从语料库中提取的词频文件作为输入，它的格式是
```
WORD1 COUNT1
WORD2 COUNT2
... 
```

词频文件也可以从 [WEBDICT](https://github.com/ling0322/webdict) 项目中下载: webdict_with_freq.txt

```sh
$ wget https://github.com/ling0322/webdict/raw/master/webdict_with_freq.txt
```

词频文件准备好后，可以使用`bin/gen_model`生成模型

```sh
$ mkdir model
$ bin/gen_model -in webdict_with_freq.txt -out model
ok
$ ls -lh model
total 5.4M
-rwxrwxrwx 1 ling0322 ling0322 863K Jul 17 17:23 cost.uni
-rwxrwxrwx 1 ling0322 ling0322 4.6M Jul 17 17:23 lexicon
-rwxrwxrwx 1 ling0322 ling0322   61 Jul 17 17:23 wordseg.conf
```

## 命令行

wordseg安装后会生成`bin/wordseg`,可以用它来对文件进行简单的分词。将模型文件准备好后即可使用：

```sh
$ echo '存钱准备买小裙子' | bin/wordseg -c model/wordseg.conf
存钱 准备 买 小 裙子
```

## API调用例子

```go
package main

import (
    "fmt"
    "log"

    "github.com/ling0322/wordseg"
)

func main() {
    s, err := wordseg.NewSegmenter("model/wordseg.conf")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(s.Seg("存钱准备买小裙子"))
}
```

