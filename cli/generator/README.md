## tl;dr

test-appのgenerator群

生成ファイルを見つけるための正規表現
```
^// Code generated .* DO NOT EDIT\.$
```

### ddl
YAML定義(`docs/entity/**/*.yaml`)からddlファイル(`db/**/*.sql`)を作成する
TODO:
Extra関連未整備
partition未対応
    schemalex使ってるからそもそもmigrate出来ない

#### Usage
```
cd ddl;go generate;cd ..
```

##### yaml
- name
    - Entity名
    - パスカルをスネークケースにしてテーブル名に使用
- package
    - Entity Package
- structure
- primary
    - PrimaryKey
- index
    - IndexKey
    - Index名は,カンマ区切りされているカラム名を_区切り
- unique
    - UniqueKey
    - Unique名は,カンマ区切りされているカラム名を_区切り

#### entity
YAML定義(`docs/entity/**/*.yaml`)からEntityファイル(`app/domain/entity/**/*.gen.go`)を作成する

TODO:
nullableにするとgetTypeWithPointerしなきゃいけないの解明してない

#### Usage
```
cd entity;go generate;cd ..
```

#### enum
YAML定義(`docs/enum/**/*.yaml`)からEntityファイル(`app/enum/**/*.gen.go`)を作成する

#### Usage
```
cd enum;go generate;cd ..
```

#### protocol
TODO: 未実装

#### repository
YAML定義(`docs/entity/**/*.yaml`)からDomainRepositoryファイル(`app/domain/repository/**/*_repository.gen.go`)を作成する
作成メソッドはPrimary、Unique、Indexの各Keyに対してのFindとSaveのみ

#### Usage
```
cd repository;go generate;cd ..
```

#### vo
YAML定義(`docs/vo/***/**/*.yaml`)からvoファイル(`app/domain/value/***/**/*.gen.go`)を作成する

#### Usage
```
cd vo;go generate;cd ..
```
