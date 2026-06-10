import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, Space, Popconfirm, message, Typography, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { getParts, createPart, updatePart, deletePart, getCategories, getBrands } from '../api'

const { Title } = Typography

export default function Parts() {
  const [data, setData] = useState([])
  const [categories, setCategories] = useState([])
  const [brands, setBrands] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState(null)
  const [search, setSearch] = useState('')
  const [filterCategory, setFilterCategory] = useState(null)
  const [filterBrand, setFilterBrand] = useState(null)
  const [form] = Form.useForm()

  const load = async () => {
    setLoading(true)
    try {
      const [parts, cats, brs] = await Promise.all([getParts(), getCategories(), getBrands()])
      setData(parts)
      setCategories(cats)
      setBrands(brs)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    setModalOpen(true)
  }

  const openEdit = (record) => {
    setEditing(record)
    form.setFieldsValue({
      name: record.name,
      article: record.article,
      description: record.description,
      category_id: record.category_id,
      brand_ids: (record.brands || []).map(b => b.id),
    })
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    const payload = {
      ...values,
      article: values.article || null,
      description: values.description || null,
      category_id: values.category_id || null,
      brand_ids: values.brand_ids || [],
    }
    try {
      if (editing) {
        await updatePart(editing.id, payload)
        message.success('Запчасть обновлена')
      } else {
        await createPart(payload)
        message.success('Запчасть добавлена')
      }
      setModalOpen(false)
      load()
    } catch (e) {
      message.error(e.message)
    }
  }

  const handleDelete = async (id) => {
    try {
      await deletePart(id)
      message.success('Запчасть удалена')
      load()
    } catch (e) {
      message.error(e.message)
    }
  }

  const filtered = data.filter(p => {
    const matchSearch =
      p.name.toLowerCase().includes(search.toLowerCase()) ||
      (p.article || '').toLowerCase().includes(search.toLowerCase())
    const matchCat = filterCategory ? p.category_id === filterCategory : true
    const matchBrand = filterBrand
      ? (p.brands || []).some(b => b.id === filterBrand)
      : true
    return matchSearch && matchCat && matchBrand
  })

  const columns = [
    { title: 'Название', dataIndex: 'name' },
    { title: 'Артикул', dataIndex: 'article', render: v => v || '—' },
    {
      title: 'Категория',
      dataIndex: 'category_name',
      render: v => v ? <Tag color="blue">{v}</Tag> : '—',
    },
    {
      title: 'Марки',
      dataIndex: 'brands',
      render: (brands) => (
        <Space size={4} wrap>
          {(brands || []).length === 0
            ? '—'
            : (brands || []).map(b => <Tag key={b.id} color="green">{b.name}</Tag>)
          }
        </Space>
      ),
    },
    { title: 'Описание', dataIndex: 'description', render: v => v || '—', ellipsis: true },
    {
      title: 'Действия',
      width: 140,
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} size="small" onClick={() => openEdit(record)} />
          <Popconfirm title="Удалить запчасть?" onConfirm={() => handleDelete(record.id)} okText="Да" cancelText="Нет">
            <Button icon={<DeleteOutlined />} size="small" danger />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Запчасти</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>Добавить</Button>
      </div>

      <Space wrap style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="Поиск по названию или артикулу"
          value={search}
          onChange={e => setSearch(e.target.value)}
          style={{ width: 280 }}
          allowClear
        />
        <Select
          placeholder="Категория"
          allowClear
          style={{ width: 180 }}
          value={filterCategory}
          onChange={setFilterCategory}
          options={categories.map(c => ({ value: c.id, label: c.name }))}
        />
        <Select
          placeholder="Марка"
          allowClear
          style={{ width: 160 }}
          value={filterBrand}
          onChange={setFilterBrand}
          options={brands.map(b => ({ value: b.id, label: b.name }))}
        />
      </Space>

      <Table
        dataSource={filtered}
        columns={columns}
        rowKey="id"
        loading={loading}
        pagination={{ pageSize: 10 }}
      />

      <Modal
        title={editing ? 'Редактировать запчасть' : 'Новая запчасть'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        okText="Сохранить"
        cancelText="Отмена"
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="Название" rules={[{ required: true, message: 'Введите название' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="article" label="Артикул">
            <Input />
          </Form.Item>
          <Form.Item name="category_id" label="Категория">
            <Select
              allowClear
              placeholder="Выберите категорию"
              options={categories.map(c => ({ value: c.id, label: c.name }))}
            />
          </Form.Item>
          <Form.Item name="brand_ids" label="Марки автомобилей">
            <Select
              mode="multiple"
              allowClear
              placeholder="Выберите марки"
              options={brands.map(b => ({ value: b.id, label: b.name }))}
            />
          </Form.Item>
          <Form.Item name="description" label="Описание">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
