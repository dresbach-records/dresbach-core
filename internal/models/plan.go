package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// BillingCycle define os ciclos de faturamento disponíveis.
type BillingCycle string

const (
	Monthly      BillingCycle = "monthly"
	Annually     BillingCycle = "annually"
	Quarterly    BillingCycle = "quarterly"
	Semiannually BillingCycle = "semiannually"
	Biennially   BillingCycle = "biennially"
)

// PlanStatus define o status de visibilidade de um plano.
type PlanStatus string

const (
	PlanStatusActive   PlanStatus = "active"
	PlanStatusHidden   PlanStatus = "hidden"
	PlanStatusArchived PlanStatus = "archived"
)

// Plan representa um produto ou serviço vendável.
type Plan struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Description    sql.NullString  `json:"description"`
	Category       sql.NullString  `json:"category"`
	Price          float64         `json:"price"`
	BillingCycle   BillingCycle    `json:"billing_cycle"`
	Features       json.RawMessage `json:"features"`
	WhmPackageName sql.NullString  `json:"whm_package_name"`
	Status         PlanStatus      `json:"status"`
	DisplayOrder   int             `json:"display_order"`
	IsFeatured     bool            `json:"is_featured"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// GetAllPlans busca todos os planos, ordenados para exibição.
func GetAllPlans(db *sql.DB) ([]Plan, error) {
	query := `SELECT id, name, description, category, price, billing_cycle, features, whm_package_name, status, display_order, is_featured, created_at, updated_at FROM plans WHERE status != 'archived' ORDER BY display_order ASC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Category, &p.Price, &p.BillingCycle, &p.Features, &p.WhmPackageName, &p.Status, &p.DisplayOrder, &p.IsFeatured, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

// CreatePlan insere um novo plano no banco de dados.
func CreatePlan(db *sql.DB, p *Plan) (int64, error) {
	query := `INSERT INTO plans (name, description, category, price, billing_cycle, features, whm_package_name, status, display_order, is_featured) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, p.Name, p.Description, p.Category, p.Price, p.BillingCycle, p.Features, p.WhmPackageName, p.Status, p.DisplayOrder, p.IsFeatured)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdatePlan atualiza um plano existente.
func UpdatePlan(db *sql.DB, p *Plan) error {
	query := `UPDATE plans SET name = ?, description = ?, category = ?, price = ?, billing_cycle = ?, features = ?, whm_package_name = ?, status = ?, display_order = ?, is_featured = ? WHERE id = ?`
	_, err := db.Exec(query, p.Name, p.Description, p.Category, p.Price, p.BillingCycle, p.Features, p.WhmPackageName, p.Status, p.DisplayOrder, p.IsFeatured, p.ID)
	return err
}

// UpdatePlanStatus atualiza apenas o status de um plano.
func UpdatePlanStatus(db *sql.DB, id int, status PlanStatus) error {
	query := `UPDATE plans SET status = ? WHERE id = ?`
	_, err := db.Exec(query, status, id)
	return err
}
