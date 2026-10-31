package model

// MCPToken 是用户为 MCP 客户端生成的长期访问令牌，仅保存哈希值。
type MCPToken struct {
	// id
	Id int64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement;comment:id"`
	// uid
	Uid string `json:"uid,omitempty" gorm:"index;size:128;comment:user uid"`
	// name
	Name string `json:"name,omitempty" gorm:"size:128;comment:token name"`
	// token_hash
	TokenHash string `json:"-" gorm:"unique;size:64;comment:sha256 hash of token"`
	// token_prefix
	TokenPrefix string `json:"token_prefix,omitempty" gorm:"size:16;comment:token prefix for display"`
	// last_used_at
	LastUsedAt int64 `json:"last_used_at,omitempty" gorm:"comment:last used at"`
	// created_at
	CreatedAt int64 `json:"created_at,omitempty" gorm:"comment:created at"`
	// updated_at
	UpdatedAt int64 `json:"updated_at,omitempty" gorm:"comment:updated at"`
	// deleted_at
	DeletedAt int64 `json:"deleted_at,omitempty" gorm:"index;comment:deleted at"`
}
