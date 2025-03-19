FROM  python:3.9.17-alpine3.18 
RUN mkdir -p /flask /config
WORKDIR /flask
RUN pip install --upgrade pip
COPY requirements.txt requirements.txt
COPY web-app-config.properties /config/web-app-config.properties
RUN pip3 install -r requirements.txt
COPY . .
CMD [ "python3", "/flask/app.py"]
