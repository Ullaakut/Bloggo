FROM nginx:alpine

COPY index.html /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY grid.scss .

RUN apk update && \
    apk upgrade && \
    apk add nodejs
RUN npm install -g sass
RUN sass grid.scss /usr/share/nginx/html/grid.css

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
