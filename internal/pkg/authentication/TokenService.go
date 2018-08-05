/*
 * Copyright © 2018 Rasmus Hansen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package authentication

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/zlepper/welp/internal/pkg/consts"
	"github.com/zlepper/welp/internal/pkg/models"
	"time"
)

type TokenServiceArgs struct {
	SecretService models.SecretService
	Logger        models.Logger
}

func NewTokenService(args TokenServiceArgs) (*TokenService, error) {
	return &TokenService{
		TokenServiceArgs: args,
	}, nil
}

type TokenService struct {
	TokenServiceArgs
}

func (s *TokenService) GenerateToken(ctx context.Context, duration time.Duration, subject interface{}) (string, error) {
	claim, err := s.generateClaim(duration, subject)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	secret, err := s.SecretService.GetSigningSecret(ctx)
	if err != nil {
		return "", err
	}

	return token.SignedString(secret)
}

func (s *TokenService) generateClaim(duration time.Duration, subject interface{}) (*jwt.StandardClaims, error) {
	sub, err := json.Marshal(subject)
	if err != nil {
		return nil, err
	}

	exp := time.Now().Add(duration)

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &jwt.StandardClaims{
		ExpiresAt: exp.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    consts.Issuer,
		Id:        id.String(),
		Subject:   string(sub),
	}, nil
}
