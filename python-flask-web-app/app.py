from flask import Flask , render_template ,request,redirect
import urllib.request , json
import requests
import configparser
from waitress import serve


app = Flask(__name__)


parser = configparser.ConfigParser()
configFilePath = r'/config/web-app-config.properties'
parser.read(configFilePath)
apiScheme=parser['todo-api']['Scheme']
apiHost=parser['todo-api']['Host']
apiPort=parser['todo-api']['Port']
flaskThread=parser['frontend']['ParallelThreads']


@app.route('/',methods=['GET','POST'])
def hello_world():
    if request.method == "POST":
       url = apiScheme+"://"+apiHost+":"+apiPort+"/create"
       input_payload = {"TaskHeader":request.form['taskhead'],"TaskBody":request.form['taskbody']}
       headers = {'Content-Type': 'application/json'}       
       jsondata = json.dumps(input_payload).encode("utf-8")
       requests.post(url, data=jsondata, headers=headers)
           
    url = apiScheme+"://"+apiHost+":"+apiPort+"/show"
    response = urllib.request.urlopen(url)
    data = response.read()
    jsondata = json.loads(data)
    if jsondata["status"] == "success":
       return render_template('index.html',alltasks=jsondata["tasks"])
    else:
       return render_template('index.html',alltasks="")


@app.route('/update',methods=['POST'])
def update_task():
    taskid = request.form['taskid']
    url = apiScheme+"://"+apiHost+":"+apiPort+"/update/"+taskid
    update_payload = {"TaskHeader":request.form['taskhead'],"TaskBody":request.form['taskbody']}
    headers = {'Content-Type': 'application/json'}    
    jsondata = json.dumps(update_payload).encode("utf-8") 
    requests.post(url, data=jsondata, headers=headers)  
    return redirect("/")


@app.route('/delete/<task_id>')
def delete_task(task_id):
    taskid = task_id
    url = apiScheme+"://"+apiHost+":"+apiPort+"/delete/"+taskid
    requests.delete(url)
    return redirect("/")
    

if __name__ == "__main__":
  serve(app,host='0.0.0.0',port=4200,threads=flaskThread)
  #app.run(debug=False,port=4200)

