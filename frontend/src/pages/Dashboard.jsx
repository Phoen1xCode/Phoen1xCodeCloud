import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Card, CardBody, Table, TableHeader, TableColumn, TableBody, TableRow, TableCell } from '@heroui/react'
import { shareAPI } from '../services/api'
import { useDispatch } from 'react-redux'
import { logout } from '../store/authSlice'

export default function Dashboard() {
  const [shares, setShares] = useState([])
  const navigate = useNavigate()
  const dispatch = useDispatch()

  useEffect(() => {
    fetchShares()
  }, [])

  const fetchShares = async () => {
    try {
      const { data } = await shareAPI.listShares()
      setShares(data)
    } catch (err) {
      console.error('Failed to fetch shares')
    }
  }

  const handleDelete = async (code) => {
    if (!confirm('Delete this share?')) return
    try {
      await shareAPI.deleteShare(code)
      fetchShares()
    } catch (err) {
      alert('Delete failed')
    }
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">My Shares</h1>
          <div className="space-x-4">
            <Button onClick={() => navigate('/')}>Upload</Button>
            <Button onClick={() => dispatch(logout())}>Logout</Button>
          </div>
        </div>

        <Card>
          <CardBody>
            <Table>
              <TableHeader>
                <TableColumn>Code</TableColumn>
                <TableColumn>Type</TableColumn>
                <TableColumn>Name</TableColumn>
                <TableColumn>Downloads</TableColumn>
                <TableColumn>Created</TableColumn>
                <TableColumn>Actions</TableColumn>
              </TableHeader>
              <TableBody>
                {shares.map((share) => (
                  <TableRow key={share.id}>
                    <TableCell>{share.share_code}</TableCell>
                    <TableCell>{share.type}</TableCell>
                    <TableCell>{share.file_name || 'Text'}</TableCell>
                    <TableCell>{share.downloads}</TableCell>
                    <TableCell>{new Date(share.created_at).toLocaleDateString()}</TableCell>
                    <TableCell>
                      <Button size="sm" color="danger" onClick={() => handleDelete(share.share_code)}>
                        Delete
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardBody>
        </Card>
      </div>
    </div>
  )
}
