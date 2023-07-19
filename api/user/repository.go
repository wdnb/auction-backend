package user

import (
	"auction-website/conf"
	db "auction-website/database/connectors/mysql"

	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (r *Repository) CreateUser(u *User) (uint32, error) {
	result, err := r.db.NamedExec(`INSERT INTO user (username, password, email)
		VALUES (:username, :password, :email)`, u)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), err
}

func (r *Repository) CreateUserLog(u *UserLog) (uint32, error) {
	result, err := r.db.NamedExec(`INSERT INTO user_log (user_id, ip, operation,init_time,details)
		VALUES (:user_id, :ip, :operation,:init_time,:details)`, u)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), err
}

func (r *Repository) GetUserByID(id uint32) (*User, error) {
	var u User
	err := r.db.Get(&u, `SELECT id,username,email FROM user WHERE id=?`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var u User
	query := "SELECT id,username,password FROM user WHERE username = ?"
	err := r.db.Get(&u, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) GetVerificationCodeByPhone(phone, kind string) (*VerificationCode, error) {
	//var code string
	var code VerificationCode
	now := time.Now().Unix()
	query := "SELECT id,code FROM verification_code WHERE phone = ? AND kind=? AND expired_at>=? AND is_used=?"
	err := r.db.Get(&code, query, phone, kind, now, "unused")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCodeNotFound
		}
		return nil, err
	}
	return &code, nil
}

// 将 is_used 字段更新为已使用
func (r *Repository) UpdateVerificationCode(id uint32) error {

	updateQuery := "UPDATE verification_code SET is_used = ? WHERE id = ?"
	_, err := r.db.Exec(updateQuery, "used", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserByPhone(phone string) (*CreateUser, error) {
	//fmt.Println(phone)
	var u CreateUser
	query := "SELECT id,username,phone FROM user WHERE phone = ?"
	err := r.db.Get(&u, query, phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	//fmt.Println(u)
	return &u, nil
}

func (r *Repository) UserNameIsExist(username string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM user WHERE username = ?"
	err := r.db.Get(&count, query, username)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) UserEmailIsExist(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM user WHERE email = ?"
	err := r.db.Get(&count, query, email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) UserPhoneIsExist(phone string) (bool, error) {
	var count int
	//fmt.Println(phone)
	query := "SELECT COUNT(*) FROM user WHERE phone = ?"
	err := r.db.Get(&count, query, phone)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CheckUserExistByID(id uint32) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM user WHERE id = ?"
	err := r.db.Get(&count, query, id)
	if err != nil {
		return false, err
	}
	//fmt.Println(count)
	return count > 0, nil
}

func (r *Repository) UpdateUserByUserID(u *UpdateUser) error {
	_, err := r.db.NamedExec(`UPDATE user SET username=:username, password=:password, email=:email WHERE id=:id`, u)
	return err
}

func (r *Repository) CreateShippingAddress(sp *CreateShippingAddress) (uint32, error) {
	result, err := r.db.NamedExec(`INSERT INTO shipping_address (user_id, phone, is_active,recipient_name,region, address)
		VALUES (:user_id, :phone, :is_active,:recipient_name,:region, :address)`, sp)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), err
}

func (r *Repository) UpdateUserAddressByID(sp *ShippingAddress) (uint32, error) {
	result, err := r.db.NamedExec(`UPDATE shipping_address SET user_id=:user_id, phone=:phone, is_active=:is_active, recipient_name=:recipient_name, region=:region ,address=:address WHERE id=:id`, sp)
	num, err := result.RowsAffected()
	//fmt.Println(id)
	if err != nil {
		return 0, err
	}
	return uint32(num), err
}

// Define a method to update the is_active field of a shipping address by its ID
func (r *Repository) UpdateShippingAddressIsActiveByUID(uid uint32, isActive bool) error {
	_, err := r.db.Exec(`UPDATE shipping_address SET is_active=? WHERE user_id=?`, isActive, uid)
	return err
}

// Define a method to get all shipping addresses of a user by their user ID
func (r *Repository) GetUserAddresses(uid uint32) ([]*ShippingAddress, error) {
	var addresses []*ShippingAddress
	err := r.db.Select(&addresses, `SELECT id,user_id,phone,is_active,recipient_name,region,address FROM shipping_address WHERE user_id=?`, uid)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// Define a method to get the active shipping address of a user by their user ID
func (r *Repository) GetActiveUserAddress(uid uint32) (*ShippingAddress, error) {
	var sp ShippingAddress
	err := r.db.Get(&sp, `SELECT id,user_id,phone,is_active,recipient_name,region,address FROM shipping_address WHERE user_id=? AND is_active=true`, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("active shipping address for user with id %d not found", uid)
		}
		return nil, err
	}
	return &sp, nil
}

// Define a method to delete a shipping address by its ID
func (r *Repository) DeleteUserAddressByID(userID, id uint32) (uint32, error) {
	result, err := r.db.Exec(`DELETE FROM shipping_address WHERE id=? AND user_id=?`, id, userID)
	num, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint32(num), err
}

// Define a method to create a user role when creating a user
func (r *Repository) CreateUserWithRole(u *CreateUser, roleName string) (uint32, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	result, err := tx.NamedExec(`INSERT INTO user (username, password,phone, email)
		VALUES (:username, :password,:phone, :email)`, u)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	var roleID uint32
	err = tx.Get(&roleID, `SELECT id FROM roles WHERE name = ?`, roleName)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(`INSERT INTO user_role (user_id, role_id) VALUES (?, ?)`, id, roleID)
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

// Define a method to get the roles of a user by their username
func (r *Repository) GetUserRoles(username string) (*Role, error) {
	//username = "8"
	//fmt.Println(username)
	var roles Role
	err := r.db.Get(&roles, `
		SELECT cr.p_type, cr.v0, cr.v1, cr.v2, r.name as role_name
		FROM casbin_rule cr
		JOIN roles r ON r.name = cr.v0
		JOIN user_role ur ON ur.role_id = r.id
		JOIN user u ON u.id = ur.user_id
		WHERE u.username = ?
	`, username)
	//fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

func (r *Repository) GetUserRolesByUid(id uint32) (*Role, error) {
	var roles Role
	err := r.db.Get(&roles, `
		SELECT cr.p_type, cr.v0, cr.v1, cr.v2, r.name as role_name
		FROM casbin_rule cr
		JOIN roles r ON r.name = cr.v0
		JOIN user_role ur ON ur.role_id = r.id
		JOIN user u ON u.id = ur.user_id
		WHERE u.id = ?
	`, id)
	//fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

// 发送手机验证码
func (r *Repository) SetVerificationCodeByPhone(phone, code, kind string) error {
	expiredAt := time.Now().Add(5 * time.Minute).Unix()
	_, err := r.db.Exec(`INSERT INTO verification_code (phone, code, kind,expired_at) VALUES (?, ?,?, ?)`, phone, code, kind, expiredAt)
	return err
}
