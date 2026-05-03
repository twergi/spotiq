FROM node:20-alpine3.18 AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci --only=production && npm cache clean --force

COPY . .

RUN npm run build

FROM node:20-alpine3.18 AS runner

WORKDIR /app

RUN npm install -g serve@14.2.6 && npm cache clean --force

COPY --from=builder /app/build ./build

RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001 && \
    chown -R nodejs:nodejs /app/build

USER nodejs

EXPOSE 3000

CMD ["serve", "-s", "build", "-l", "3000"]