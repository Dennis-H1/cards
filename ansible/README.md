# Deploying karteikarten

Assumes the VPS already has Docker + the `docker compose` plugin installed,
the deploy SSH user is in the `docker` group (no `become` used below), and an
external Docker network named `proxy` already exists (the one your existing
nginx-proxy + acme-companion setup uses -- referenced as `nginx-proxy`
internally in `../docker-compose.yml`, mapped to the real `proxy` network via
`name:`). If any of that doesn't hold, adjust `deploy.yml` /
`../docker-compose.yml` accordingly.

`LETSENCRYPT_HOST`/`LETSENCRYPT_EMAIL` are set on the `web` service so
acme-companion requests a real cert for `VIRTUAL_HOST` instead of nginx-proxy
falling back to its self-signed default (which Cloudflare will reject in
Full Strict mode).

There are two ways to run `deploy.yml`: manually from your machine, or
automatically from GitHub Actions on every push to `main` (and on manual
`workflow_dispatch`). Both run the exact same playbook; they just supply
`mcp_api_key` and `virtual_host` differently.

## Manual (from your machine)

One-time setup:

1. `cp inventory.ini.example inventory.ini`, then edit it with the real VPS
   hostname and SSH user. `inventory.ini` is gitignored -- it holds your
   real VPS IP/hostname, which shouldn't sit in a public repo (it would let
   anyone bypass Cloudflare Zero Trust and hit the origin directly).
2. Create the vault file with your MCP API key, login credentials, real
   hostname, and the email Let's Encrypt should use for cert
   notices/registration:
   ```
   cp group_vars/karteikarten/vault.yml.example group_vars/karteikarten/vault.yml
   # edit mcp_api_key to a real random value (e.g. `openssl rand -hex 32`),
   # auth_username/auth_password to the login credentials for the app itself,
   # session_secret to another long random value (e.g. `openssl rand -hex 32`),
   # virtual_host to the real hostname you want nginx-proxy/Cloudflare to
   # route, and letsencrypt_email to a real address you control
   ansible-vault encrypt group_vars/karteikarten/vault.yml
   ```

Then, any time you want to deploy:

```
ansible-playbook deploy.yml --ask-vault-pass
```

## Automatic (GitHub Actions)

The `deploy` job in `.github/workflows/ci.yml` runs after the image is
pushed to GHCR, using an **inline** inventory built from secrets -- it
never reads `inventory.ini` or the local vault. It needs these repo secrets
(Settings -> Secrets and variables -> Actions):

| Secret | Value |
|---|---|
| `SSH_PRIVATE_KEY` | Private half of a dedicated deploy keypair. Generate one just for this (`ssh-keygen -t ed25519 -f deploy_key -C karteikarten-ci`), and add the public half to the VPS deploy user's `~/.ssh/authorized_keys`. Don't reuse your personal key. |
| `VPS_HOST` | The VPS hostname or IP. |
| `VPS_USER` | The SSH user to deploy as. |
| `MCP_API_KEY` | Same value you'd put in the vault for a manual deploy. |
| `AUTH_USERNAME` | Login username for the app's own login screen (separate from the MCP API key). |
| `AUTH_PASSWORD` | Login password for the app's own login screen. |
| `SESSION_SECRET` | Long random value used to sign session cookies, e.g. `openssl rand -hex 32`. Rotating it invalidates all existing sessions. |
| `VIRTUAL_HOST` | The public hostname nginx-proxy/Cloudflare routes to this app, e.g. `karteikarten.yourdomain.com`. Not sensitive in the security sense (it's a public-facing URL), but kept as a secret rather than committed so the repo doesn't publicly tie your personal domain to this project. |
| `LETSENCRYPT_EMAIL` | Email address acme-companion registers with Let's Encrypt for this cert (expiry/problem notices). |

The job runs `ssh-keyscan` against `VPS_HOST` to populate `known_hosts`
before connecting (trust-on-first-use -- fine for a personal single-VPS
setup; pin the host key yourself first if you want stronger guarantees).

Trigger it by pushing to `main`, or manually via the Actions tab -> CI ->
"Run workflow".

## One-time GHCR step

After the *first* push to `main`, two images are published --
`ghcr.io/dennis-h1/cards-backend` (the Go binary) and
`ghcr.io/dennis-h1/cards-web` (nginx serving the built frontend and
proxying to the backend) -- but both start out **private** by GHCR
default. Go to *each* package's GitHub
settings and set visibility to Public once (matches the "public" choice
made for this repo), otherwise both the VPS pull and any manual
`docker pull` will fail with 401/403.
