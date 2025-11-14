import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Card, CardBody, Button } from '@heroui/react'
import { shareAPI } from '../services/api'

export default function Share() {
  const { code } = useParams()
  const [share, setShare] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchShare = async () => {
      try {
        const { data } = await shareAPI.getShare(code)
        setShare(data)
      } catch (err) {
        setError(err.response?.data?.error || 'Share not found')
      } finally {
        setLoading(false)
      }
    }
    fetchShare()
  }, [code])

  const handleDownload = () => {
    window.location.href = `/api/share/${code}`
  }

  if (loading) return <div className="min-h-screen flex items-center justify-center">Loading...</div>
  if (error) return <div className="min-h-screen flex items-center justify-center text-red-500">{error}</div>

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto">
        <Card>
          <CardBody className="p-8">
            <h1 className="text-2xl font-bold mb-4">Shared Content</h1>
            {share.type === 'file' ? (
              <div>
                <p className="mb-4">File: {share.file_name}</p>
                <p className="mb-4">Size: {(share.file_size / 1024).toFixed(2)} KB</p>
                <Button color="primary" onClick={handleDownload}>
                  Download File
                </Button>
              </div>
            ) : (
              <div>
                <pre className="bg-gray-800 text-white p-4 rounded overflow-auto">
                  {share.text_content}
                </pre>
              </div>
            )}
            <p className="mt-4 text-sm text-gray-500">Downloads: {share.downloads}</p>
          </CardBody>
        </Card>
      </div>
    </div>
  )
}
