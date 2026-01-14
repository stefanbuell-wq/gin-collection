# Gin Collection SaaS - Frontend

Modern React + TypeScript frontend for the Gin Collection SaaS platform.

## Tech Stack

- **Framework:** React 18 + TypeScript
- **Build Tool:** Vite
- **Routing:** React Router v6
- **State Management:** Zustand
- **Styling:** Tailwind CSS
- **HTTP Client:** Axios
- **Icons:** Lucide React
- **PWA:** Vite PWA Plugin + Workbox

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn/pnpm

### Installation

```bash
# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Start development server
npm run dev
```

The app will be available at `http://localhost:3000`

### Build for Production

```bash
npm run build
npm run preview  # Preview production build locally
```

## Project Structure

```
src/
├── api/              # API client and service functions
│   ├── client.ts     # Axios instance with interceptors
│   └── services.ts   # All API endpoints
├── components/       # Reusable React components
│   ├── Layout.tsx    # Main app layout with navigation
│   └── ProtectedRoute.tsx  # Auth guard
├── pages/            # Page components (lazy loaded)
│   ├── Login.tsx
│   ├── Register.tsx
│   ├── Dashboard.tsx
│   ├── GinList.tsx
│   └── ...
├── routes/           # React Router configuration
│   └── index.tsx
├── stores/           # Zustand stores
│   ├── authStore.ts  # Authentication state
│   └── ginStore.ts   # Gin collection state
├── types/            # TypeScript type definitions
│   └── index.ts
├── App.tsx           # Root component
├── main.tsx          # Entry point
└── index.css         # Global styles + Tailwind
```

## Features

### Authentication
- JWT-based authentication
- Automatic token refresh
- Subdomain-based tenant detection
- Protected routes

### State Management
- Zustand for global state
- Persistent auth state (localStorage)
- Optimistic updates

### API Integration
- Axios instance with request/response interceptors
- Automatic tenant header injection
- Error handling with retry logic
- Upgrade prompts for tier-limited features

### PWA Support
- Service Worker for offline functionality
- App manifest for install prompts
- Network-first caching strategy for API
- Background sync (coming soon)

## Environment Variables

Create a `.env` file based on `.env.example`:

```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_ENV=development
```

## Development

### Code Style
- TypeScript strict mode enabled
- ESLint for code quality
- React Hooks rules enforced

### Testing (Coming Soon)
- Vitest for unit tests
- React Testing Library for component tests
- Cypress for E2E tests

## Deployment

### Docker

```bash
# Build production image
docker build -t gin-collection-frontend .

# Run
docker run -p 80:80 gin-collection-frontend
```

### Nginx

The production build can be served with any static file server. Example nginx config:

```nginx
server {
  listen 80;
  server_name _;

  root /usr/share/nginx/html;
  index index.html;

  # Serve React app
  location / {
    try_files $uri $uri/ /index.html;
  }

  # Proxy API requests
  location /api {
    proxy_pass http://api:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

## Contributing

See main project README for contribution guidelines.
