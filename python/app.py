# app.py
from flask import Flask
import random
import time

app = Flask(__name__)

@app.route('/',methods=['GET'])
def hello():
    return "Hello from Flask!"

@app.route('/metrics',methods=['GET'])
def metrics():
    return f'custom_metric_1 {random.randint(0, 100)}\n'
    #return f'custom_metric_1 5\n'

if __name__ == "__main__":
  app.run(debug=False,host="0.0.0.0", port=8080)
