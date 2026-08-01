package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey      []byte
	expiryMinutes  int
	issuer         string
	audience       string
	driverAudience string
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string, expiryMinutes int, issuer string, audience string, driverAudience string) *JWTManager {
	return &JWTManager{
		secretKey:      []byte(secret),
		expiryMinutes:  expiryMinutes,
		issuer:         issuer,
		audience:       audience,
		driverAudience: driverAudience,
	}
}

// Generate stamps a role-specific audience: driver tokens get driverAudience
// (what driver-request-handler/websocket-gateway's driver route expect),
// everything else gets the existing client/rider audience. role is already
// passed in by every caller today (application/driver/login.go: "driver",
// application/user/login.go: "rider"), so this needed no caller changes.
func (m *JWTManager) Generate(userID string, email string, role string) (string, error) {
	audience := m.audience
	if role == "driver" {
		audience = m.driverAudience
	}

	now := time.Now().UTC()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    m.issuer,
			Audience:  []string{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.expiryMinutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// Parse backs this service's own auth middleware, which gates both rider
// and driver routes through the same JWTManager (role-specific access is
// enforced separately, via Claims.Role in interfaces/http/middleware). Since
// Generate now stamps one of two different audiences depending on role, this
// can't use jwt.WithAudience(m.audience) as a parser option (that only ever
// accepts one fixed value) — it accepts a token whose audience is either of
// this service's two valid audiences, then leaves role-gating to the caller.
func (m *JWTManager) Parse(tokenValue string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	}, jwt.WithIssuer(m.issuer), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if !audienceContains(claims.RegisteredClaims.Audience, m.audience) && !audienceContains(claims.RegisteredClaims.Audience, m.driverAudience) {
		return nil, fmt.Errorf("token audience does not match %q or %q", m.audience, m.driverAudience)
	}

	return claims, nil
}

func audienceContains(audience jwt.ClaimStrings, value string) bool {
	for _, a := range audience {
		if a == value {
			return true
		}
	}
	return false
}
