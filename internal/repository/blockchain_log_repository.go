package repository

import (
	"database/sql"
	"go-service/internal/model"
)

type BlockchainLogRepository struct {
	db *sql.DB
}

func NewBlockchainLogRepository(db *sql.DB) *BlockchainLogRepository {
	return &BlockchainLogRepository{db: db}
}

func (r *BlockchainLogRepository) Create(log model.BlockchainLog) error {
	_, err := r.db.Exec(`
        INSERT INTO task_blockchain_logs (task_id, action, tx_hash, chain_id, contract_address, status)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, log.TaskID, log.Action, log.TxHash, log.ChainID, log.ContractAddress, log.Status)
	return err
}

func (r *BlockchainLogRepository) FindAll() ([]model.BlockchainLog, error) {
	rows, err := r.db.Query(`
        SELECT id, task_id, action, tx_hash, chain_id, contract_address, status, created_at
        FROM task_blockchain_logs
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]model.BlockchainLog, 0)
	for rows.Next() {
		var l model.BlockchainLog
		if err := rows.Scan(
			&l.ID, &l.TaskID, &l.Action, &l.TxHash,
			&l.ChainID, &l.ContractAddress, &l.Status, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
