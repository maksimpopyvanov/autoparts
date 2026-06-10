import { useState } from 'react'
import { Layout, Menu, Typography } from 'antd'
import {
  TagsOutlined,
  ToolOutlined,
  InboxOutlined,
  ArrowDownOutlined,
  ArrowUpOutlined,
  CarOutlined,
} from '@ant-design/icons'
import Categories from './pages/Categories'
import Brands from './pages/Brands'
import Parts from './pages/Parts'
import Stock from './pages/Stock'
import Income from './pages/Income'
import Outcome from './pages/Outcome'

const { Sider, Content } = Layout
const { Title } = Typography

const PAGES = {
  parts: <Parts />,
  categories: <Categories />,
  brands: <Brands />,
  stock: <Stock />,
  income: <Income />,
  outcome: <Outcome />,
}

const menuItems = [
  { key: 'parts', icon: <ToolOutlined />, label: 'Запчасти' },
  { key: 'categories', icon: <TagsOutlined />, label: 'Категории' },
  { key: 'brands', icon: <CarOutlined />, label: 'Марки авто' },
  { key: 'stock', icon: <InboxOutlined />, label: 'Склад (остатки)' },
  { key: 'income', icon: <ArrowDownOutlined />, label: 'Приходы' },
  { key: 'outcome', icon: <ArrowUpOutlined />, label: 'Расходы' },
]

export default function App() {
  const [page, setPage] = useState('parts')

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider width={220} theme="dark">
        <div style={{ padding: '20px 16px 12px' }}>
          <Title level={5} style={{ color: '#fff', margin: 0 }}>Учёт автозапчастей</Title>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[page]}
          items={menuItems}
          onClick={({ key }) => setPage(key)}
        />
      </Sider>
      <Layout>
        <Content style={{ margin: 24, padding: 24, background: '#fff', borderRadius: 8, minHeight: 360 }}>
          {PAGES[page]}
        </Content>
      </Layout>
    </Layout>
  )
}
