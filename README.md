### TL;DR
小松原Golang勉強用プロジェクト

### Usage
#### サーバー起動
1. ./init.sh
2. ./update.sh
3. ./run.sh

#### generate
```cd cli/generator/;./run.sh```

swagger
```cd app;swag init -g main.go -o ../docs/protocol/api;cd ../```

### もっと整備したいところ
#### 環境周り
- local環境切り分け
    - ENGがlocal環境を立てる分にはIDEでgoを実行すれば問題ないが、PLANやARTがlocal環境を立てる際にapp(go)もコンテナとして立てる環境を作ったほうが良い
- mysql
    - ユーザー管理
        - めんどくてrootのままアクセスしてる
        - ローカルだし優先度低い
    - AUTOCOMMIT=1になってる問題
- テスト環境構築
    - go、k6のコンテナを作って疎通させる

#### 実装周り
- user db sharding
    - コンテナだけは1、2用意してるがappから動的に接続先を切り替えてない
    - SINはspannerだから考えなくていい？
- Openのスレッドセーフがどういう意味か調査する必要がある
    - 1リクエスト内でスレッドセーフなのかリクエスト跨いでスレッドセーフなのかとか
    - Goのmapはスレッドセーフじゃないらしいが、GinでHTTPリクエストを受け取るとGoroutine化(=スレッド化)されないんだろうか？
- UserDBのSpanner化
- masterのgorm接続
  - 多分VO周りでコケるからgeneratorから修正が必要
    - item_masterのScanエラーはint64を定義したら解消した
    - でもinsert時にエラーになる
      - valuerの問題？
- swagger
    - VOがprimitiveとしてswaggerに表示されない
    - titleが省略できない
    - Container化
- generator
    - protocol
    - infrastructure/repository
        - domainでメソッド自体を縛っているし可能な希ガス
- gRPC化
- マスターインポート
    - refactor
    - 実行前にdbを空にするように
    - validation
    - save系を動的にしたい
        - Golangの時点で完全に動的は無理
        - genにimporter用の定義ファイルとか構造とかを出力するようにする必要がありそう
- マスターコンバート
  - 既存ファイルへの書き込みがAccess is denied.になる
- repository
  - DIPの関係でdomain repositoryに定義しているメソッドしか参照できない
  - gen側をもっと拡張しないと開発に耐えられなそう
    - リスト取得とかの条件とか
    - LIMITとか
      - これはinfra repository側でどうにかしてもいい気がする
  - そもそもrepositoryまでgenする必要ない？
    - genで縛ることの
      - メリット
        - 楽
        - 各indexでのfindしかないからindex貼り忘れとかなくなる
      - デメリット
        - 自由度下がる
          - 自由度保つには何もかもgen出来るようにする必要がある
        - indexを貼らないとfind出来ない
          - これに関してはテーブル設計の問題もある

### 参考資料
#### Go

https://zenn.dev/jy8752/scraps/0149823a972676
Go言語100Tipsを読んで

#### Spanner
詳解 google-cloud-go/spanner — トランザクション編
https://medium.com/google-cloud-jp/%E8%A9%B3%E8%A7%A3-google-cloud-go-spanner-%E3%83%88%E3%83%A9%E3%83%B3%E3%82%B6%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3%E7%B7%A8-6b63099bd7fe

Cloud Spannerの論理シャーディングを理解する
https://zenn.dev/facengineer/articles/6c66870d484407

Go で Spanner とよろしくやるためにガチャガチャやっている話
https://chidakiyo.hatenablog.com/entry/2020/12/14/go-spanner-tools

Google Cloud Spanner用のORMを公開しました
https://developers.10antz.co.jp/archives/2237

cloudspannerecosystem/yo
https://github.com/cloudspannerecosystem/yo

【備忘】spanner の mutation と Statement DML のどちらを使うか
https://blog.framinal.life/entry/2022/09/29/000101

Cloud Spanner でのタグ付けのススメ
https://zenn.dev/google_cloud_jp/articles/756f47aa39b2b6

Cloud Spannerのスプリット分散をわかった気になる
https://zenn.dev/facengineer/articles/bca8790087b0e4

SpannerのTransactionでアプリケーションからAbortedを返してはいけない
https://zenn.dev/ryo_yamaoka/articles/316f26b0047f58

GCP SpannerにおけるStatement DMLとMutationの挙動の違い
https://zenn.dev/ryo_yamaoka/articles/1dd9799da26440

ツールを使った Cloud Spanner のウォームアップ
https://zenn.dev/google_cloud_jp/articles/1adfdfc1fa5a6c

Cloud Spanner における各種トランザクションの使い分け
https://zenn.dev/google_cloud_jp/articles/15d34df66becfe

【Go】Service層でSpannerのトランザクション管理をしたい
https://zenn.dev/ymtdzzz/articles/66288138744973

Cloud Spannerのパフォーマンスチューニングの勘所
https://zenn.dev/facengineer/articles/cc0cab5c7e9a1c

Cloud Spanner のロックについて
https://zenn.dev/apstndb/articles/a62ac78b3b91bb
