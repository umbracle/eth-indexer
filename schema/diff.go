package schema

type Diff struct {
	// table to which the difference belongs
	Table string `protobuf:"bytes,1,opt,name=table,proto3" json:"table,omitempty"`
	// whether this entry is new or its an update
	Creation bool `protobuf:"varint,2,opt,name=creation,proto3" json:"creation,omitempty"`
	// primary keys of the object
	Keys map[string]string `protobuf:"bytes,3,rep,name=keys,proto3" json:"keys,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// new values of the object
	Vals map[string]string `protobuf:"bytes,4,rep,name=vals,proto3" json:"vals,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}
