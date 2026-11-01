import React, { useState, useRef, useCallback, useEffect } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import { clsx } from 'clsx'
import {
  BookOpen,
  Monitor,
  Hammer,
  ChevronDown,
  ChevronRight,
  X,
} from 'lucide-react'
import { useAuthStore } from '@/store/auth'

interface MenuItem {
  label: string
  url?: string
  icon?: React.ReactNode
  children?: MenuItem[]
  visibleOn?: (roleId: number) => boolean
}

const menuItems: MenuItem[] = [
  {
    label: '我的课程',
    icon: <BookOpen size={18} />,
    children: [
      { label: '课程列表', url: '/task/list' },
    ],
  },
  {
    label: '我的收藏',
    icon: <BookOpen size={18} />,
    children: [
      { label: '收藏列表', url: '/collect/list' },
    ],
  },
  {
    label: '极客课程',
    icon: <Monitor size={18} />,
    visibleOn: (roleId) => roleId === 1,
    children: [
      { label: '体系/公开/线下大会', url: '/product/pvip' },
      { label: '每日一课', url: '/product/lesson' },
      { label: '大厂案例', url: '/product/case' },
    ],
  },
  {
    label: '系统设置',
    icon: <Hammer size={18} />,
    children: [
      { label: '系统配置', url: '/setting', visibleOn: (roleId) => roleId === 1 },
      { label: '用户管理', url: '/user/list', visibleOn: (roleId) => roleId === 1 },
      { label: 'MCP Token 管理', url: '/mcp/token' },
    ],
  },
]

interface DrawerProps {
  isOpen: boolean
  onClose: () => void
  openMenus: Set<string>
  onToggleMenu: (label: string) => void
}

export const Drawer: React.FC<DrawerProps> = ({ isOpen, onClose, openMenus, onToggleMenu }) => {
  const user = useAuthStore((state) => state.user)
  const navigate = useNavigate()
  const [drawerWidth, setDrawerWidth] = useState<number>(288) // w-72 = 288px
  const isResizing = useRef(false)
  const startX = useRef(0)
  const startWidth = useRef(0)

  const isVisible = (item: MenuItem) => {
    if (item.visibleOn && user?.role_id) {
      return item.visibleOn(user.role_id)
    }
    return true
  }

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    isResizing.current = true
    startX.current = e.clientX
    startWidth.current = drawerWidth

    const handleMouseMove = (moveEvent: MouseEvent) => {
      if (!isResizing.current) return
      const deltaX = moveEvent.clientX - startX.current
      const newWidth = startWidth.current + deltaX
      const minWidth = 240
      const maxWidth = window.innerWidth * 0.8
      setDrawerWidth(Math.min(maxWidth, Math.max(minWidth, newWidth)))
    }

    const handleMouseUp = () => {
      isResizing.current = false
      document.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseup', handleMouseUp)
      document.body.style.cursor = ''
      document.body.style.userSelect = ''
    }

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
  }, [drawerWidth])

  useEffect(() => {
    const handleResize = () => {
      const newMaxWidth = window.innerWidth * 0.8
      if (drawerWidth > newMaxWidth) {
        setDrawerWidth(newMaxWidth)
      }
    }
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [drawerWidth])

  return (
    <>
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40 lg:hidden"
          onClick={onClose}
        />
      )}
      <aside
        className={clsx(
          'fixed inset-y-0 left-0 z-50 bg-white/95 backdrop-blur-xl transform transition-transform duration-300 ease-in-out lg:hidden',
          isOpen ? 'translate-x-0' : '-translate-x-full'
        )}
        style={{ width: drawerWidth }}
      >
        <div
          className="absolute right-0 top-0 bottom-0 w-4 cursor-grab hover:bg-primary-500/20 transition-colors group active:cursor-grabbing"
          onMouseDown={handleMouseDown}
        >
          <div className="absolute right-1 top-1/2 -translate-y-1/2 w-1 h-8 bg-gray-300 group-hover:bg-primary-500 rounded-full transition-colors" />
        </div>
        <div className="p-5 border-b border-gray-100 flex items-center justify-between">
          <button
            onClick={() => {
              navigate('/')
              onClose()
            }}
            className="flex items-center gap-3 hover:opacity-80 transition-opacity"
          >
            <div className="p-2.5 bg-primary-100 rounded-xl">
              <BookOpen className="w-5 h-5 text-primary-600" />
            </div>
            <h1 className="text-xl font-bold text-primary-700">
              我的极客时间
            </h1>
          </button>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X size={20} className="text-gray-500" />
          </button>
        </div>
        <nav className="p-3 overflow-y-auto h-[calc(100vh-80px)]">
          {menuItems.filter(isVisible).map((item) => (
            <div key={item.label} className="mb-1">
              {item.children ? (
                <>
                  <button
                    onClick={() => onToggleMenu(item.label)}
                    className="w-full flex items-center justify-between px-3 py-2.5 text-left rounded-xl hover:bg-gray-100 transition-colors text-gray-700"
                  >
                    <div className="flex items-center gap-3">
                      <span className="text-primary-500">{item.icon}</span>
                      <span className="font-medium">{item.label}</span>
                    </div>
                    {openMenus.has(item.label) ? (
                      <ChevronDown size={16} className="text-gray-400" />
                    ) : (
                      <ChevronRight size={16} className="text-gray-400" />
                    )}
                  </button>
                  {openMenus.has(item.label) && (
                    <div className="ml-5 mt-1">
                      {item.children.filter(isVisible).map((child) => (
                        <NavLink
                          key={child.url}
                          to={child.url!}
                          onClick={onClose}
                          className={({ isActive }) =>
                            clsx(
                              'block px-3 py-2.5 rounded-xl text-sm transition-colors',
                              isActive
                                ? 'bg-primary-50 text-primary-700 font-medium'
                                : 'text-gray-600 hover:bg-gray-100'
                            )
                          }
                        >
                          {child.label}
                        </NavLink>
                      ))}
                    </div>
                  )}
                </>
              ) : (
                <NavLink
                  to={item.url!}
                  onClick={onClose}
                  className={({ isActive }) =>
                    clsx(
                      'flex items-center gap-3 px-3 py-2.5 rounded-xl transition-colors',
                      isActive
                        ? 'bg-primary-50 text-primary-700 font-medium'
                        : 'text-gray-600 hover:bg-gray-100'
                    )
                  }
                >
                  <span className="text-primary-500">{item.icon}</span>
                  <span className="font-medium">{item.label}</span>
                </NavLink>
              )}
            </div>
          ))}
        </nav>
      </aside>
    </>
  )
}
