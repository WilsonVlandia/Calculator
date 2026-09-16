# --- build stage ---
FROM node:20-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

# VITE_API_BASE_URL must be a build-time ARG, not a runtime ENV: Vite
# inlines every import.meta.env.VITE_* value directly into the
# compiled JS during `vite build` (there is no server process left at
# runtime to read an environment variable from — the final stage below
# is just nginx serving static files). Setting it as an ENV here makes
# it visible to the `npm run build` step that follows.
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN npm run build

# --- final stage ---
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html

# Matches the container port docker-compose.yml maps for this service.
EXPOSE 80
