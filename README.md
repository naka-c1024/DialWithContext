# DialWithContext

`net.Dial`関数の上位互換として、引数にコンテキストを追加し、コンテキストのキャンセルを検知する`DialWithContext`関数を作成。

ナレッジワークのインターン「Enablement Workshop for Gophers」の初日チーム課題。

## 実装

### コンテキストのキャンセル検知

2023年6月現在最新の Go1.21 で導入された`context.AfterFunc`関数を使用して、キャンセルが発生したら接続を閉じるという実装を行った。

[Go1.21 Release Notes](https://tip.golang.org/doc/go1.21)
[AfterFunc](https://pkg.go.dev/context@master#AfterFunc)

初期の実装方針では、`net.Conn`型の`SetDeadline`メソッドを使用して読み書きをキャンセルする方針だったが、

キャンセルが呼ばれているなら、コネクションの再利用は考えにくいと判断し、

`SetDeadline`を使用せずに直接接続をクローズする方針を採用した。

### Go1.20以下のバージョンへの対応

`context.AfterFunc`関数を使用せず、`Go1.20`以下のバージョンにも対応した。

`DialWithContext`の中で新しくキャンセルを検知するゴールーチンを立ち上げる実装にした。

リークを回避するため、`DialWithContext`の戻り値にゴールーチンをクローズする`cancelFunc`関数を追加した。
