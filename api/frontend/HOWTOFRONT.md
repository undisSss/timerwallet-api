HOWTO: Frontend (TimerWallet Mini App)

Overview
- Frontend is in `api/frontend` built with Vite + React + TypeScript.
- Client-side routes: `/`, `/game`, `/profile`.
- Auth flow: Telegram Mini App passes `initData` to the frontend (`window.Telegram.WebApp.initData`). Frontend posts `initData` to `POST /api/auth/telegram` (see `src/services/auth.ts`). Server validates `initData`, creates/loads user, issues tokens and sets http-only Secure cookies. Frontend then calls `GET /api/me` to fetch current user.

Files to edit
- Pages: `src/pages/Home.tsx`, `src/pages/Game.tsx`, `src/pages/Profile.tsx`.
- Shared components: `src/components/Header.tsx`.
- Auth service: `src/services/auth.ts`.
- Types: `src/types/user.ts`.
- Styles: add `src/styles.css` as needed.

Dev workflow
1. Install dependencies:

```bash
cd api/frontend
npm install
```

2. Run dev server (Vite):

```bash
npm run dev
```

- Dev server runs on `http://localhost:3000` and proxies `/api` to `http://localhost:8080` (see `vite.config.ts`).
- When testing in Telegram Mini App you can build/serve production and expose your machine via ngrok or use a server.

Auth flow details (initData verification)
1. Telegram Mini App includes a query string `initData` that contains fields like `user`, `auth_date`, `hash`. The server must verify `hash` according to Telegram docs:
   - Build a data-check-string by sorting all fields except `hash` in alphabetical order and joining as `key=value` with newlines.
   - Compute HMAC-SHA256 using a secret key derived from your `bot_token`:
     - Secret key = SHA256(bot_token)
     - HMAC = HMAC-SHA256(secret_key, data-check-string)
   - Compare hex (or base16) of HMAC with the `hash` field (use constant-time comparison).
2. If verification passes, extract `user.id` and use it to find or create a user in DB.
3. Generate access and refresh tokens server-side (JWT or other) and set them as `HttpOnly; Secure; SameSite=None` cookies.

Frontend responsibilities
- Only send `initData` to the server: `POST /api/auth/telegram` (implemented in `src/services/auth.ts`). Include `credentials: 'include'` so cookies are saved.
- After auth, call `GET /api/me` to retrieve current user. If the server returns 401, redirect to home or ask user to re-authenticate.
- The frontend should never try to read tokens from cookies (they are httpOnly).

Working with Telegram WebApp locally
- When opened inside Telegram the WebApp injects `window.Telegram.WebApp.initData`.
- For local testing you can mock `window.Telegram.WebApp.initData` in browser console before clicking Login.
- Alternatively, capture the `initData` from an actual Mini App launch and paste it into a browser input implemented for testing only.

Docker / Production
- `Dockerfile` builds the app and serves via nginx. nginx proxies `/api` to the backend server (see `nginx.conf`).

- Important (docker-compose): when running the full stack with `docker-compose`, the frontend's `nginx.conf` is configured to proxy `/api` to `http://app:8080/` — that is, the `app` service name inside the compose network. This allows the frontend container to reach the backend by service name and works correctly when you start the entire stack with `docker-compose up`.

- Local / other setups: if you run the frontend container outside the same compose network (or test nginx locally), `nginx.conf` may need `proxy_pass http://host.docker.internal:8080/` or another backend host/IP where your API is reachable. Adjust `nginx.conf` accordingly before building the image for that environment.

Where to change logic
- To change auth endpoint or add headers, edit `src/services/auth.ts`.
- To change the place where initData is read (for example, read from query param), change `src/pages/Home.tsx` login handler.

Testing
- Unit testing is not preconfigured; add Jest/Testing Library if needed.

Next steps you may want
- Add better UI and styles.
- Add token refresh flow via a silent endpoint `POST /api/auth/refresh`.
- Add centralized auth context (React Context) to store user state.

Contact
- If server endpoints differ from `/api/auth/telegram` or `/api/me`, update `vite.config.ts` proxy and `src/services/auth.ts` accordingly.
