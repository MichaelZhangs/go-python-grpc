
import requests

headers = {
    "contentType": "application/json"
}

def register(port,name, id, address):
    url = "http://192.168.237.176:8500/v1/agent/service/register"
    rsp =  requests.put(url, headers=headers,json={
        "Name": name,
        "ID": id,
        "Tags": ["mxshop", "web"],
        "Address": address,
        "Port": port,
        "Check":{
            "HTTP": f"http://{address}:{port}/health",
            "Timeout": "5s",
            "Interval": "5s",
            "DeregisterCriticalServiceAfter": "5s"
        }

    })
    print(rsp.status_code)
    if rsp.status_code == 200:
        print("注册成功")
    else:
        print("注册失败")

def deregister(id):
    url = f"http://192.168.237.176:8500/v1/agent/service/deregister/{id}"
    rsp = requests.put(url, headers=headers)
    if rsp.status_code == 200:
        print("注销成功")
    else:
        print(f"注销失败：{rsp.status_code}")
def filter_service(name):
    url = "http://192.168.237.176:8500/v1/agent/services"
    params = {
        "filter": f'Service == "{name}"'
    }
    rsp = requests.get(url, params=params).json()
    for key, value in rsp.items():
        print(key)

def send_to_user():

    for i in range(20):
        url = f"http://127.0.0.1:8888/api_scheduler/v1/test?user_id=abc{i}"
        rsp = requests.get(url=url).json()
        print(rsp)


if __name__ == '__main__':
    # register( 8083,"user_srv","user_srv", "192.168.1.107",)
    # send_to_user()
    deregister("user_srv")