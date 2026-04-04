#!/usr/bin/env bash
# Self-signed cert for testing only. Browsers still show "not private" (untrusted CA),
# but names must match or you get ERR_CERT_COMMON_NAME_INVALID.
#
# Includes *.clickly.cv so tunnel hosts like uveulbt.clickly.cv match the cert.
# Production: use Let's Encrypt wildcard (*.clickly.cv) via DNS-01 — HTTP-01 cannot issue wildcard.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SSL="$ROOT/nginx/ssl"
mkdir -p "$SSL"

if [[ -f "$SSL/fullchain.pem" && -f "$SSL/privkey.pem" && "${1:-}" != "-f" ]]; then
  echo "nginx/ssl/*.pem already exist. Re-run with: $0 -f"
  exit 0
fi

TMP=$(mktemp)
trap 'rm -f "$TMP"' EXIT
cat >"$TMP" <<'EOF'
[req]
distinguished_name = dn
x509_extensions = v3_req
prompt = no

[dn]
CN = clickly.cv

[v3_req]
subjectAltName = @alt_names
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment

[alt_names]
DNS.1 = clickly.cv
DNS.2 = *.clickly.cv
DNS.3 = www.clickly.cv
EOF

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout "$SSL/privkey.pem" \
  -out "$SSL/fullchain.pem" \
  -config "$TMP"

echo "Wrote $SSL/fullchain.pem (SAN: clickly.cv, *.clickly.cv, www.clickly.cv)"
echo "Restart nginx: docker compose restart nginx"
