# natsume-cabocha-bindings

Simple Go bindings for CaboCha.
This package was developed in tandem with [natsume-server](https://github.com/borh/natsume-server) as its CaboCha provider.

# Usage

An example HTTP and WebSocket server that serves CaboCha in normal lattice or JSON output is provided in the examples subfolder:

```bash
$ go run examples/example-server.go
```

Then in another shell:

```bash
$ curl http://localhost:8080/ -d レスポンスを返す
* 0 1D 0/1 0.000000
レスポンス      名詞,普通名詞,一般,*,*,*,レスポンス,レスポンス,レスポンス,レスポンス,レスポンス,外,レスポンス,レスポンス,レスポンス,レスポンス,*,*,*,*,*,*,"1,3",C1,*   O
を      助詞,格助詞,*,*,*,*,ヲ,を,を,オ,ヲ,和,を,オ,ヲ,ヲ,*,*,*,*,*,*,*,"動詞%F2@0,名詞%F1,形容詞%F2@-1",*      O
* 1 -1D 0/0 0.000000
返す    動詞,一般,*,*,五段-サ行,終止形-一般,カエス,返す,返す,カエス,カエス,和,返す,カエス,カエス,カエス,*,*,*,*,*,*,1,C1,*      O
EOS
```

# Version

0.1

# TODO

- write tests
- stress-test the WebSocket implementation
- decide on a standard JSON format
- streamline usage
