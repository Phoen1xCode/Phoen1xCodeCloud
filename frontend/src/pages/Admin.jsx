import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Card, CardBody } from '@heroui/react'
import { adminAPI } from '../services/api'
import { useDispatch } from 'react-redux'
import { logout } from '../store/authSlice'

export default function Admin() {
  const [stats, setStats] = useState(null)
  const navigate = useNavigate()
  const dispatch = useDispatch()

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      const { data } = await adminAPI.getStats()
      setStats(data)
    } catch (err) {
      console.error('Failed to fetch stats')
    }
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Admin Dashboard</h1>
          <div className="space-x-4">
            <Button onClick={() => navigate('/')}>Upload</Button>
            <Button onClick={() => navigate('/dashboard')}>My Shares</Button>
            <Button onClick={() => dispatch(logout())}>Logout</Button>
          </div>
        </div>

        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <Card>
              <CardBody className="p-6">
                <h3 className="text-lg font-semibold mb-2">Total Users</h3>
                <p className="text-3xl font-bold">{stats.users}</p>
              </CardBody>
            </Card>
            <Card>
              <CardBody className="p-6">
                <h3 className="text-lg font-semibold mb-2">Total Shares</h3>
                <p className="text-3xl font-bold">{stats.total_shares}</p>
                <p className="text-sm text-gray-500 mt-2">
                  Files: {stats.file_shares} | Text: {stats.text_shares}
                </p>
              </CardBody>
            </Card>
            <Card>
              <CardBody className="p-6">
                <h3 className="text-lg font-semibold mb-2">Storage Used</h3>
                <p className="text-3xl font-bold">
                  {(stats.total_file_size / 1024 / 1024).toFixed(2)} MB
                </p>
              </CardBody>
            </Card>
          </div>
        )}
      </div>
    </div>
  )
}
