package repository

import "database/sql"

type NFTOrder struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	Price           string `json:"price"`
	RecipientWallet string `json:"recipientWallet"`
	CreatorWallet   string `json:"creatorWallet"`

	Sold           bool   `json:"sold"`            // ⭐ 新增
	PurchaseTxHash string `json:"purchaseTxHash"` // ⭐ 新增
}

type NFTOrderRepository struct {
	db *sql.DB
}

func NewNFTOrderRepository(db *sql.DB) *NFTOrderRepository {
	return &NFTOrderRepository{db: db}
}

func (r *NFTOrderRepository) Create(order NFTOrder) (int, error) {
	var id int

	err := r.db.QueryRow(`
		INSERT INTO nft_orders (
			title, description, image, price, recipient_wallet, creator_wallet
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		order.Title,
		order.Description,
		order.Image,
		order.Price,
		order.RecipientWallet,
		order.CreatorWallet,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *NFTOrderRepository) FindAll() ([]NFTOrder, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, image, price, recipient_wallet, creator_wallet, sold, purchase_tx_hash
		FROM nft_orders
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]NFTOrder, 0)

	for rows.Next() {
		var order NFTOrder
		err := rows.Scan(
	&order.ID,
	&order.Title,
	&order.Description,
	&order.Image,
	&order.Price,
	&order.RecipientWallet,
	&order.CreatorWallet,
	&order.Sold,
	&order.PurchaseTxHash,
)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *NFTOrderRepository) UpdatePurchase(id int, txHash string) error {
	_, err := r.db.Exec(`
		UPDATE nft_orders
		SET sold = true,
		    purchase_tx_hash = $1,
		    updated_at = now()
		WHERE id = $2
	`, txHash, id)

	return err
}