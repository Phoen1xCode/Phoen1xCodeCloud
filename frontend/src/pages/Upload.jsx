import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Input, Textarea, Card, CardBody, Tabs, Tab } from '@heroui/react'
import { shareAPI } from '../services/api'
import { useDispatch } from 'react-redux'
import { logout } from '../store/authSlice'

export default function Upload() {
  const [file, setFile] = useState(null)
  const [text, setText] = useState('')
  const [shareCode, setShareCode] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const dispatch = useDispatch()

  const handleFileUpload = async (e) => {
    e.preventDefault()
    if (!file) return
    setLoading(true)
    try {
      const formData = new FormData()
      formData.append('file', file)
      const { data } = await shareAPI.uploadFile(formData)
      setShareCode(data.share_code)
    } catch (err) {
      alert('Upload failed')
    } finally {
      setLoading(false)
    }
  }

  const handleTextUpload = async (e) => {
    e.preventDefault()
    if (!text) return
    setLoading(true)
    try {
      const { data } = await shareAPI.createText({ content: text })
      setShareCode(data.share_code)
    } catch (err) {
      alert('Upload failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Phoen1xCodeCloud</h1>
          <div className="space-x-4">
            <Button onClick={() => navigate('/dashboard')}>My Shares</Button>
            <Button onClick={() => dispatch(logout())}>Logout</Button>
          </div>
        </div>

        <Card>
          <CardBody className="p-8">
            <Tabs>
              <Tab key="file" title="Upload File">
                <form onSubmit={handleFileUpload} className="space-y-4 mt-4">
                  <Input
                    type="file"
                    onChange={(e) => setFile(e.target.files[0])}
                    required
                  />
                  <Button type="submit" color="primary" isLoading={loading}>
                    Upload File
                  </Button>
                </form>
              </Tab>
              <Tab key="text" title="Share Text">
                <form onSubmit={handleTextUpload} className="space-y-4 mt-4">
                  <Textarea
                    label="Text/Code"
                    value={text}
                    onChange={(e) => setText(e.target.value)}
                    rows={10}
                    required
                  />
                  <Button type="submit" color="primary" isLoading={loading}>
                    Share Text
                  </Button>
                </form>
              </Tab>
            </Tabs>

            {shareCode && (
              <div className="mt-6 p-4 bg-green-100 rounded">
                <p className="font-bold">Share Code: {shareCode}</p>
                <p className="text-sm mt-2">
                  Share URL: {window.location.origin}/share/{shareCode}
                </p>
              </div>
            )}
          </CardBody>
        </Card>
      </div>
    </div>
  )
}
