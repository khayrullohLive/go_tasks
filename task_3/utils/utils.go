package utils

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================
// JWT UTILS
// ============================================================

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Token yaratish
func GenerateToken(userID uint, username string, isAdmin bool, secret string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)), // 7 kun
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Token tekshirish
func ValidateToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("kutilmagan imzolash usuli: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("yaroqsiz token")
	}

	return claims, nil
}

// ============================================================
// PASSWORD UTILS
// ============================================================

// Parolni hash qilish
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Parolni tekshirish
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ============================================================
// SLUG UTILS
// ============================================================

// Slug yaratish (masalan: "Hello World" -> "hello-world")
func GenerateSlug(title string) string {
	// Kichik harflarga o'tkazish
	slug := strings.ToLower(title)

	// Faqat harf, raqam va bo'sh joylarni qoldirish
	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' {
			result.WriteRune('-')
		}
	}

	// Ketma-ket tirechalarni bitta qilish
	re := regexp.MustCompile(`-+`)
	cleaned := re.ReplaceAllString(result.String(), "-")

	// Bosh va oxirdagi tirechalarni olib tashlash
	cleaned = strings.Trim(cleaned, "-")

	// Unique qilish uchun timestamp qo'shish
	return fmt.Sprintf("%s-%d", cleaned, time.Now().UnixMilli())
}

// ============================================================
// O'QISH VAQTINI HISOBLASH
// ============================================================

// O'rtacha o'qish tezligi: 200 so'z/daqiqa
func CalculateReadingTime(content string) int {
	words := len(strings.Fields(content))
	minutes := math.Ceil(float64(words) / 200.0)
	if minutes < 1 {
		return 1
	}
	return int(minutes)
}

// ============================================================
// PAGINATION UTILS
// ============================================================

func GetOffset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}

func GetTotalPages(total int64, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

// ============================================================
// EXCERPT YARATISH
// ============================================================

// Matndan qisqa excerpt olish
func GenerateExcerpt(content string, maxLength int) string {
	if len(content) <= maxLength {
		return content
	}

	// So'z o'rtasida kesilmasin
	truncated := content[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}
