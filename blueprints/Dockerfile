FROM node:6

RUN npm install -g aglio
COPY . .

CMD ["bash", "-c", "aglio -h 0.0.0.0 -i main.apib --theme-variables slate -s"]
