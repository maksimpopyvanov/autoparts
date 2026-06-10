import { useEffect, useState } from 'react'
import { Table, Input, Space, Typography, Tag } from 'antd'
import { getStock } from '../api'

const { Title } = Typography

export default function Stock() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')

  useEffect(() => {
    setLoading(true)
    getStock().then(setData).finally(() => setLoading(false))
  }, [])

  const filtered = data.filter(s =>
    s.part_name.toLowerCase().includes(search.toLowerCase()) ||
    (s.article || '').toLowerCase().includes(search.toLowerCase()) ||
    (s.category_name || '').toLowerCase().includes(search.toLowerCase())
  )

  const columns = [
    { title: 'Запчасть', dataIndex: 'part_name' },
    { title: 'Артикул', dataIndex: 'article', render: v => v || '—' },
    {
      title: 'Категория',
      dataIndex: 'category_name',
      render: v => v ? <Tag color="blue">{v}</Tag> : '—',
    },
    {
      title: 'Остаток (шт.)',
      dataIndex: 'quantity',
      render: v => (
        <Tag color={v === 0 ? 'red' : v < 5 ? 'orange' : 'green'}>
          {v}
        </Tag>
      ),
    },
  ]

  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Остатки на складе</Title>
      </div>

      <Input.Search
        placeholder="Поиск по названию, артикулу или категории"
        value={search}
        onChange={e => setSearch(e.target.value)}
        style={{ marginBottom: 16, maxWidth: 400 }}
        allowClear
      />

      <Table
        dataSource={filtered}
        columns={columns}
        rowKey="part_id"
        loading={loading}
        pagination={{ pageSize: 15 }}
      />
    </>
  )
}
