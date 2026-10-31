package mcptoken

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/zkep/my-geektime/internal/global"
	"github.com/zkep/my-geektime/internal/model"
	"github.com/zkep/my-geektime/internal/types/user"
)

const (
	Prefix          = "mgk_"
	prefixShowChars = 12
	lastUsedRefresh = 60
)

// Generate 返回明文 Token、哈希和用于展示的前缀。
func Generate() (plain, hash, prefix string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", "", err
	}
	plain = Prefix + base64.RawURLEncoding.EncodeToString(buf)
	hash = Hash(plain)
	prefix = plain[:prefixShowChars]
	return plain, hash, prefix, nil
}

func Hash(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

// Authenticate 校验长期 Token，返回所属用户；用户被禁用或删除时一并失效。
func Authenticate(ctx context.Context, plain string) (uid string, roleId int32, err error) {
	var token model.MCPToken
	if err = global.DB.WithContext(ctx).
		Where("token_hash = ? AND deleted_at = ?", Hash(plain), 0).
		First(&token).Error; err != nil {
		return "", 0, err
	}
	var info model.User
	if err = global.DB.WithContext(ctx).
		Where("uid = ? AND status = ? AND deleted_at = ?", token.Uid, user.UserStatusActive, 0).
		First(&info).Error; err != nil {
		return "", 0, err
	}
	now := time.Now().Unix()
	if token.LastUsedAt == 0 || now-token.LastUsedAt >= lastUsedRefresh {
		global.DB.WithContext(ctx).Model(&model.MCPToken{}).
			Where("id = ?", token.Id).
			Update("last_used_at", now)
	}
	return info.Uid, info.RoleId, nil
}
