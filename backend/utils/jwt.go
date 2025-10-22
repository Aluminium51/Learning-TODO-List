// backend/utils/jwt.go
package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT creates a new JWT token for a given user ID
func GenerateJWT(userID uint) (string, error) {
	// ดึงค่า Secret Key จาก .env
	jwtSecret := os.Getenv("JWT_SECRET")

	// สร้าง "Claims" หรือข้อมูลที่จะเก็บไว้ในบัตรผ่าน
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // บัตรหมดอายุใน 24 ชั่วโมง

	// สร้าง Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// เซ็นลายเซ็นอิเล็กทรอนิกส์ลงบนบัตรด้วย Secret Key ของเรา
	return token.SignedString([]byte(jwtSecret))
}
