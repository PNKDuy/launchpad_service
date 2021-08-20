package token

import (
	"errors"
	"github.com/google/uuid"
	"launchpad_service/model"
	"time"
)

type Token struct {
	Id uuid.UUID `json:"id"`
	Index int `json:"index"`
	Name string `json:"name"`
	SymbolToken string `json:"symbol_token"`
	Icon string `json:"icon"`
	TotalSupply string `json:"total_supply"`
	Description string `json:"description"`
	Website string `json:"website"`
	Twitter string `json:"twitter"`
	Facebook string `json:"facebook"`
	Telegram string `json:"telegram"`
	Reddit string `json:"reddit"`
	CoinMarketCap string `json:"coin_market_cap"`
	CoinGecko string `json:"coin_gecko"`
	Address string `json:"address"`
	ChainId string `json:"chain_id,omitempty"`
	ChainName string `json:"chain_name,omitempty"`
	Locked bool `json:"locked"`
	BaseCrypto string `json:"base_crypto,omitempty"`
	Decimal int64 `json:"decimal,omitempty"`
	DecimalBase int64 `json:"decimal_base"`
	AddressBase string `json:"address_base"`
	LaunchPadAmount string `json:"launch_pad_amount"`
	LaunchPadPrice float64 `json:"launch_pad_price"`
	MaxBuy string `json:"max_buy"`
	MinBuy string `json:"min_buy"`
	TimeStart int64 `json:"time_start"`
	TimeEnd int64	`json:"time_end"`
	ImageBanner string `json:"image_banner"`
	IsDeleted bool `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Get() ([]Token, error) {
	var tokens []Token
	db, err := model.ConnectToPostgres()
	if err != nil {
		return nil, err
	}

	if err = db.Where("is_deleted = ?", false).Find(&tokens).Error; err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Close()
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

func GetById(id string) (Token, error) {
	var token Token
	db, err := model.ConnectToPostgres()
	if err != nil {
		return token, err
	}

	result := db.First(&token, "id = ?", id)
	if result.Error != nil {
		return token, result.Error
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Close()
	if err != nil {
		return token, err
	}

	return token, nil
}

func Create(token Token) (Token, error) {
	db, err := model.ConnectToPostgres()
	if err != nil {
		return token, err
	}

	isExisted, err := checkIfNameOrSymbolTokenExists(token)
	if err != nil {
		return token, err
	}

	if isExisted == true {
		return token, errors.New("token name or symbol token is already existed")
	}

	token.Id = uuid.New()
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()

	result := db.Create(&token)
	if result.Error != nil {
		return token, result.Error
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Close()
	if err != nil {
		return token, err
	}

	return token, nil
}

func (token Token)Update() (Token, error){
	db, err := model.ConnectToPostgres()
	if err != nil {
		return token, err
	}
    token.UpdatedAt = time.Now()
	result := db.Model(&token).Updates(token)
	if result.Error != nil {
		return token, err
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Close()
	if err != nil {
		return token, err
	}

	return token, nil
}

func (token Token)DeactivateToken() (bool, error) {
	db, err := model.ConnectToPostgres()
	if err != nil {
		return false, err
	}
	result := db.Model(&token).Update("is_deleted", true)
	if result.Error != nil {
		return false, err
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Close()
	if err != nil {
		return false, err
	}

	return true, nil
}

func checkIfNameOrSymbolTokenExists(token Token) (bool, error) {
	db, err := model.ConnectToPostgres()
	var tokens []Token
	if err != nil {
		return true, err
	}
	result := db.Where("is_deleted = false AND ( name = ? OR symbol_token = ?)", token.Name, token.SymbolToken).Find(&tokens)
	if result.RowsAffected != 0 {
		return true, nil
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Close()
	if err != nil {
		return true, err
	}

	return false, nil
}
