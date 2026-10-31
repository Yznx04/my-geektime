import React, { useEffect, useState } from 'react'
import { Copy, KeyRound, Plus, Trash2 } from 'lucide-react'
import { Button, Card, Input, Modal, Spinner, Table } from '@/components/ui'
import { useToast } from '@/components/ui/Toast'
import {
  createMCPToken,
  deleteMCPToken,
  getMCPTokens,
  MCPTokenItem,
} from '@/api/mcpToken'

const formatTime = (value: number) =>
  value > 0 ? new Date(value * 1000).toLocaleString('zh-CN') : '-'

export const McpToken: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [items, setItems] = useState<MCPTokenItem[]>([])
  const [showCreate, setShowCreate] = useState(false)
  const [name, setName] = useState('')
  const [creating, setCreating] = useState(false)
  const [newToken, setNewToken] = useState('')
  const [showToken, setShowToken] = useState(false)
  const [deleteItem, setDeleteItem] = useState<MCPTokenItem | null>(null)
  const { addToast } = useToast()

  const endpoint = `${window.location.origin}/mcp`

  const loadData = async () => {
    setLoading(true)
    try {
      const res = await getMCPTokens()
      setItems(res.rows || [])
    } catch (error) {
      console.error('Failed to load MCP tokens', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  const handleCreate = async () => {
    if (!name.trim()) {
      addToast('请输入 Token 名称', 'warning')
      return
    }
    setCreating(true)
    try {
      const res = await createMCPToken({ name: name.trim() })
      setShowCreate(false)
      setName('')
      setNewToken(res.token)
      setShowToken(true)
      loadData()
    } catch (error: any) {
      const errorMsg = error?.response?.msg || error?.message || '生成失败'
      addToast(errorMsg, 'error')
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteItem) return
    try {
      await deleteMCPToken({ id: deleteItem.id })
      addToast('删除成功', 'success')
      setDeleteItem(null)
      loadData()
    } catch (error: any) {
      const errorMsg = error?.response?.msg || error?.message || '删除失败'
      addToast(errorMsg, 'error')
    }
  }

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(newToken)
      addToast('已复制到剪贴板', 'success')
    } catch {
      addToast('复制失败，请手动选择复制', 'error')
    }
  }

  const columns = [
    {
      key: 'name',
      label: '名称',
      render: (item: MCPTokenItem) => (
        <span className="font-medium">{item.name}</span>
      ),
    },
    {
      key: 'token_prefix',
      label: 'Token',
      render: (item: MCPTokenItem) => (
        <code className="font-mono text-sm text-gray-600">
          {item.token_prefix}...
        </code>
      ),
    },
    {
      key: 'created_at',
      label: '创建时间',
      render: (item: MCPTokenItem) => formatTime(item.created_at),
    },
    {
      key: 'last_used_at',
      label: '最近使用',
      render: (item: MCPTokenItem) =>
        item.last_used_at > 0 ? formatTime(item.last_used_at) : '从未使用',
    },
    {
      key: 'action',
      label: '操作',
      render: (item: MCPTokenItem) => (
        <Button
          variant="danger"
          size="sm"
          onClick={() => setDeleteItem(item)}
        >
          <Trash2 size={14} className="mr-1" />
          删除
        </Button>
      ),
    },
  ]

  return (
    <div>
      <Card
        header="MCP Token 管理"
        actions={
          <Button onClick={() => setShowCreate(true)}>
            <Plus size={16} className="mr-1" />
            生成 Token
          </Button>
        }
      />

      <Card className="mt-4">
        <div className="p-5 space-y-3 text-sm text-gray-600">
          <div className="flex items-start gap-3 rounded-xl bg-blue-50 border border-blue-200 p-4 text-blue-800">
            <KeyRound size={18} className="mt-0.5 flex-shrink-0" />
            <div className="space-y-1">
              <p>这里的 Token 不会过期，适用于无法通过账号密码自动刷新 JWT 的 MCP 客户端。</p>
              <p>
                MCP 地址：<code className="font-mono">{endpoint}</code>
              </p>
              <p>
                请求头：<code className="font-mono">Authorization: Bearer &lt;Token&gt;</code>
              </p>
              <p>Token 仅在创建时显示一次，请及时保存；删除后立即失效。</p>
            </div>
          </div>
        </div>

        {loading ? (
          <div className="flex justify-center py-8">
            <Spinner />
          </div>
        ) : (
          <>
            <Table
              columns={columns}
              data={items}
              rowKey="id"
              className="table-db table-striped"
            />
            {items.length === 0 && (
              <div className="flex flex-col items-center justify-center py-20 text-gray-500">
                <div className="text-lg">还没有 Token</div>
                <div className="mt-1 text-sm">点击右上角“生成 Token”创建第一个</div>
              </div>
            )}
          </>
        )}
      </Card>

      <Modal
        isOpen={showCreate}
        onClose={() => {
          setShowCreate(false)
          setName('')
        }}
        title="生成 MCP Token"
        size="md"
      >
        <div className="space-y-4">
          <Input
            label="名称"
            placeholder="例如：Claude Desktop"
            value={name}
            onChange={(e) => setName(e.target.value)}
            maxLength={64}
          />
          <p className="text-sm text-gray-500">名称用于区分不同客户端，不影响 Token 权限。</p>
          <div className="flex gap-3">
            <Button
              variant="secondary"
              onClick={() => {
                setShowCreate(false)
                setName('')
              }}
              className="flex-1"
            >
              取消
            </Button>
            <Button
              variant="primary"
              onClick={handleCreate}
              disabled={creating}
              className="flex-1"
            >
              {creating ? '生成中...' : '确定'}
            </Button>
          </div>
        </div>
      </Modal>

      <Modal
        isOpen={showToken}
        onClose={() => {
          setShowToken(false)
          setNewToken('')
        }}
        title="Token 已生成"
        size="md"
      >
        <div className="space-y-4">
          <p className="text-sm text-amber-700 bg-amber-50 border border-amber-200 rounded-lg p-3">
            请立即复制并妥善保存，关闭窗口后将无法再次查看该 Token。
          </p>
          <textarea
            readOnly
            value={newToken}
            rows={4}
            className="w-full px-3 py-2 font-mono text-sm border border-gray-200 rounded-lg bg-gray-50 break-all"
            onFocus={(e) => e.target.select()}
          />
          <div className="flex gap-3">
            <Button variant="secondary" onClick={handleCopy} className="flex-1">
              <Copy size={16} className="mr-1" />
              复制 Token
            </Button>
            <Button
              variant="primary"
              onClick={() => {
                setShowToken(false)
                setNewToken('')
              }}
              className="flex-1"
            >
              我已保存
            </Button>
          </div>
        </div>
      </Modal>

      <Modal
        isOpen={!!deleteItem}
        onClose={() => setDeleteItem(null)}
        title="确认删除"
        size="sm"
      >
        <div className="space-y-4">
          <p className="text-sm text-gray-600">
            确定要删除 Token {deleteItem?.name} 吗？删除后使用该 Token 的客户端会立即失去访问权限。
          </p>
          <div className="flex gap-3">
            <Button
              variant="secondary"
              onClick={() => setDeleteItem(null)}
              className="flex-1"
            >
              取消
            </Button>
            <Button variant="danger" onClick={handleDelete} className="flex-1">
              确定删除
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}
