# app.py
from flask import Flask
import random
import time

app = Flask(__name__)

@app.route('/')
def hello():
    return "Hello from Flask!"

@app.route('/metrics')
def metrics():
    return f'custom_metric_1 {random.randint(0, 100)}\n'
