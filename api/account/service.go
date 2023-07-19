package account

import (
	"auction-website/conf"
)

type Service struct {
	repo *Repository
}

func NewService(c *conf.Config) *Service {
	return &Service{
		repo: NewRepository(c),
	}
}

// Deposit adds the specified amount to the account balance
func (s *Service) Deposit(accountID uint32, amount float64) error {
	// Call the corresponding function in the Repository struct to update the account balance
	err := s.repo.UpdateBalance(accountID, amount)
	if err != nil {
		return err
	}
	return nil
}

// Withdraw subtracts the specified amount from the account balance
func (s *Service) Withdraw(accountID uint32, amount float64) error {
	// Call the corresponding function in the Repository struct to check if there is enough balance and update the account balance
	orderNo := "orderNoxxx"
	err := s.repo.Withdraw(accountID, orderNo, amount)
	//fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

// GetBalance retrieves the current balance of the specified account
func (s *Service) GetBalance(accountID uint32) (float64, error) {
	// Call the corresponding function in the Repository struct to get the current balance of the account
	balance, err := s.repo.GetBalance(accountID)
	if err != nil {
		return 0.0, err
	}
	return balance, nil
}

// Define a type to hold the withdrawal records
//type WithdrawalRecord struct {
//	AccountID uint32
//	OrderNo   string
//	Amount    float64
//}

// GetWithdrawalRecords retrieves all withdrawal records for the specified account
func (s *Service) WithdrawalRecord(accountID, page, pageSize uint32) ([]WithdrawalRecord, error) {
	// Call the corresponding function in the Repository struct to get all withdrawal records for the account
	withdrawals, err := s.repo.GetWithdrawalRecord(accountID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
