# Knowledge Production Line System (KPLS)

知的生産を製造ライン方式で管理するCLIツール

## 概要

KPLSは、知的生産（企画書・設計方針・比較検討・議事録など）を製造ライン（工程固定＋検査ゲート＋WIP制御）として実装し、アウトプットを連続生産できる状態を作るシステムです。

## 機能

- **ジョブ管理**: ジョブチケットの作成、一覧表示、詳細表示、工程進行、差し戻し
- **品質ゲート**: IQC（受入検査）、IPQC-1（骨格ゲート）、IPQC-2（根拠ゲート）、FQC（最終検査）
- **テンプレート管理**: 1pager、比較検討、PRDの3種類のテンプレート
- **統計・可視化**: WIP状況、不良率、工程別滞留時間の可視化

## インストール

```bash
go mod download
go build -o kpls
```

## 使用方法

### ジョブの作成

```bash
kpls job create
```

### ジョブ一覧（カンバン表示）

```bash
kpls job list
```

### ジョブ詳細

```bash
kpls job show <job-id>
```

### 工程を進める

```bash
kpls job advance <job-id>
```

### 差し戻し

```bash
kpls job reject <job-id> --reason "理由" --defects D01,D02
```

### 品質ゲートチェック

```bash
kpls check iqc <job-id>
kpls check ipqc1 <job-id>
kpls check ipqc2 <job-id>
kpls check fqc <job-id>
```

### 統計表示

```bash
kpls stats
```

### テンプレート一覧

```bash
kpls template list
kpls template show <template-name>
```

## ワークフロー

```
Backlog → IQC → Skeleton → IPQC-1 → Draft → IPQC-2 → Packaging → FQC → Done
```

## ディレクトリ構成

```
kpls/
  templates/          # 出力テンプレート
  quality/            # 品質チェックリスト
  ops/                # 運用手順書
data/
  jobs/               # ジョブデータ（JSON）
  defects/            # 不良記録（JSON）
```

## 開発

```bash
# テスト実行
go test ./...

# ビルド
go build -o kpls

# 実行
./kpls --help
```

## ライセンス

内部利用
