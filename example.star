
load('indexer.star', 'indexer')
load('schema.star', 'schema')

schema.add("token", [
    {"name": "tokenid", "type": schema.int64}
])

def impl(evnt):
    print(evnt)
    
    token = indexer.get("token", "id1")
    token.set("val", "1")

indexer.index({
    "event": "",
    "impl": impl
})
