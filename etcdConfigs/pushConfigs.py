import etcd3
import json

def main():
    f = open("config.json", "r")
    data = json.loads(f.read())
    f.close()
    
    etcd = etcd3.client(host='localhost', port=2379)
    for k, v in data.items():
        etcd.put(k, json.dumps(v))

if __name__ == "__main__":
    main()
