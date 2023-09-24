FROM node:20-alpine3.18

WORKDIR /source
COPY . .

RUN npm install
RUN npm install -g serve@14.2.1
RUN npm run build