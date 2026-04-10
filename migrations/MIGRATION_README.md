# onchain-task-tracker 初始化 SQL / migration

這份 migration 是依照這次實際跑出來的錯誤與已確認可用的 API 行為整理出的 **best-effort 初版**。

已確認需求：
- `tasks` 至少需要：`task_id`, `description`, `status`, `priority`, `due_date`
- `task_blockchain_logs` 需要：`task_id`, `action`, `tx_hash`, `chain_id`, `contract_address`, `status`

## 檔案
- `001_init_schema.sql`

## 執行方式
在 PowerShell：

```powershell
docker exec -i onchain-postgres psql -U postgres -d postgres < 001_init_schema.sql
```

如果你是在 Windows PowerShell，且 SQL 檔在目前目錄，也可以：

```powershell
Get-Content .\001_init_schema.sql | docker exec -i onchain-postgres psql -U postgres -d postgres
```

## 建議目錄結構
```text
onchain-task-tracker_go-service/
  migrations/
    001_init_schema.sql
```

## 建議後續做法
1. 把這個 SQL 放進專案的 `migrations/` 目錄
2. 再從程式碼搜尋 `FROM tasks`, `INSERT INTO tasks`, `UPDATE tasks`, `FROM task_blockchain_logs`
3. 若發現還有其他欄位需求，再補成 `002_add_xxx.sql`

## 我刻意保守的地方
因為這次我沒有直接讀到你完整 repo 原始碼，所以：
- `tasks` 的欄位只先放 **目前已經明確看到會用到** 的那幾個
- 沒有硬塞未證實的欄位
- migration 盡量做成可重跑、對既有 DB 破壞較小
