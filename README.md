# Phoen1xCodeCloud

A modern file sharing platform built with React and Go, enabling users to share files and code snippets with simple share codes.

## Features

- **File & Text Sharing**: Upload files or share text/code snippets
- **Share Code System**: Generate unique codes for easy content sharing
- **User Authentication**: Secure registration and login with JWT
- **Admin Dashboard**: Monitor service statistics and manage content
- **Flexible Storage**: Support for local storage or S3-compatible services (Cloudflare R2, AWS S3)
- **Download Tracking**: Track how many times content has been accessed

## Tech Stack

### Frontend
- **Vite** - Fast build tool and dev server
- **React 18** - UI library with hooks
- **Redux Toolkit** - State management
- **React Router** - Client-side routing
- **TailwindCSS** - Utility-first CSS framework
- **HeroUI** - Modern React component library
- **Axios** - HTTP client

### Backend
- **Go 1.21+** - High-performance backend language
- **Gin** - Fast HTTP web framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Relational database
- **AWS SDK v2** - S3-compatible storage integration
- **JWT** - Secure authentication tokens
- **bcrypt** - Password hashing

### DevOps
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **Nginx** - Reverse proxy and static file serving

## Architecture

```
phoen1xcodecloud/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── config/          # Configuration management
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── middleware/      # Auth & CORS middleware
│   │   ├── models/          # Database models
│   │   └── services/        # Business logic
│   └── pkg/
│       ├── storage/         # Storage abstraction (local/S3)
│       └── utils/           # Helper functions
├── frontend/
│   └── src/
│       ├── components/      # Reusable UI components
│       ├── pages/           # Route pages
│       ├── services/        # API client
│       └── store/           # Redux state management
└── docker/                  # Docker configurations
```

## API Endpoints

### Authentication
- `POST /api/register` - Register new user
- `POST /api/login` - User login

### Shares (Public)
- `GET /api/share/:code` - Get shared content

### Shares (Authenticated)
- `POST /api/upload` - Upload file
- `POST /api/text` - Create text share
- `GET /api/shares` - List user's shares
- `DELETE /api/share/:code` - Delete share

### Admin (Admin Only)
- `GET /api/admin/stats` - Service statistics
- `GET /api/admin/shares` - List all shares
- `GET /api/admin/users` - List all users

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 20+ (for local development)

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/phoen1xcode/phoen1xcodecloud.git
cd phoen1xcodecloud
```

2. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start all services:
```bash
docker-compose up -d
```

4. Access the application:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

### Local Development

#### Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Run database migrations (ensure PostgreSQL is running)
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/phoen1xcloud?sslmode=disable"

# Start the server
go run cmd/server/main.go
```

#### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Backend server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `JWT_SECRET` | **Required**: Strong secret key (min 32 chars) | - |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins | `http://localhost:3000,http://localhost:5173` |
| `STORAGE_TYPE` | Storage backend (`local` or `s3`) | `local` |
| `LOCAL_STORAGE_PATH` | Local file storage path | `./uploads` |
| `S3_BUCKET` | S3 bucket name | - |
| `S3_REGION` | S3 region | `auto` |
| `S3_ENDPOINT` | S3 endpoint URL | - |
| `S3_ACCESS_KEY` | S3 access key | - |
| `S3_SECRET_KEY` | S3 secret key | - |

### Cloudflare R2 Configuration

To use Cloudflare R2 for storage:

```env
STORAGE_TYPE=s3
S3_BUCKET=your-bucket-name
S3_REGION=auto
S3_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
S3_ACCESS_KEY=your-r2-access-key
S3_SECRET_KEY=your-r2-secret-key
```

## Database Schema

### Users Table
- `id` - Primary key
- `username` - Unique username
- `email` - Unique email
- `password` - Hashed password
- `is_admin` - Admin flag
- `created_at`, `updated_at`, `deleted_at`

### Shares Table
- `id` - Primary key
- `share_code` - Unique 8-character code
- `user_id` - Foreign key to users
- `type` - Content type (`file` or `text`)
- `file_name`, `file_size`, `file_path` - File metadata
- `text_content` - Text/code content
- `downloads` - Download counter
- `expires_at` - Optional expiration
- `created_at`, `updated_at`, `deleted_at`

## Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT Authentication**: 24-hour token expiration with strong secret enforcement (min 32 chars)
- **CORS Protection**: Configurable allowed origins (no wildcards in production)
- **SQL Injection Prevention**: GORM parameterized queries
- **Admin Authorization**: Middleware-based access control
- **Path Traversal Protection**: File path sanitization and validation
- **Rate Limiting**: 100 requests per minute per IP address
- **File Upload Validation**: Size limits (100MB) and type restrictions
- **Error Handling**: All errors properly handled and logged

## Production Deployment

### Build for Production

```bash
# Build backend
cd backend
CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Build frontend
cd frontend
npm run build
```

### Docker Production Build

```bash
docker-compose -f docker-compose.yml up -d --build
```

### Security Recommendations

1. **Generate a strong JWT_SECRET**: `openssl rand -base64 32` (required, minimum 32 characters)
2. **Use HTTPS in production**: Set up SSL/TLS certificates
3. **Configure CORS origins**: Set `CORS_ALLOWED_ORIGINS` to your production domain
4. **Set up database backups**: Regular automated backups
5. **Use environment-specific configurations**: Separate configs for dev/staging/prod
6. **Monitor rate limiting**: Adjust limits based on your traffic patterns
7. **Add virus scanning**: Integrate with ClamAV or similar for uploaded files
8. **Enable logging**: Set up centralized logging and monitoring
9. **Regular updates**: Keep dependencies updated for security patches
10. **Security headers**: Add security headers in nginx/reverse proxy

## Development

### Adding New Features

1. Backend: Add handlers in `internal/handlers/`
2. Frontend: Add pages in `src/pages/` and components in `src/components/`
3. State: Add Redux slices in `src/store/`
4. API: Update `src/services/api.js`

### Code Style

- **Go**: Follow standard Go conventions, use `gofmt`
- **JavaScript**: Use ES6+ features, functional components with hooks
- **CSS**: Use TailwindCSS utility classes

## Troubleshooting

### Database Connection Issues
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# View logs
docker-compose logs postgres
```

### Frontend Build Errors
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Storage Issues
```bash
# Check upload directory permissions
chmod 755 ./uploads

# For S3, verify credentials and endpoint
```

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## Author

**Phoen1xCode** - [@Phoen1xCode](https://github.com/phoen1xcode)

## Acknowledgments

- Built with modern web technologies
- Inspired by file sharing services
- Community-driven development
