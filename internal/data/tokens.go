package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"sort"
	"time"

	"gorm.io/gorm"
)

var (
	ScopeAuthentication = "authentication"
	ScopeActivation     = "activation"
)

type Token struct {
	CoreModel
	Scope     string    `json:"scope" gorm:"not null"`
	Plaintext string    `json:"token" gorm:"-"`
	Hash      []byte    `json:"-" gorm:"not null"`
	Code      string    `json:"code" gorm:"-"`
	Expiry    time.Time `json:"expiry"`
	UserID    int64     `json:"user_id" gorm:"not null"`
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Scope:  scope,
		Expiry: time.Now().Add(ttl),
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	// example plain token: Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)

	// one way hash with no salt, user will send plain token...
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenModel struct {
	DB *gorm.DB
}

func (m TokenModel) Insert(token *Token) error {
	return m.DB.Create(&token).Error
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) ActivateUserAndDeleteToken(tokenPlaintext string, user *User) error {
	tx := m.DB.Begin()
	if err := tx.Model(user).Updates(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := m.DB.Where("user_id=? and scope=?", user.ID, ScopeActivation).Delete(Token{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m TokenModel) KeepLastFiveAuthTokens(userID int64) error {
	var tokens []Token

	if err := m.DB.Where("user_id=? and scope=?", userID, ScopeAuthentication).Find(&tokens).Error; err != nil {
		return err
	}

	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].CreatedAt.After(tokens[j].CreatedAt)
	})

	if len(tokens) > 5 {
		args := []interface{}{userID, ScopeAuthentication, tokens[4].CreatedAt}
		if err := m.DB.Where("user_id=? and scope=? and created_at < ?", args...).Delete(&Token{}).Error; err != nil {
			return err
		}
	}

	return nil
}
