package v2

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zkep/my-geektime/internal/global"
	"github.com/zkep/my-geektime/internal/mcptoken"
	"github.com/zkep/my-geektime/internal/model"
	"github.com/zkep/my-geektime/internal/types/user"
)

type MCPToken struct{}

func NewMCPToken() *MCPToken {
	return &MCPToken{}
}

func (m *MCPToken) List(c *gin.Context) {
	uid := c.GetString(global.Identity)
	if uid == "" {
		global.FAIL(c, "fail.msg", "no auth")
		return
	}
	var rows []*model.MCPToken
	if err := global.DB.
		Where("uid = ? AND deleted_at = ?", uid, 0).
		Order("id DESC").
		Find(&rows).Error; err != nil {
		global.FAIL(c, "fail.msg", err.Error())
		return
	}
	ret := user.MCPTokenListResponse{Rows: make([]user.MCPToken, 0, len(rows))}
	for _, row := range rows {
		ret.Rows = append(ret.Rows, toMCPTokenItem(row))
	}
	global.OK(c, ret)
}

func (m *MCPToken) Create(c *gin.Context) {
	var req user.MCPTokenCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FAIL(c, "fail.msg", err.Error())
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		global.FAIL(c, "fail.msg", "名称不能为空")
		return
	}
	uid := c.GetString(global.Identity)
	if uid == "" {
		global.FAIL(c, "fail.msg", "no auth")
		return
	}
	plain, hash, prefix, err := mcptoken.Generate()
	if err != nil {
		global.FAIL(c, "fail.msg", err.Error())
		return
	}
	now := time.Now().Unix()
	row := &model.MCPToken{
		Uid:         uid,
		Name:        req.Name,
		TokenHash:   hash,
		TokenPrefix: prefix,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err = global.DB.Create(row).Error; err != nil {
		global.FAIL(c, "fail.msg", err.Error())
		return
	}
	global.OK(c, user.MCPTokenCreateResponse{Token: plain, Item: toMCPTokenItem(row)})
}

func (m *MCPToken) Delete(c *gin.Context) {
	var req user.MCPTokenDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.FAIL(c, "fail.msg", err.Error())
		return
	}
	uid := c.GetString(global.Identity)
	if uid == "" {
		global.FAIL(c, "fail.msg", "no auth")
		return
	}
	now := time.Now().Unix()
	tx := global.DB.Model(&model.MCPToken{}).
		Where("id = ? AND uid = ? AND deleted_at = ?", req.Id, uid, 0).
		Updates(map[string]any{"deleted_at": now, "updated_at": now})
	if tx.Error != nil {
		global.FAIL(c, "fail.msg", tx.Error.Error())
		return
	}
	if tx.RowsAffected == 0 {
		global.FAIL(c, "fail.msg", "Token 不存在")
		return
	}
	global.OK(c, nil)
}

func toMCPTokenItem(row *model.MCPToken) user.MCPToken {
	return user.MCPToken{
		Id:          row.Id,
		Name:        row.Name,
		TokenPrefix: row.TokenPrefix,
		LastUsedAt:  row.LastUsedAt,
		CreatedAt:   row.CreatedAt,
	}
}
