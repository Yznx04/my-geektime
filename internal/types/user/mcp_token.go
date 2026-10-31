package user

type MCPToken struct {
	// id
	Id int64 `json:"id"`
	// name
	Name string `json:"name"`
	// token_prefix
	TokenPrefix string `json:"token_prefix"`
	// last_used_at
	LastUsedAt int64 `json:"last_used_at"`
	// created_at
	CreatedAt int64 `json:"created_at"`
}

type MCPTokenListResponse struct {
	Rows []MCPToken `json:"rows"`
}

type MCPTokenCreateRequest struct {
	// name
	Name string `json:"name" binding:"required,max=64"`
}

type MCPTokenCreateResponse struct {
	// token 仅在创建时返回一次
	Token string   `json:"token"`
	Item  MCPToken `json:"item"`
}

type MCPTokenDeleteRequest struct {
	// id
	Id int64 `json:"id" binding:"required"`
}
