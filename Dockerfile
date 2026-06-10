FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
COPY --from=frontend-builder /app/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o server .

FROM scratch
COPY --from=builder /app/server /server
EXPOSE 80
CMD ["/server"]
