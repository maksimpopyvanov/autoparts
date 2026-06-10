import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Space, Popconfirm, message, Typography } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { getBrands, createBrand, updateBrand, deleteBrand } from '../api'

const { Title } = Typography

export default function Brands() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState(null)
  const [search, setSearch] = useState('')
  const [form] = Form.useForm()

  const load = async () => {
    setLoading(true)
    try {
      setData(await getBrands())
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
    form.setFieldsValue({ name: record.name })
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    try {
      if (editing) {
        await updateBrand(editing.id, values)
        message.success('Марка обновлена')
      } else {
        await createBrand(values)
        message.success('Марка добавлена')
      }
      setModalOpen(false)
      load()
    } catch (e) {
      message.error(e.message)
    }
  }

  const handleDelete = async (id) => {
    try {
      await deleteBrand(id)
      message.success('Марка удалена')
      load()
    } catch (e) {
      message.error(e.message)
    }
  }

  const filtered = data.filter(b =>
    b.name.toLowerCase().includes(search.toLowerCase())
  )

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: 'Название', dataIndex: 'name' },
    {
      title: 'Действия',
      width: 140,
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} size="small" onClick={() => openEdit(record)} />
          <Popconfirm title="Удалить марку?" onConfirm={() => handleDelete(record.id)} okText="Да" cancelText="Нет">
            <Button icon={<DeleteOutlined />} size="small" danger />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Марки автомобилей</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>Добавить</Button>
      </div>

      <Input.Search
        placeholder="Поиск по названию"
        value={search}
        onChange={e => setSearch(e.target.value)}
        style={{ marginBottom: 16, maxWidth: 320 }}
        allowClear
      />

      <Table
        dataSource={filtered}
        columns={columns}
        rowKey="id"
        loading={loading}
        pagination={{ pageSize: 10 }}
      />

      <Modal
        title={editing ? 'Редактировать марку' : 'Новая марка'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        okText="Сохранить"
        cancelText="Отмена"
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="Название" rules={[{ required: true, message: 'Введите название' }]}>
            <Input placeholder="Например: Toyota" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
