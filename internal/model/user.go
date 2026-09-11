package model

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils/random"
	"github.com/OpenListTeam/go-cache"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pkg/errors"
)

const (
	GENERAL = iota
	GUEST   // only one exists
	ADMIN
)

const (
	StaticHashSalt = "https://github.com/alist-org/alist"

	InvalidUsernameOrPassword = "Invalid username or password"
	Invalid2FACode            = "Invalid 2FA code"
	TooManyAttempts           = "Too many unsuccessful sign-in attempts have been made using an incorrect username or password, Try again later."
	GuestCannotUpdateProfile  = "Guest user can not update profile"
	GuestCannotGenerate2FA    = "Guest user can not generate 2FA code"
)

var LoginCache = cache.NewMemCache[int]()

var (
	DefaultLockDuration   = time.Minute * 5
	DefaultMaxAuthRetries = 5
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`                      // unique key
	Username string `json:"username" gorm:"unique" binding:"required"` // username
	PwdHash  string `json:"-"`                                         // password hash
	PwdTS    int64  `json:"-"`                                         // password timestamp
	Salt     string `json:"-"`                                         // unique salt
	Password string `json:"password"`                                  // password
	BasePath string `json:"base_path"`                                 // base path
	Role     int    `json:"role"`                                      // user's role
	Disabled bool   `json:"disabled"`
	// Determine permissions by bit
	//   0:  can mkdir and upload
	//   1:  can remove
	//   2:  webdav read
	//   3:  webdav manage
	Permission int32  `json:"permission"`
	OtpSecret  string `json:"-"`
	Authn      string `gorm:"type:text" json:"-"`
}

func (u *User) IsGuest() bool {
	return u.Role == GUEST
}

func (u *User) IsAdmin() bool {
	return u.Role == ADMIN
}

func (u *User) ValidateRawPassword(password string) error {
	return u.ValidatePwdStaticHash(StaticHash(password))
}

func (u *User) ValidatePwdStaticHash(pwdStaticHash string) error {
	if pwdStaticHash == "" {
		return errors.WithStack(errs.EmptyPassword)
	}
	if u.PwdHash != HashPwd(pwdStaticHash, u.Salt) {
		return errors.WithStack(errs.WrongPassword)
	}
	return nil
}

func (u *User) SetPassword(pwd string) *User {
	u.Salt = random.String(16)
	u.PwdHash = TwoHashPwd(pwd, u.Salt)
	u.PwdTS = time.Now().Unix()
	return u
}

func CanWriteContent(permission int32) bool {
	return permission&1 == 1
}

func (u *User) CanWriteContent() bool {
	return CanWriteContent(u.Permission)
}

func CanRemove(permission int32) bool {
	return (permission>>1)&1 == 1
}

func (u *User) CanRemove() bool {
	return CanRemove(u.Permission)
}

func CanWebdavRead(permission int32) bool {
	return (permission>>2)&1 == 1
}

func (u *User) CanWebdavRead() bool {
	return CanWebdavRead(u.Permission)
}

func CanWebdavManage(permission int32) bool {
	return (permission>>3)&1 == 1
}

func (u *User) CanWebdavManage() bool {
	return CanWebdavManage(u.Permission)
}

func (u *User) JoinPath(reqPath string) (string, error) {
	return utils.JoinBasePath(u.BasePath, reqPath)
}

func StaticHash(password string) string {
	return utils.HashData(utils.SHA256, []byte(fmt.Sprintf("%s-%s", password, StaticHashSalt)))
}

func HashPwd(static string, salt string) string {
	return utils.HashData(utils.SHA256, []byte(fmt.Sprintf("%s-%s", static, salt)))
}

func TwoHashPwd(password string, salt string) string {
	return HashPwd(StaticHash(password), salt)
}

func (u *User) WebAuthnID() []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(u.ID))
	return bs
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.Username
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	var res []webauthn.Credential
	err := json.Unmarshal([]byte(u.Authn), &res)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func (u *User) WebAuthnIcon() string {
	return "https://res.oplist.org/logo/logo.svg"
}
