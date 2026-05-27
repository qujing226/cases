package Dynamicpb

import ()

func main() {
	md := Person{}
	// 假设你已经有一个 MessageDescriptor
	var md protoreflect.MessageDescriptor = getDescriptorSomehow()

	// 创建动态消息
	msg := dynamicpb.NewMessage(md)

	// 设置字段值
	msg.Set(md.Fields().ByName("id"), protoreflect.ValueOfInt32(123))

	// 序列化
	bytes, _ := proto.Marshal(msg)

	// 反序列化
	_ = proto.Unmarshal(bytes, msg)

	// 读取字段
	val := msg.Get(md.Fields().ByName("id")).Int()

}
