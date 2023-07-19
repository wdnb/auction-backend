package user

import (
	"auction-website/conf"
	"auction-website/middleware"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(c *conf.Config) *Service {
	return &Service{
		repo: NewRepository(c),
	}
}

func (s *Service) UserNameIsExist(u *User) (bool, error) {
	exist, err := s.repo.UserNameIsExist(u.Username)
	if err != nil {
		return false, err
	}
	if exist {
		return true, nil
	}
	return false, nil
}

func (s *Service) UserEmailIsExist(u *User) (bool, error) {
	exist, err := s.repo.UserEmailIsExist(u.Email)
	if err != nil {
		return false, err
	}
	if exist {
		return true, nil
	}
	return false, nil
}

func (s *Service) UserPhoneIsExist(u *VerificationCode) (bool, error) {
	exist, err := s.repo.UserPhoneIsExist(u.Phone)
	//fmt.Println(exist)
	if err != nil {
		return false, err
	}
	if exist {
		return true, nil
	}
	return false, nil
}

func (s *Service) CreateUser(u *CreateUser) (string, error) {
	hashedPassword, err := s.HashPassword(u.Password)
	if err != nil {
		return "", err
	}
	u.Password = hashedPassword
	uid, err := s.repo.CreateUserWithRole(u, "buyer") //默认角色 buyer
	if err != nil {
		return "", err
	}
	//u.ID = uid
	roles, err := s.repo.GetUserRolesByUid(uid)
	if err != nil {
		return "", err
	}
	//fmt.Println(roles)
	//fmt.Println(u.ID, u.Username, roles.RoleName)
	token, err := middleware.GenerateToken(uid, roles.RoleName)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Login(username, password string) (string, error) {
	u, err := s.repo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	roles, err := s.repo.GetUserRoles(username)
	//fmt.Println(roles)
	if err != nil {
		return "", err
	}
	passwordMatch := s.CheckPasswordHash(password, u.Password)
	if passwordMatch != true {
		return "", ErrIncorrectPassword
	}
	//t := user.RoleUser{
	//	User: *(u),
	//	Role: *(roles),
	//}
	//token, err := utils.GenerateToken(t)
	token, err := middleware.GenerateToken(u.ID, roles.RoleName)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) LoginByCode(clu *VerificationCode, ip string) (string, error) {
	phoneIsExist, err := s.UserPhoneIsExist(clu)
	if err != nil {
		return "", err
	}
	//检验验证码
	err = s.VerificationCodeByPhone(clu.Phone, clu.Code, clu.Kind)
	if err != nil {
		return "", err
	}
	var token string
	//走登录逻辑
	if phoneIsExist {
		u, err := s.repo.GetUserByPhone(clu.Phone)
		if err != nil {
			return "", err
		}
		roles, err := s.repo.GetUserRolesByUid(u.ID)
		if err != nil {
			return "", err
		}
		token, err = middleware.GenerateToken(u.ID, roles.RoleName)
		_ = s.LogUserOperation(u.ID, ip, OperationLogin, "手机号登录")
	} else {
		//走注册逻辑
		var cu CreateUser
		cu.Password, _ = generateRandomPassword(16)
		cu.Phone = &clu.Phone
		token, err = s.CreateUser(&cu)
		_ = s.LogUserOperation(0, ip, OperationRegister, "用户首次注册")
	}
	return token, err
}

func generateRandomPassword(length int) (string, error) {
	buffer := make([]byte, length)
	r := rand.Reader
	_, err := io.ReadFull(r, buffer)
	if err != nil {
		return "", err
	}
	password := base64.URLEncoding.EncodeToString(buffer)
	return password[:length], nil
}

func (s *Service) SendVerificationCode(phone, kind string) (string, error) {
	code := generateVerificationCode(6)
	err := s.repo.SetVerificationCodeByPhone(phone, code, kind)
	if err != nil {
		return "", err
	}
	return code, nil
}

func generateVerificationCode(length int) string {
	buffer := make([]byte, length)
	r := rand.Reader
	_, _ = io.ReadFull(r, buffer)
	code := ""
	for _, b := range buffer {
		code += strconv.Itoa(int(b) % 10) // 将字节转换为数字并取模10
	}
	return code[:length]
}

func (s *Service) VerificationCodeByPhone(phone, verificationCode, kind string) error {
	isValid, err := s.CheckVerificationCode(phone, verificationCode, kind)
	if !isValid {
		return err
	}
	return nil
}

func (s *Service) CheckVerificationCode(phone, verificationCode, kind string) (bool, error) {
	storedCode, err := s.repo.GetVerificationCodeByPhone(phone, kind)
	if err != nil {
		return false, err
	}
	//fmt.Println(err)
	if verificationCode == storedCode.Code {
		_ = s.repo.UpdateVerificationCode(storedCode.ID)
		return true, nil
	} else {
		return false, ErrInvalidVerificationCode
	}

}

func (s *Service) UpdateUserByUserID(uid uint32, u *UpdateUser) error {
	hashedPassword, err := s.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	u.ID = uid
	err = s.repo.UpdateUserByUserID(u)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserByID(id uint32) (*User, error) {
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *Service) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// todo 可升级为黑名单检测
func (s *Service) CheckUserExist(uid uint32) (bool, error) {
	exist, err := s.repo.CheckUserExistByID(uid)
	//fmt.Println(exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// This service method creates a new shipping address for the user with the given ID
func (s *Service) CreateShippingAddress(sp *CreateShippingAddress) (uint32, error) {
	isActive := sp.IsActive
	//当用户设置默认地址 清除原先的默认地址
	if isActive == true {
		err := s.repo.UpdateShippingAddressIsActiveByUID(sp.UserID, false)
		if err != nil {
			return 0, err
		}
	}
	//todo 设置收货地址数量上限
	id, err := s.repo.CreateShippingAddress(sp)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// This service method updates the address of a user with the given ID
func (s *Service) UpdateUserAddressByID(sp *ShippingAddress) error {
	isActive := sp.IsActive
	if isActive == true {
		err := s.repo.UpdateShippingAddressIsActiveByUID(sp.UserID, false)
		if err != nil {
			return err
		}
	}
	num, err := s.repo.UpdateUserAddressByID(sp)
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("目标不存在")
	}
	return nil
}

// This service method retrieves all shipping addresses for the user with the given ID
func (s *Service) GetUserAddresses(userID uint32) ([]*ShippingAddress, error) {
	addresses, err := s.repo.GetUserAddresses(userID)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// This service method retrieves the active shipping address for the user with the given ID
func (s *Service) GetActiveUserAddress(userID uint32) (*ShippingAddress, error) {
	address, err := s.repo.GetActiveUserAddress(userID)
	if err != nil {
		return nil, err
	}
	return address, nil
}

// This service method deletes the shipping address with the given ID for the user with the given ID
func (s *Service) DeleteUserAddressByID(userID, addressID uint32) error {
	num, err := s.repo.DeleteUserAddressByID(userID, addressID)
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("收货地址不存在")
	}
	return nil
}

func (s *Service) LogUserOperation(userId uint32, ip, operation, details string) error {
	// 创建日志对象
	userLog := UserLog{
		UserID:    userId,
		IP:        ip,
		Operation: operation,
		InitTime:  time.Now().Unix(),
		Details:   details,
	}
	_, err := s.repo.CreateUserLog(&userLog)
	if err != nil {
		return err
	}
	return nil
}
