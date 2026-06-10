import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, InputNumber, Select, Input, message, Typography } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { getOutcome, createOutcome, getParts } from '../api'

const { Title } = Typography

export default function Outcome() {
  const [data, setData] = useState([])
  const [parts, setParts] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [search, setSearch] = useState('')
  const [form] = Form.useForm()

  const load = async () => {
    setLoading(true)
    try {
      const [out, p] = await Promise.all([getOutcome(), getParts()])
      setData(out)
      setParts(p)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [])

  const handleSubmit = async () => {
    const values = await form.validateFields()
    try {
      await createOutcome({
        ...values,
        date: values.date || null,
        comment: values.comment || null,
      })
      message.success('Расход добавлен')
      setModalOpen(false)
      form.resetFields()
      load()
    } catch (e) {
      message.error(e.message)
    }
  }

  const filtered = data.filter(o =>
    o.part_name.toLowerCase().includes(search.toLowerCase())
  )

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: 'Запчасть', dataIndex: 'part_name' },
    { title: 'Количество (шт.)', dataIndex: 'quantity' },
    { title: 'Дата', dataIndex: 'date' },
    { title: 'Комментарий', dataIndex: 'comment', render: v => v || '—', ellipsis: true },
  ]

  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Расходы</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { form.resetFields(); setModalOpen(true) }}>
          Добавить расход
        </Button>
      </div>

      <Input.Search
        placeholder="Поиск по запчасти"
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
        pagination={{ pageSize: 15 }}
      />

      <Modal
        title="Новый расход"
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        okText="Сохранить"
        cancelText="Отмена"
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="part_id" label="Запчасть" rules={[{ required: true, message: 'Выберите запчасть' }]}>
            <Select
              showSearch
              placeholder="Выберите запчасть"
              filterOption={(input, option) =>
                option.label.toLowerCase().includes(input.toLowerCase())
              }
              options={parts.map(p => ({ value: p.id, label: p.name + (p.article ? ` (${p.article})` : '') }))}
            />
          </Form.Item>
          <Form.Item name="quantity" label="Количество" rules={[{ required: true, message: 'Введите количество' }]}>
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="date" label="Дата (необязательно)">
            <Input type="date" />
          </Form.Item>
          <Form.Item name="comment" label="Комментарий">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
