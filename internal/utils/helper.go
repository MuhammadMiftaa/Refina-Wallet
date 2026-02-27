package utils

import (
	"errors"
	"time"

	"refina-wallet/config/env"
	"refina-wallet/internal/types/dto"
	"refina-wallet/internal/types/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func ConvertToResponseType(data any) any {
	switch v := data.(type) {
	case model.Wallets:
		return dto.WalletsResponse{
			ID:             v.ID.String(),
			UserID:         v.UserID.String(),
			WalletTypeID:   v.WalletTypeID.String(),
			WalletType:     string(v.WalletType.Type),
			WalletTypeName: v.WalletType.Name,
			Name:           v.Name,
			Number:         v.Number,
			Balance:        v.Balance,
		}
	case model.WalletTypes:
		return dto.WalletTypesResponse{
			ID:          v.ID.String(),
			Name:        v.Name,
			Type:        dto.WalletType(v.Type),
			Description: v.Description,
		}
	default:
		return nil
	}
}

func VerifyToken(jwtToken string) (dto.UserData, error) {
	token, _ := jwt.Parse(jwtToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("parsing token error occured")
		}
		return []byte(env.Cfg.Server.JWTSecretKey), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return dto.UserData{}, errors.New("token is invalid")
	}

	return dto.UserData{
		ID:       claims["id"].(string),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}, nil
}

func ParseUUID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedID, nil
}

func Ms(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / 1e6
}
