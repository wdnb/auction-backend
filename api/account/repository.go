package account

import (
	"auction-website/conf"
	db "auction-website/database/connectors/mysql"
	"auction-website/utils"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(c *conf.Config) *Repository {
	return &Repository{
		db: db.GetClient(c.Mysql),
	}
}

func (r *Repository) UpdateBalance(accountID uint32, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM account WHERE user_id = ?", accountID).Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}

	if count == 0 {
		_, err = tx.Exec("INSERT INTO account (user_id, balance, status) VALUES (?, ?, ?)", accountID, amount, "normal")
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.Exec("UPDATE account SET balance = balance + ? WHERE user_id = ?", amount, accountID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

//func (r *Repository) CheckBalance(accountID uint32) (float64, error) {
//	var balance float64
//	err := r.db.QueryRow("SELECT balance FROM account WHERE user_id = ?", accountID).Scan(&balance)
//	if err != nil {
//		return 0, err
//	}
//
//	return balance, nil
//}

//func (r *Repository) CheckWithdrawalEligibility(accountID uint32, withdrawalAmount float64) (bool, error) {
//	var balance float64
//	err := r.db.QueryRow("SELECT balance FROM account WHERE user_id = ?", accountID).Scan(&balance)
//	if err != nil {
//		return false, err
//	}
//
//	if balance >= withdrawalAmount {
//		return true, nil
//	} else {
//		return false, nil
//	}
//}
//
// We need to add a new method to the Repository struct called Withdraw that will handle the withdrawal of funds from an account.
// We will use a transaction to ensure that the update to the account balance and the creation of a transaction record occur atomically.
// Here's the code for the Withdraw method:

func (r *Repository) Withdraw(accountID uint32, orderNo string, withdrawalAmount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var balance float64
	err = tx.QueryRow("SELECT balance FROM account WHERE user_id = ?", accountID).Scan(&balance)
	if err != nil {
		tx.Rollback()
		return err
	}

	if balance < withdrawalAmount {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	_, err = tx.Exec("UPDATE account SET balance = balance - ? WHERE user_id = ?", withdrawalAmount, accountID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO withdraw_log (user_id,order_no, amount) VALUES (?,?,?)", accountID, orderNo, withdrawalAmount)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// We need to add a new method to the Repository struct called GetBalance that will retrieve the current balance of an account.
// Here's the code for the GetBalance method:

func (r *Repository) GetBalance(accountID uint32) (float64, error) {
	var balance float64
	err := r.db.QueryRow("SELECT balance FROM account WHERE user_id = ?", accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *Repository) GetWithdrawalRecord(accountID, page, pageSize uint32) ([]WithdrawalRecord, error) {
	var withdrawals []WithdrawalRecord
	offset := utils.Offset(page, pageSize)
	err := r.db.Select(&withdrawals, `SELECT * FROM withdraw_log WHERE user_id = ?  LIMIT ? OFFSET ?`, accountID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
