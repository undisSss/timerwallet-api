# TimerWallet Frontend

Vite + React + TypeScript frontend for the TimerWallet Telegram Mini App.

Quick start:

1. Install dependencies:

```bash
cd api/frontend
npm install
```

2. Dev server:

```bash
npm run dev
```

3. Build for production:

```bash
npm run build
```

4. Preview production build:

```bash
npm run preview
```

Notes:
- The dev server proxies `/api` to `http://localhost:8080` (see `vite.config.ts`).
- When running the full stack with `docker-compose`, the production nginx in the frontend container is configured to proxy `/api` to `http://app:8080/` (service name `app` inside the compose network). If you run nginx outside of the compose network, you may need to change `nginx.conf` to use `host.docker.internal:8080` or a specific backend host/IP.
- The Telegram Mini App will provide `window.Telegram.WebApp.initData` when opened inside Telegram.
- The frontend sends `initData` to `POST /api/auth/telegram` with `credentials: 'include'`. The server should set httpOnly cookies for tokens.
