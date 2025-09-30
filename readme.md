<img width="1920" height="2294" alt="homepage" src="https://github.com/user-attachments/assets/6994ca0e-0fbe-4bbe-a688-1b324fa91552" /># OpenSSL UI - Production Grade Web Interface

A comprehensive, production-ready web interface for OpenSSL operations with modern UI design and commercial features.

## Features

- 🔐 Complete Certificate Management
- 🔑 Advanced Key Generation & Management
- 🛡️ Encryption/Decryption Tools
- 🔍 Hash Functions & Integrity Verification
- 🌐 TLS/SSL Testing & Analysis
- 👥 User Management & Authentication
- 💳 Subscription-based Pricing
- 🏢 Enterprise Features
- 📊 Usage Analytics & Billing

## Technology Stack

### Backend
- **Go 1.21+** with Gin framework
- **PostgreSQL** for data persistence
- **Redis** for caching and rate limiting
- **JWT** authentication
- **Stripe** for payments

### Frontend
- **Next.js 14** with TypeScript
- **Tailwind CSS** + Shadcn/ui components
- **React Query** for API state management
- **Framer Motion** for animations

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- OpenSSL installed on system

### Development Setup

1. **Clone and setup**
   ```bash
   git clone <repository>
   cd web-openssl
   ```

2. **Start services with Docker**
   ```bash
   docker-compose up -d postgres redis
   ```

3. **Run backend**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

4. **Run frontend**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

5. **Access application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Docs: http://localhost:8080/swagger

## Screenshots
<img width="1920" height="2294" alt="homepage" src="https://github.com/user-attachments/assets/6ff55e45-6b0c-4341-86df-7a363f6cc681" />)
<img width="1920" height="911" alt="loginpage" src="https://github.com/user-attachments/assets/8a8624e9-c1e2-44b9-a20f-6aae0f461cf5" />
<img width="1920" height="911" alt="encription" src="https://github.com/user-attachments/assets/92683bee-4a8a-4c47-8434-91e25ca3d88a" />
<img width="1920" height="911" alt="certificate page" src="https://github.com/user-attachments/assets/8f655011-8039-46eb-9705-761ae9da761e" />
<img width="1920" height="911" alt="dashboard" src="https://github.com/user-attachments/assets/e55bf033-86f7-4445-8930-3071c04ed0ec" />

![Dashboard](<img width="1920" height="2294" alt="homepage" src="https://github.com/user-attachments/assets/6ff55e45-6b0c-4341-86df-7a363f6cc681" />)

## API Documentation

The API follows RESTful conventions with comprehensive OpenAPI documentation available at `/swagger` endpoint.

## Security

- Zero-knowledge architecture
- End-to-end encryption
- GDPR compliant
- SOC 2 Type II ready
- OWASP security practices

## License

Commercial license - Contact for pricing and enterprise options.

## Support

- Documentation: `/docs`
- Issues: GitHub Issues
- Enterprise: enterprise@company.com
