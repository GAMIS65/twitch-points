A simple tool for making giveaways on Twitch, where viewers can buy entries with their channel points.

<img width="2540" height="1264" alt="showcase" src="https://github.com/user-attachments/assets/84e87142-e0ea-4742-864e-6fe02fe73c5a" />


## Project Structure

```
twitch-points/
├── cmd/
│   └── main.go              # Back-end entry point
├── internal/
│   ├── api/                 # API routes
│   ├── db/                  # SQLC generated code
│   ├── sql/                 # SQL migrations
│   ├── twitch/              # Twitch API integration
│   └── util/                # Utility functions
├── frontend/                
│   ├── src/
│   └── package.json
├── docker-compose.yml       # Multi-container configuration
├── backend.dockerfile       # Backend container build
├── db.dockerfile            # Database container build
├── migrate.dockerfile       # Database migration container
├── go.mod / go.sum          # Go dependencies
├── sqlc.yaml                # sqlc configuration
└── .env.example             # Environment variables template
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Twitch Affiliate Account
- (Optional) Node.js 18+ and Go 1.24+ for local development

### Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/GAMIS65/twitch-points.git
   cd twitch-points
   ```

2. **Configure environment variables:**
   ```bash
   cp .env.example .env
   ```

3. **Update `.env` with your Twitch credentials:**
   ```
   # Twitch OAuth
   TWITCH_CLIENT_ID=your_client_id
   TWITCH_CLIENT_SECRET=your_client_secret
   TWITCH_WEBHOOK_SECRET=your_webhook_secret
   TWITCH_WEBHOOK_URL=https://yourdomain.com/webhooks/twitch

   # Database
   DB_USER=postgres
   DB_PASSWORD=your_secure_password
   DB_NAME=twitch_points
   DB_HOST=db
   DB_PORT=5432

   # Application
   HOST=http://localhost:8080
   FRONTEND_URL=http://localhost:5173
   BACKEND_DOMAIN_NAME=http://localhost:8080
   SESSION_KEY=your_session_key
   ```

4. **Start the application:**
   ```bash
   docker-compose up
   ```

### Development

**Backend:**
```bash
cd cmd
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```
