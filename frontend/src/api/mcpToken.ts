import request from '@/utils/request'

export interface MCPTokenItem {
  id: number
  name: string
  token_prefix: string
  last_used_at: number
  created_at: number
}

export const getMCPTokens = () => {
  return request.get<any, { rows: MCPTokenItem[] }>('/user/mcp/tokens')
}

export const createMCPToken = (params: { name: string }) => {
  return request.post<any, { token: string; item: MCPTokenItem }>('/user/mcp/tokens', params)
}

export const deleteMCPToken = (params: { id: number }) => {
  return request.post('/user/mcp/tokens/delete', params)
}
