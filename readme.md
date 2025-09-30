# OpenSSL UI - Production Grade Web Interface

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

## Pricing Plans

### Free Tier
- 50 operations/month
- Basic certificate operations
- Community support

### Pro Tier - $19/month
- 5,000 operations/month
- Advanced features
- Priority support
- API access

### Enterprise - $99/month
- Unlimited operations
- White-label solution
- SSO integration
- SLA guarantee
- 24/7 support

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